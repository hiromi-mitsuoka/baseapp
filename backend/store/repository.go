package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hiromi-mitsuoka/baseapp/clock"
	"github.com/hiromi-mitsuoka/baseapp/config"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Clocker clock.Clocker
}

// NOTE* インターフェース経由でのみメソッドの実行を可能にする
type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// NOTE: For preparer statement
type Preparer interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

// NOTE: 「参照のみのため，DBへの更新なし」と，インターフェースを分割すると，コードリーディングが容易に
type Queryer interface {
	Preparer
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
	QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...any) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error
}

var (
	// インターフェースが期待通りに宣言されているか確認
	_ Beginner = (*sqlx.DB)(nil)
	_ Preparer = (*sqlx.DB)(nil)
	_ Queryer  = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.DB)(nil)
	_ Execer   = (*sqlx.Tx)(nil)
)

func New(ctx context.Context, cfg *config.Config) (*sqlx.DB, func(), error) {
	// sqlx.Connectを使用すると内部でpingする
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			// NOTE: parseTime=trueがないと，time.Time型のフィールドに正しい時刻情報が取得できない
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBName,
		),
	)
	if err != nil {
		return nil, nil, err
	}
	// https://qiita.com/taizo/items/69d3de8622eabe8da6a2#timeout%E3%81%A8cancel
	// 特定時間になったタイミングでキャンセルされる
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	// Openは実際に接続テストを行っていない為，PingContextを利用して明示的に疎通確認
	if err := db.PingContext(ctx); err != nil {
		return nil, func() { _ = db.Close() }, err
	}
	// https://pkg.go.dev/github.com/jmoiron/sqlx@v1.3.5#NewDb
	// NewDb returns a new sqlx DB wrapper for a pre-existing *sql.DB.
	xdb := sqlx.NewDb(db, "mysql")
	// NOTE: *sql.DB型の値はRDBMSの利用終了後にCloseしてコネクションを正しく終了する必要あり
	//       New関数呼び出し元で終了処理できるよう，戻り値に*sql.DB.Closeメソッドを実行する無名関数を返しておく
	return xdb, func() { _ = db.Close() }, nil
}
