package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hiromi-mitsuoka/baseapp/config"
)

func TestNewMux(t *testing.T) {
	// NOTE: httptestパッケージを使って，ServeHTTP関数の引数に渡すためのモックを生成
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	// TODO: NewMux関数をパスするために，ctx, cfgを作成したが，利用の仕方が正しいか不明な状態
	//       cfgのエラーチェックをmux_test.goで行っていることは違う気がする
	ctx, _ := context.WithCancel(context.Background())
	cfg, err := config.New()
	if err != nil {
		t.Error(err)
	}
	sut, _, err := NewMux(ctx, cfg)
	if err != nil {
		t.Error(err)
	}

	sut.ServeHTTP(w, r)
	resp := w.Result()
	t.Cleanup(func() { _ = resp.Body.Close() })

	if resp.StatusCode != http.StatusOK {
		t.Error("want status code 200, but", resp.StatusCode)
	}

	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	want := `{"status": "OK"}`
	if string(got) != want {
		// https://qiita.com/rock619/items/14eb2b32f189514b5c3c#q-1
		// %q : Goの文法上のエスケープをした文字列
		t.Errorf("want %q, but got %q", want, got)
	}
}
