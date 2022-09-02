package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/hiromi-mitsuoka/baseapp/auth"
	"github.com/hiromi-mitsuoka/baseapp/clock"
	"github.com/hiromi-mitsuoka/baseapp/config"
	"github.com/hiromi-mitsuoka/baseapp/handler"
	"github.com/hiromi-mitsuoka/baseapp/service"
	"github.com/hiromi-mitsuoka/baseapp/store"
)

// https://qiita.com/huji0327/items/c85affaf5b9dbf84c11e
// muxとは
// マルチプレクサ: ふたつ以上の入力をひとつの信号として出力する機構．通信分野では多重通信の入口の装置．muxと略される．
// ハンドラ（Go言語における）: ServeHTTPというメソッドを持ったインターフェースのこと．

// https://www.twihike.dev/docs/golang-web/handlers
// マルチプレクサ：リクエストを受け付け、URLに対応するハンドラへ転送する
// ハンドラ：クリエストに応じた処理をして、レスポンスを返却する

// NOTE: 戻り値を*http.ServeMux型の値ではなく，http.Handlerインターフェースにしておくことで，内部実装に依存しない関数シグネチャになる
func NewMux(ctx context.Context, cfg *config.Config) (http.Handler, func(), error) {
	// mux := http.NewServeMux()
	// NOTE: 標準パッケージhttp.ServeMux型の場合ルーティング設定の表現に乏しいため，
	//       ルーティングのみの機能を提供する,net/httpパッケージの型定義に準拠するgo-chi/chiパッケージを利用
	mux := chi.NewRouter()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// NOTE: 静的解析のエラーを回避するため明示的に戻り値を捨てている
		_, _ = w.Write([]byte(`{"status": "OK"}`))
	})

	v := validator.New()
	clocker := clock.RealClocker{}

	// RDBMS
	db, cleanup, err := store.New(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	r := store.Repository{Clocker: clocker}

	// Redis
	rcli, err := store.NewKVS(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}

	// JWT
	jwter, err := auth.NewJWTer(rcli, clocker)
	if err != nil {
		return nil, cleanup, err
	}

	// /login
	l := &handler.Login{
		Service: &service.Login{
			DB:             db,
			Repo:           &r,
			TokenGenerator: jwter,
		},
		Validator: v,
	}
	mux.Post("/login", l.ServeHTTP)

	// /tasks
	at := &handler.AddTask{
		Service:   &service.AddTask{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/tasks", at.ServeHTTP)
	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: &r},
	}
	mux.Get("/tasks", lt.ServeHTTP)

	// /register
	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

	return mux, cleanup, nil
}
