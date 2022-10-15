package service

import (
	"context"

	"github.com/hiromi-mitsuoka/baseapp/entity"
	"github.com/hiromi-mitsuoka/baseapp/store"
)

// NOTE: storeパッケージの直接参照を避けるためのインターフェース

//go:generate go run github.com/matryer/moq -out moq_test.go . TaskAdder TaskLister TaskUpdater TaskDeleter UserRegister UserRegister UserGetter TokenGenerator
type TaskAdder interface {
	AddTask(ctx context.Context, db store.Execer, t *entity.Task) error
}

type TaskLister interface {
	ListTasks(ctx context.Context, db store.Queryer, id entity.UserID) (entity.Tasks, error)
}

type TaskUpdater interface {
	UpdateTask(ctx context.Context, db store.Execer, t *entity.Task, tid int64, uid entity.UserID) error
}

type TaskDeleter interface {
	DeleteTask(ctx context.Context, db store.Execer, tid int64, uid entity.UserID) error
}

type UserRegister interface {
	RegisterUser(ctx context.Context, db store.Execer, u *entity.User) error
}

type UserGetter interface {
	GetUser(ctx context.Context, db store.Queryer, name string) (*entity.User, error)
}

type TokenGenerator interface {
	GenerateToken(ctx context.Context, u entity.User) ([]byte, error)
}

// admin
type AdminTaskLister interface {
	AdminListTask(ctx context.Context, db store.Queryer) (entity.Tasks, error)
}
