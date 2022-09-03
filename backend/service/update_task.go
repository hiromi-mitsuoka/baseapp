package service

import (
	"context"
	"fmt"
	"log"

	"github.com/hiromi-mitsuoka/baseapp/auth"
	"github.com/hiromi-mitsuoka/baseapp/entity"
	"github.com/hiromi-mitsuoka/baseapp/store"
)

type UpdateTask struct {
	DB   store.Execer
	Repo TaskUpdater
}

// まずはtitleの変更のみ
func (u *UpdateTask) UpdateTask(ctx context.Context, title string, tid int64) (*entity.Task, error) {
	log.Printf("====== Service UpdateTask =======")
	uid, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("user_id not found")
	}
	t := &entity.Task{
		Title: title,
	}
	err := u.Repo.UpdateTask(ctx, u.DB, t, tid, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to edit task: %w", err)
	}
	return t, nil
}
