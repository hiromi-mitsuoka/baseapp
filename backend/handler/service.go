package handler

import (
	"context"

	"github.com/hiromi-mitsuoka/baseapp/entity"
)

// NOTE: 密結合になっていたhandlerパッケージからビジネスロジックと永続化処理を取り除く
//       → リクエストの解釈とレスポンスを組み立てる処理のみに

// NOTE: 構造体や関数ではなく，インターフェースを定義する2つの理由
// 1. 他のパッケージへの参照を取り除いて，疎なパッケージにするため
// 2. インターフェースを介して特定の型に依存しないことで，モックに処理を入れ替えたテストを行うため

// https://qiita.com/yaegashi/items/d1fd9f7d0c75b2bb7446
// NOTE: ソースコードを自動生成するコマンド
//go:generate go run github.com/matryer/moq -out moq_test.go . ListTasksService AddTaskService RegisterUserService LoginService
type ListTasksService interface {
	ListTasks(ctx context.Context) (entity.Tasks, error)
}

type AddTaskService interface {
	AddTask(ctx context.Context, title string) (*entity.Task, error)
}

type RegisterUserService interface {
	RegisterUser(ctx context.Context, name, password, role string) (*entity.User, error)
}

type LoginService interface {
	Login(ctx context.Context, name, pw string) (string, error)
}

// TODO: このファイルはなぜservice.go??
