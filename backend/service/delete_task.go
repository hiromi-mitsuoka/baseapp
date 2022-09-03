package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hiromi-mitsuoka/baseapp/auth"
	"github.com/hiromi-mitsuoka/baseapp/entity"
	"github.com/hiromi-mitsuoka/baseapp/store"
)

type DeleteTask struct {
	DB   store.Execer
	Repo TaskDeleter
}

func (d *DeleteTask) DeleteTask(ctx context.Context, tid int64) (int64, error) {
	uid, ok := auth.GetUserID(ctx)
	if !ok {
		// https://qiita.com/uenosy/items/ba9dbc70781bddc4a491#delete
		return http.StatusNotFound, fmt.Errorf("user_id not found")
	}
	err := d.Repo.DeleteTask(ctx, d.DB, tid, entity.UserID(uid))
	if err != nil {
		// https://qiita.com/uenosy/items/ba9dbc70781bddc4a491#delete
		return http.StatusUnprocessableEntity, fmt.Errorf("failed to delete task: %w", err)
	}
	// https://developer.mozilla.org/ja/docs/Web/HTTP/Methods/DELETE
	// 204 (No Content) は、処理は完了しておりかつ、さらなる情報が提供されない場合のステータスコード
	return http.StatusNoContent, nil
}
