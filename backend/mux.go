package main

import (
	"context"
	"log"
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

	// Elasticsearch
	_, res, err := store.NewES(cfg)
	if err != nil {
		return nil, cleanup, err
	}
	// TODO: res.Body.Close() もmysql同様にcleanupを用意するべきか検討
	defer res.Body.Close()
	log.Println(res)

	// JWT
	jwter, err := auth.NewJWTer(rcli, clocker)
	if err != nil {
		return nil, cleanup, err
	}

	// /register
	ru := &handler.RegisterUser{
		Service:   &service.RegisterUser{DB: db, Repo: &r},
		Validator: v,
	}
	mux.Post("/register", ru.ServeHTTP)

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
	lt := &handler.ListTask{
		Service: &service.ListTask{DB: db, Repo: &r},
	}
	ut := &handler.UpdateTask{
		Service:   &service.UpdateTask{DB: db, Repo: &r},
		Validator: v,
	}
	dt := &handler.DeleteTask{
		Service:   &service.DeleteTask{DB: db, Repo: &r},
		Validator: v,
	}
	// ログインユーザーのみ認可するため，サブルーター作成
	mux.Route("/tasks", func(r chi.Router) {
		// chi.Router.Use: サブルーターのエンドポイント全体にミドルウェアを適用
		r.Use(handler.AuthMiddleware(jwter))
		r.Post("/", at.ServeHTTP)
		r.Get("/", lt.ServeHTTP)
		// https://github.com/go-chi/chi#examples
		r.Route("/{taskID}", func(r chi.Router) {
			r.Use(handler.TaskCtx)
			r.Put("/", ut.ServeHTTP)
			r.Delete("/", dt.ServeHTTP)
		})
	})

	mux.Route("/admin", func(r chi.Router) {
		// NOTE: ミドルウェアの適用順序注意
		r.Use(handler.AuthMiddleware(jwter), handler.AdminMiddleware)
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			_, _ = w.Write([]byte(`{"message": "admin only"}`))
		})
	})

	return mux, cleanup, nil
}
