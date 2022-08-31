package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hiromi-mitsuoka/baseapp/entity"
)

type AddTask struct {
	// NOTE: 永続化処理をAddTaskServiceインターフェース型に委譲
	// DB        *sqlx.DB
	// Repo      *store.Repository
	Service   AddTaskService
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// NOTE: JSONの構造が巨大・複雑だった場合，Unmarshalでは検証が大変なため，Validatorパッケージを使用
	if err := at.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// NOTE: 永続化処理をAddTaskServiceインターフェース型に委譲
	// t := &entity.Task{
	// 	Title:  b.Title,
	// 	Status: entity.TaskStatusTodo,
	// }
	t, err := at.Service.AddTask(ctx, b.Title)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		ID entity.TaskID `json:"id"`
	}{ID: t.ID}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
