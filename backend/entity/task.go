package entity

import (
	"context"
	"strconv"
	"time"
)

// NOTE: 独自の型を挟むことで，謝った代入を防ぐ
type TaskID int64
type TaskStatus string

const (
	TaskStatusTodo  TaskStatus = "todo"
	TaskStatusDoing TaskStatus = "doing"
	TaskStatusDone  TaskStatus = "done"
)

// NOTE: github.com/jmoiron/sqlxパッケージを利用する場合は，構造体の各フィールドにタグでテーブルカラム名に対応したメタデータを設定する
type Task struct {
	ID       TaskID     `json:"id" db:"id"`
	UserID   UserID     `json:"user_id" db:"user_id"`
	Title    string     `json:"title" db:"title"`
	Status   TaskStatus `json:"status" db:"status"`
	Created  time.Time  `json:"created" db:"created"`
	Modified time.Time  `json:"modified" db:"modified"`
}

type Tasks []*Task

type TaskIDKey struct{}

// TODO: 共通処理を記述すべきファイルはどこ？確認?
func GetTaskID(ctx context.Context) (int64, error) {
	// TODO: 上手にentity.TaskIDにキャストできない
	tid_key := ctx.Value(TaskIDKey{}).(string)
	// TODO: SQLインジェクション対策になっている？
	tid, err := strconv.ParseInt(tid_key, 10, 64)
	return tid, err
}
