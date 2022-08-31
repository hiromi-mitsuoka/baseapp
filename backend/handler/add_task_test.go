package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/hiromi-mitsuoka/baseapp/entity"
	"github.com/hiromi-mitsuoka/baseapp/testutil"
)

// NOTE: Table Driven Test
// 複数の入力や期待値の組み合わせを共通化した実行手順で実行させるテストの実装パターン

// NOTE: ゴールデンテスト
// テストの入力や期待値を別ファイルとして保存したテストのこと

func TestAddTask(t *testing.T) {
	t.Parallel()
	type want struct {
		status  int
		rspFile string
	}
	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"ok": {
			reqFile: "testdata/add_task/ok_req.json.golden",
			want: want{
				status:  http.StatusOK,
				rspFile: "testdata/add_task/ok_rsp.json.golden",
			},
		},
		"badRequest": {
			reqFile: "testdata/add_task/bad_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/add_task/bad_req_rsp.json.golden",
			},
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/tasks",
				bytes.NewReader(testutil.LoadFile(t, tt.reqFile)),
			)

			// NOTE: matryer/moqで自動生成したAddTaskServiceMock型を使用
			moq := &AddTaskServiceMock{}
			moq.AddTaskFunc = func(ctx context.Context, title string) (*entity.Task, error) {
				if tt.want.status == http.StatusOK {
					return &entity.Task{ID: 1}, nil
				}
				return nil, errors.New("error from mock")
			}

			sut := AddTask{
				// NOTE: 永続化処理をAddTaskServiceインターフェース型に委譲した関係で削除
				// Store: &store.TaskStore{
				// 	Tasks: map[entity.TaskID]*entity.Task{},
				// },
				Service:   moq,
				Validator: validator.New(),
			}
			sut.ServeHTTP(w, r)

			resp := w.Result()
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
