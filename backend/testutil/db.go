package testutil

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
)

// NOTE: CI環境変数はGitHub Actions上しか定義されていない想定
// 	     ローカルマシン環境やGitHub Actions上の環境に対して，ポート番号を切り替える

func OpenDBForTest(t *testing.T) *sqlx.DB {
	t.Helper()

	port := 33306
	// https://hawksnowlog.blogspot.com/2019/04/lookup-env-golang.html
	// 環境変数「CI」が設定してある場合は，3306を使用
	if _, defined := os.LookupEnv("CI"); defined {
		port = 3306
	}
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("user:password@tcp(127.0.0.1:%d)/baseapp?parseTime=true", port),
	)
	if err != nil {
		// https://simple-minds-think-alike.moritamorie.com/entry/go-testing-error-fatal
		// t.Fatalf() : FailNow→ 対象の関数のテストに失敗した記録を残し、後続のテストは実行しない。
		t.Fatal(err)
	}
	// https://qiita.com/wifecooky/items/cf8ae640eb526e7413eb
	// deferは並列で動いているサブテストの終了を待たずに実行されてしまう問題あり
	// 全てのサブテストが終了してから後処理を行たいなら t.Cleanup()
	t.Cleanup(
		func() { _ = db.Close() },
	)
	return sqlx.NewDb(db, "mysql")
}
