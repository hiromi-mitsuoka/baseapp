package service

import (
	"context"
	"fmt"

	"github.com/hiromi-mitsuoka/baseapp/entity"
	"github.com/hiromi-mitsuoka/baseapp/store"
)

// NOTE: storeパッケージの特定の方に依存せずインターフェースをDIする設計

// DI: Dependency Injection(オブジェクトの注入)
// https://qiita.com/yoshinori_hisakawa/items/a944115eb77ed9247794
// 処理に必要なオブジェクトを外部から注入できるようにするデザインパターン
// オブジェクトをインタフェースとして定義し、使う側は実装オブジェクトでなく、インタフェースを利用するようにする。
// 実装オブジェクトは外部からそのインタフェースに外部から注入する事で、実装を入れ替えたりできる。

type AddTask struct {
	DB   store.Execer
	Repo TaskAdder
}

func (a *AddTask) AddTask(ctx context.Context, title string) (*entity.Task, error) {
	t := &entity.Task{
		Title:  title,
		Status: entity.TaskStatusTodo,
	}
	err := a.Repo.AddTask(ctx, a.DB, t)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return t, nil
}
