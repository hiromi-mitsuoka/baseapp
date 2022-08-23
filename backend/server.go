package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	srv *http.Server
	l   net.Listener
}

func NewServer(l net.Listener, mux http.Handler) *Server {
	return &Server{
		srv: &http.Server{Handler: mux},
		l:   l,
	}
}

// https://zenn.dev/hsaki/books/golang-context/viewer/definition
// Context型
// 処理の締め切り・キャンセル信号・API境界やプロセス間を横断する必要のあるリクエストスコープな値を伝達させる
// 基本的に、Goでは「異なるゴールーチン間での情報共有は、ロックを使ってメモリを共有するよりも、チャネルを使った伝達を使うべし」という考え方を取っている
// 「複数ゴールーチン間で安全に、そして簡単に情報伝達を行いたい」という要望は、チャネルによる伝達だけ実現しようとすると難しい
// → ゴールーチン上で起動される関数の第一引数に、context.Context型を1つ渡す」だけで簡単に実現できるようになっている
func (s *Server) Run(ctx context.Context) error {
	// グレースフルシャットダウン
	// 何らかの処理を実行中に終了シグナルを受け付けた場合，アプリケーションプロセスは処理を正しく終了させるまで終了しないことが望ましい．
	// http.Server型はShutdownメソッドを呼ぶと，グレースフルシャットダウンを開始する
	// 割り込みシグナル（SIGINT）, 終了シグナル（SIGTERM）
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// https://zenn.dev/ikawaha/articles/20211218-f37638b56e5807#golang.org%2Fx%2Fsync%2Ferrgroup
	// 複数の goroutine を実行して、それらのうちにエラーがあったときにエラーを知る、ということを可能にしてくれるライブラリ
	// 標準パッケージのsync.WaitGroup型の場合は，別ゴルーチン上で実行する関数から戻り値でエラーを受け取ることができない
	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTPサーバーを起動する
	eg.Go(func() error {
		// http.ErrServerClosedは，http.Server.Shutdown()が正常に終了したことを示すので異常ではない
		// errgroup は sync.WaitGroup+error といったイメージで、どれかの goroutine でエラーがあったら最初のひとつを知ることができる
		// ctx を組み合わせても使えるようになっているので、goroutine のどれかがエラーになったら処理を切り上げる、という使い方ができて便利
		if err := s.srv.Serve(s.l); err != nil &&
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
	if err := s.srv.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	// Goメソッドで起動した別ゴルーチンの終了を待つ
	// https://zenn.dev/ikawaha/articles/20211218-f37638b56e5807#golang.org%2Fx%2Fsync%2Ferrgroup
	// すべての goroutine が終わるのを待って、エラーが発生していれば（最初の）エラーを返す
	return eg.Wait()
}
