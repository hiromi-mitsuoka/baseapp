package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hiromi-mitsuoka/baseapp/entity"
)

type UpdateTask struct {
	Service   UpdateTaskService
	Validator *validator.Validate
}

func (ut *UpdateTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tid, err := entity.GetTaskID(ctx)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	var b struct {
		Title  string            `json:"title"`
		Status entity.TaskStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	if err := ut.Validator.Struct(b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	t, err := ut.Service.UpdateTask(ctx, tid, b.Title, b.Status)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	// TODO: Update時の適切な返し方確認
	rsp := struct {
		ID entity.TaskID `json:"id"`
	}{ID: t.ID}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
