package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hiromi-mitsuoka/baseapp/entity"
)

type DeleteTask struct {
	Service   DeleteTaskService
	Validator *validator.Validate
}

func (dt *DeleteTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tid, err := entity.GetTaskID(ctx)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	status_code, err := dt.Service.DeleteTask(ctx, tid)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := struct {
		StatusCode int64 `json:"status_code"`
	}{StatusCode: status_code}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
