package auth

import (
	"context"
	_ "embed"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hiromi-mitsuoka/baseapp/clock"
	"github.com/hiromi-mitsuoka/baseapp/entity"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const (
	RoleKey     = "role"
	UserNameKey = "user_name"
)

// https://future-architect.github.io/articles/20210208/
// go:embed
// ファイル埋め込みをサポートするためのパッケージ
// 単一のファイルを埋め込みするだけなら，_ "embed"として先頭に_をつけてインポートすることが推奨
// 埋め込みファイルの場所を指示する記述

// os.ReadFile を使うと，実行バイナリの他にファイルも適切なファイルパスで実行環境に展開しておく運用が必要になる
// → Goのシングルバイナリで実行可能なメリットがなくなる → go:embed の場合，変数に鍵ファイルの内容の埋め込み可能
// → go build コマンドによって生成されたシングルバイナリを実行環境に配置するだけでデプロイ完了

// NOTE:
// 実際は，開発環境・本番環境で利用する鍵ファイルは異なるものにしたい
// go1.18時点では，環境変数などを使って指定するファイル名を変更できない
// 実際のデプロイでは，デプロイパイプライン上でgo build コマンド実行前に，その環境の鍵ファイルを指定のパスに配置するステップの作成が必要

//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte

type JWTer struct {
	PrivateKey, PublicKey jwk.Key
	// DI（依存性注入: Dependency Injection）する構造に
	// https://qiita.com/okazuki/items/a0f2fb0a63ca88340ff6
	// テストはできる限り単体テストにしたい → 疎結合な実装に
	// → 内部で依存先を new するのではなく、外から依存先の実装を設定してもらうという考え方が Dependency (依存性) Injection (注入)
	Store   Store
	Clocker clock.Clocker
}

//go:generate go run github.com/matryer/moq -out moq_test.go . Store
type Store interface {
	Save(ctx context.Context, key string, userID entity.UserID) error
	Load(ctx context.Context, key string) (entity.UserID, error)
}

// NOTE: アプリケーション起動時に，「鍵」として読み込んだデータを保持する
func NewJWTer(s Store, c clock.Clocker) (*JWTer, error) {
	j := &JWTer{Store: s}
	privkey, err := parse(rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: private key: %w", err)
	}
	pubkey, err := parse(rawPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed in NewJWTer: public key: %w", err)
	}
	j.PrivateKey = privkey
	j.PublicKey = pubkey
	j.Clocker = clock.RealClocker{}
	return j, nil
}

// NOTE: lestrrat-go/jwx/v2/jwkパッケージの jwk.ParseKey関数を使って，鍵の情報が含まれるバイト配列から，
//       lestrrat-go/jwxパッケージで利用可能な jwk.Keyインターフェースを満たす型の値を取得
func parse(rawKey []byte) (jwk.Key, error) {
	key, err := jwk.ParseKey(rawKey, jwk.WithPEM(true))
	if err != nil {
		return nil, err
	}
	return key, nil
}

// NOTE: 引数で渡すユーザー情報で，JWTに署名してトークン文字列を作成
func (j *JWTer) GenerateToken(ctx context.Context, u entity.User) ([]byte, error) {
	// https://github.com/lestrrat-go/jwx
	// https://kamichidu.github.io/post/2017/01/24-about-json-web-token/
	tok, err := jwt.NewBuilder().
		JwtID(uuid.New().String()).
		Issuer(`github.com/hiromi-mitsuoka/baseapp`). // 発行者
		Subject("access_token").                      // 用途
		IssuedAt(j.Clocker.Now()).
		Expiration(j.Clocker.Now().Add(30*time.Minute)). // 失効時間
		Claim(RoleKey, u.Role).
		Claim(UserNameKey, u.Name).
		Build()
	if err != nil {
		return nil, fmt.Errorf("GetToken: failed to build token: %w", err)
	}
	if err := j.Store.Save(ctx, tok.JwtID(), u.ID); err != nil {
		return nil, err
	}

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.RS256, j.PrivateKey))
	if err != nil {
		return nil, err
	}
	return signed, nil
}

// NOTE: HTTPリクエストのAuthorizationリクエストヘッダーに付与されているJWTトークンを取得する
func (j *JWTer) GetToken(ctx context.Context, r *http.Request) (jwt.Token, error) {
	token, err := jwt.ParseRequest(
		r,
		// WithKey: 署名を検証するアルゴリズムと利用する鍵を指定
		jwt.WithKey(jwa.RS256, j.PublicKey),
		// NOTE: *auth.JWTer.Clocker フィールドをベースに検証を行うため，jwt.WithValidate関数の検証は無効化
		jwt.WithValidate(false),
	)
	if err != nil {
		return nil, err
	}
	if err := jwt.Validate(token, jwt.WithClock(j.Clocker)); err != nil {
		return nil, fmt.Errorf("GetToken: failed to validate token: %w", err)
	}
	// Redisから削除して手動でexpireさせていることも想定
	if _, err := j.Store.Load(ctx, token.JwtID()); err != nil {
		return nil, fmt.Errorf("GetToken: %q expired: %w", token.JwtID(), err)
	}
	return token, nil
}

// NOTE: アプリケーションコードで常にJWTを扱うのは冗長なため，context.Context型の値にJWTから取得したユーザーIDとロール権限を設定する

// context.WithValue で値設定する時，Defined Type で指定することで誤入力を防ぐ
type userIDKey struct{}
type roleKey struct{}

func SetUserID(ctx context.Context, uid entity.UserID) context.Context {
	return context.WithValue(ctx, userIDKey{}, uid)
}

func GetUserID(ctx context.Context) (entity.UserID, bool) {
	id, ok := ctx.Value(userIDKey{}).(entity.UserID)
	return id, ok
}

func SetRole(ctx context.Context, tok jwt.Token) context.Context {
	get, ok := tok.Get(RoleKey)
	if !ok {
		return context.WithValue(ctx, roleKey{}, "")
	}
	return context.WithValue(ctx, roleKey{}, get)
}

func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(roleKey{}).(string)
	return role, ok
}

func (j *JWTer) FillContext(r *http.Request) (*http.Request, error) {
	token, err := j.GetToken(r.Context(), r)
	if err != nil {
		return nil, err
	}
	uid, err := j.Store.Load(r.Context(), token.JwtID())
	if err != nil {
		return nil, err
	}
	ctx := SetUserID(r.Context(), uid)

	ctx = SetRole(ctx, token)
	clone := r.Clone(ctx)
	return clone, nil
}

func IsAdmin(ctx context.Context) bool {
	role, ok := GetRole(ctx)
	if !ok {
		return false
	}
	// adminかどうかは1行で十分
	return role == "admin"
}
