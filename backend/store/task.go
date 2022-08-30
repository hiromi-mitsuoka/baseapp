package store

import (
	"context"

	"github.com/hiromi-mitsuoka/baseapp/entity"
)

// RDBMSへ書き込みを実行するため，Execerインターフェースを指定
func (r *Repository) AddTask(ctx context.Context, db Execer, t *entity.Task) error {
	t.Created = r.Clocker.Now()
	t.Modified = r.Clocker.Now()
	sql := `INSERT INTO task
						(title, status, created, modified)
					VALUES (?, ?, ?, ?)`
	result, err := db.ExecContext(
		ctx,
		sql,
		t.Title,
		t.Status,
		t.Created,
		t.Modified,
	)
	if err != nil {
		return err
	}
	// NOTE: idを明示的に構造体に追加，呼び出し元に発行されたIDを伝える
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = entity.TaskID(id)
	return nil
}

// Repositoryインターフェース型に紐付け
// 参照系のため，定義したQueryerインターフェース型の値を受け取る．Queryerを指定することで，DBの書き込みの可能性がなくなる
func (r *Repository) ListTasks(ctx context.Context, db Queryer) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT
						id, title, status, created, modified
					FROM task;`
	// db.SelectContext: 複数のレコードを取得．各レコードを1つひとつの構造体に代入したスライスを返す．github.com/jmoiron/sqlxパッケージの拡張メソッド
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}
	return tasks, nil
}
