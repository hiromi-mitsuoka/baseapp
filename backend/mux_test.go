package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TEstNewMux(t *testing.T) {
	// NOTE: httptestパッケージを使って，ServeHTTP関数の引数に渡すためのモックを生成
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/health", nil)
	sut := NewMux()
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
		t.Errorf("want %q, got %q", want, got)
	}
}
