package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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

	// CORS
	// https://github.com/go-chi/cors#usage
	// https://qiita.com/taito-ITO/items/f4fdfa6a031b91beb080
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:18080"}, // swagger
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// NOTE: 静的解析のエラーを回避するため明示的に戻り値を捨てている
		_, _ = w.Write([]byte(`{"status": "OK"}`))
	})
	mux.Get("/cors", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("cors"))
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
	escli, res, err := store.NewES(ctx, cfg)
	if err != nil {
		return nil, cleanup, err
	}
	// TODO: res.Body.Close() もmysql同様にcleanupを用意するべきか検討
	defer res.Body.Close()
	log.Println(res)

	// TODO: 実装の見通したったら，別ファイルに移動
	// temp code
	// https://www.elastic.co/jp/blog/the-go-client-for-elasticsearch-introduction
	// indexing
	req := esapi.IndexRequest{
		Index:      "user-index",
		DocumentID: "1",
		Body:       strings.NewReader(`{"name":"test-user"}`),
	}
	req.Do(ctx, escli.Cli)
	req = esapi.IndexRequest{
		Index:      "user-index",
		DocumentID: "2",
		Body:       strings.NewReader(`{"name":"test-userrrr"}`),
	}
	req.Do(ctx, escli.Cli)

	// search
	// https://developer.okta.com/blog/2021/04/23/elasticsearch-go-developers-guide
	var buffer bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": "test",
			},
		},
	}
	json.NewEncoder(&buffer).Encode(query)
	response, err := escli.Cli.Search(
		escli.Cli.Search.WithIndex("user-index"),
		escli.Cli.Search.WithBody(&buffer),
	)
	var result map[string]interface{}
	json.NewDecoder(response.Body).Decode(&result)
	fmt.Println("====", result)
	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		craft := hit.(map[string]interface{})["_source"].(map[string]interface{})
		fmt.Println("-----", craft)
	}

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
