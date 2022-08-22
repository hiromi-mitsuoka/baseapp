package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {
	// NOTE: 初期のHTTPサーバーを起動するだけの実装
	// err := http.ListenAndServe(
	// 	":8000", // localhostを省略した記法
	// 	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		// https://qiita.com/taji-taji/items/77845ef744da7c88a6fe#%E6%8E%A5%E9%A0%AD%E8%BE%9E-f
	// 		// 接頭辞F : 書き込み先を明示的に指定
	// 		// fmt.Println(r.URL)
	// 		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	// 	}),
	// )
	// if err != nil {
	// 	// https://qiita.com/taji-taji/items/77845ef744da7c88a6fe#%E6%8E%A5%E5%B0%BE%E8%BE%9E-f
	// 	// 接尾辞f : フォーマットを指定
	// 	fmt.Printf("failed to terminate server: %v", err)
	// 	// https://qiita.com/umisama/items/7be04949d670d8cdb99c
	// 	// https://tech-up.hatenablog.com/entry/2018/12/13/154143
	// 	// OSに値を返してプロセスを切る．deferも実行されない．
	// 	// 0 = 正常終了．1 = 以上終了
	// 	os.Exit(1)
	// }

	if err := run(context.Background()); err != nil {
		// https://waman.hatenablog.com/entry/2017/09/29/011614#logPrintf-%E9%96%A2%E6%95%B0
		// fmt.Printf関数のように，フォーマットを指定してログメッセージを出力する
		log.Printf("failed to terminate server: %v", err)
	}
}

// https://zenn.dev/hsaki/books/golang-context/viewer/definition
// Context型
// 処理の締め切り・キャンセル信号・API境界やプロセス間を横断する必要のあるリクエストスコープな値を伝達させる
// 基本的に、Goでは「異なるゴールーチン間での情報共有は、ロックを使ってメモリを共有するよりも、チャネルを使った伝達を使うべし」という考え方を取っている
// 「複数ゴールーチン間で安全に、そして簡単に情報伝達を行いたい」という要望は、チャネルによる伝達だけ実現しようとすると難しい
// → ゴールーチン上で起動される関数の第一引数に、context.Context型を1つ渡す」だけで簡単に実現できるようになっている
func run(ctx context.Context) error {
	// net/httpパッケージの，*http.Server型を使用 → サーバーのタイムアウト時間なども設定可能
	// https://pkg.go.dev/net/http#pkg-overview
	// More control over the server's behavior is available by creating a custom Server
	s := &http.Server{
		Addr: ":8000",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	// https://zenn.dev/ikawaha/articles/20211218-f37638b56e5807#golang.org%2Fx%2Fsync%2Ferrgroup
	// 複数の goroutine を実行して、それらのうちにエラーがあったときにエラーを知る、ということを可能にしてくれるライブラリ
	// 標準パッケージのsync.WaitGroup型の場合は，別ゴルーチン上で実行する関数から戻り値でエラーを受け取ることができない
	eg, ctx := errgroup.WithContext(ctx)

	// 別ゴルーチンでHTTPサーバーを起動する
	eg.Go(func() error {
		// http.ErrServerClosedは，http.Server.Shutdown()が正常に終了したことを示すので異常ではない
		// errgroup は sync.WaitGroup+error といったイメージで、どれかの goroutine でエラーがあったら最初のひとつを知ることができる
		// ctx を組み合わせても使えるようになっているので、goroutine のどれかがエラーになったら処理を切り上げる、という使い方ができて便利
		if err := s.ListenAndServe(); err != nil &&
			// https://pkg.go.dev/net/http#example-Server.Shutdown
			// https://qiita.com/t2y/items/acd86fe24a25e996dbda#shutdown-%E3%81%99%E3%82%8B%E3%81%A8-serve-%E3%81%8B%E3%82%89-errserverclosed-%E3%81%8C%E8%BF%94%E3%82%8B
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// チャネルからの通知（終了通知）を待機する
	<-ctx.Done()
	// NOTE: contextを通じて処理の中断命令を検知した際は，Shutdownメソッドで終了する．そのShutdownに失敗した時
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	// Goメソッドで起動した別ゴルーチンの終了を待つ
	// https://zenn.dev/ikawaha/articles/20211218-f37638b56e5807#golang.org%2Fx%2Fsync%2Ferrgroup
	// すべての goroutine が終わるのを待って、エラーが発生していれば（最初の）エラーを返す
	return eg.Wait()
}
