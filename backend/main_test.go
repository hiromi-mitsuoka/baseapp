package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestRun(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("failed to listen port %v", err)
	}
	// https://www.wakuwakubank.com/posts/867-go-context/
	// WithCancel : 新しいコンテキストとともに，キャンセル関数を取得する
	// Background : 空のContextを生成
	ctx, cancel := context.WithCancel(context.Background())
	// https://www.fullstory.com/blog/why-errgroup-withcontext-in-golang-server-handlers/
	// WithContext returns a new Group and an associated Context derived from ctx. The derived Context is canceled the first time a function passed to Go returns a non-nil error or the first time Wait returns, whichever occurs first.
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx, l)
	})
	in := "message"
	url := fmt.Sprintf("http://%s/%s", l.Addr().String(), in)
	// どんなポート番号でリッスンしているのか確認
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)
	if err != nil {
		// https: //simple-minds-think-alike.moritamorie.com/entry/go-testing-error-fatal
		// t.Errorf : Fail→ 対象の関数のテストに失敗した記録を残すが、後続のテストは実行する
		// https://blog.y-yuki.net/entry/2017/05/02/000000
		// %+v : 構造体を出力する際に、%vの内容にフィールド名も加わる
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	// https://budougumi0617.github.io/2021/02/22/update_capacity/
	// io.ReadAll : io.Readerインターフェイスからすべてのデータを読み出す関数
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		// %v : デフォルトフォーマットで対象データの情報を埋め込む
		t.Fatalf("failed to read body: %v", err)
	}

	// HTTPサーバーの戻り値を検証する
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	// run関数に終了通知を送信する
	cancel()
	// run関数の戻り値を検証する
	if err := eg.Wait(); err != nil {
		// https://simple-minds-think-alike.moritamorie.com/entry/go-testing-error-fatal
		// Fatal : FailNow→ 対象の関数のテストに失敗した記録を残し、後続のテストは実行しない
		t.Fatal(err)
	}
}
