package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hiromi-mitsuoka/baseapp/auth"
	"github.com/hiromi-mitsuoka/baseapp/entity"
)

// https://zenn.dev/tanaka_takeru/articles/aecd36a805886d
// 認証(Authentication)： 「あなたは誰ですか？」 を確認
// 認可(Authorization) ： 「あなたには、リソースにアクセスする権限がありますか？」 を確認

// NOTE: アクセストークンを確認し，context.Context型にユーザー情報を埋め込む
//       アクセストークンが見つからなかった場合，リクエストの処理を終了するので認証も兼ねている
func AuthMiddleware(j *auth.JWTer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			req, err := j.FillContext(r)
			if err != nil {
				RespondJSON(r.Context(), w, ErrResponse{
					Message: "not find auth info",
					Details: []string{err.Error()},
				}, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

// NOTE: context.Context型の値にユーザー情報が含まれていることを前提にしているため，ミドルウェアを適用する順序に注意
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !auth.IsAdmin(r.Context()) {
			RespondJSON(r.Context(), w, ErrResponse{
				Message: "not admin",
			}, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// https://github.com/go-chi/chi#examples
// NOTE: taskIDをmiddlewareで取得
func TaskCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tid := chi.URLParam(r, "taskID")
		ctx := context.WithValue(r.Context(), entity.TaskIDKey{}, tid)
		req := r.Clone(ctx)
		next.ServeHTTP(w, req)
	})
}
