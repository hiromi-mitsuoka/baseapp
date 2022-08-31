package service

import (
	"context"

	"github.com/hiromi-mitsuoka/baseapp/entity"
	"github.com/hiromi-mitsuoka/baseapp/store"
)

// NOTE: storeパッケージの直接参照を避けるためのインターフェース

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskAdder TaskLister
type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer) (entity.Tasks, error)
}
