package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// NOTE: HTTPハンドラーのテストで，レスポンスと期待値のJSONを比較してテスト結果を検証
//       大きなJSON構造を文字列として比較すると差異がわかりにくいため，JSON文字列をUnmarshalして差分を比較
func AssertJSON(t *testing.T, want, got []byte) {
	t.Helper()

	var jw, jq any
	if err := json.Unmarshal(want, &jw); err != nil {
		t.Fatalf("cannot unmarshal want %q: %v", want, err)
	}
	if err := json.Unmarshal(got, &jq); err != nil {
		t.Fatalf("cannot unmarshal got %q: %v", got, err)
	}
	// 差分の比較．go-cmp/cmpパッケージを利用
	if diff := cmp.Diff(jq, jw); diff != "" {
		t.Errorf("got differs: (-got +want)\n%s", diff)
	}
}

func AssertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper()
	// https://qiita.com/wifecooky/items/cf8ae640eb526e7413eb
	// deferは並列で動いているサブテストの終了を待たずに実行されてしまう問題あり
	// 全てのサブテストが終了してから後処理を行たいなら t.Cleanup()
	t.Cleanup(func() { _ = got.Body.Close() })
	gb, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}
	if got.StatusCode != status {
		t.Fatalf("want status %d, but got %d, body: %q", status, got.StatusCode, gb)
	}

	if len(gb) == 0 && len(body) == 0 {
		// 期待としても実体とsてもレスポンスボディがないので，AssertJSONを呼ぶ必要はない
		return
	}
	AssertJSON(t, body, gb)
}

// NOTE: ゴールデンテストで使用
func LoadFile(t *testing.T, path string) []byte {
	t.Helper()

	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read from %q: %v", path, err)
	}
	return bt
}
