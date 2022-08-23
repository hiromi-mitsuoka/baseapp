package main

import "net/http"

// https://qiita.com/huji0327/items/c85affaf5b9dbf84c11e
// muxとは
// マルチプレクサ: ふたつ以上の入力をひとつの信号として出力する機構．通信分野では多重通信の入口の装置．muxと略される．
// ハンドラ（Go言語における）: ServeHTTPというメソッドを持ったインターフェースのこと．

// https://www.twihike.dev/docs/golang-web/handlers
// マルチプレクサ：リクエストを受け付け、URLに対応するハンドラへ転送する
// ハンドラ：クリエストに応じた処理をして、レスポンスを返却する

// NOTE: 戻り値を*http.ServeMux型の値ではなく，http.Handlerインターフェースにしておくことで，内部実装に依存しない関数シグネチャになる
func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// NOTE: 静的解析のエラーを回避するため明示的に戻り値を捨てている
		_, _ = w.Write([]byte(`{"status": "OK"}`))
	})
	return mux
}
