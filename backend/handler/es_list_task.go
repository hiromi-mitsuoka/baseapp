package handler

import (
	"encoding/json"
	"net/http"

	"github.com/hiromi-mitsuoka/baseapp/entity"
)

type EsListTask struct {
	Service EsListTaskService
}

type es_task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (elt *EsListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := elt.Service.EsListTask(ctx)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}
	var rsp map[string]interface{}
	json.NewDecoder(tasks.Body).Decode(&rsp)
	// for _, hit := range rsp["hits"].(map[string]interface{})["hits"].([]interface{}) {
	// 	craft := hit.(map[string]interface{})["_source"].(map[string]interface{})
	// 	fmt.Println("-----", craft)
	// }
	// NOTE: 現状は'hits'の中を返す
	RespondJSON(ctx, w, rsp["hits"], http.StatusOK)

}
