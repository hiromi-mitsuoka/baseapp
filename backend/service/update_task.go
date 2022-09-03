package service

import (
	"context"
	"fmt"

	"github.com/hiromi-mitsuoka/baseapp/auth"
	"github.com/hiromi-mitsuoka/baseapp/entity"
	"github.com/hiromi-mitsuoka/baseapp/store"
)

type UpdateTask struct {
	DB   store.Execer
	Repo TaskUpdater
}

func (u *UpdateTask) UpdateTask(ctx context.Context, tid int64, title string, status entity.TaskStatus) (*entity.Task, error) {
	uid, ok := auth.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("user_id not found")
	}
	if status != "todo" && status != "doing" && status != "done" {
		return nil, fmt.Errorf("Inappropriate task status")
	}
	t := &entity.Task{
		Title:  title,
		Status: status,
	}
	err := u.Repo.UpdateTask(ctx, u.DB, t, tid, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to edit task: %w", err)
	}
	return t, nil
}
