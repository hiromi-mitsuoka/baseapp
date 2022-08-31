package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hiromi-mitsuoka/baseapp/entity"
	"github.com/hiromi-mitsuoka/baseapp/testutil"
)

func TestListTask(t *testing.T) {
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		tasks []*entity.Task
		want  want
	}{
		"ok": {
			tasks: []*entity.Task{
				{
					ID:     1,
					Title:  "test1",
					Status: entity.TaskStatusTodo,
				},
				{
					ID:     2,
					Title:  "test2",
					Status: entity.TaskStatusDone,
				},
			},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_task/ok_rsp.json.golden",
			},
		},
		"empty": {
			tasks: []*entity.Task{},
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/list_task/empty_rsp.json.golden",
			},
		},
	}
	for n, tt := range tests {
		// https://zenn.dev/ucwork/articles/cd26d933978080
		// Testxxxとt.Run()をどちらもt.Parallel()で並行処理させると想定外の挙動をする，再宣言することで解消される
		// https://engineering.mercari.com/blog/entry/how_to_use_t_parallel/
		// func TestXXX(t *testing.T)のシグニチャを持つトップレベルのテスト関数と、トップレベルのテスト関数内で、t.Run()を用いて記述するサブテスト関数がある
		tt := tt
		t.Run(n, func(t *testing.T) {
			// https: //qiita.com/marnie_ms4/items/8706f43591fb23dd4e64
			// testing.T.Parallel() を呼び出したテストは並列に実行される
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/tasks", nil)

			moq := &ListTasksServiceMock{}
			moq.ListTasksFunc = func(ctx context.Context) (entity.Tasks, error) {
				if tt.tasks != nil {
					return tt.tasks, nil
				}
				return nil, errors.New("error from mock")
			}
			sut := ListTask{Service: moq}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(
				t,
				resp,
				tt.want.status,
				testutil.LoadFile(t, tt.want.rspFile),
			)
		})
	}
}
