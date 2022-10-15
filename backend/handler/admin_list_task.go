package handler

import (
	"net/http"

	"github.com/hiromi-mitsuoka/baseapp/entity"
)

type AdminListTask struct {
	Service AdminListTaskService
}

// NOTE: admin権限により，user_idも確認できるように作成
type a_task struct {
	ID     entity.TaskID     `json:"id"`
	UserID entity.UserID     `json:"user_id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (alt *AdminListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := alt.Service.AdminListTask(ctx)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	rsp := []a_task{}
	for _, t := range tasks {
		rsp = append(rsp, a_task{
			ID:     t.ID,
			UserID: t.UserID,
			Title:  t.Title,
			Status: t.Status,
		})
	}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
