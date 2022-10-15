package store

import (
	"context"

	"github.com/hiromi-mitsuoka/baseapp/entity"
)

func (r *Repository) AdminListTask(
	ctx context.Context, db Queryer,
) (entity.Tasks, error) {
	tasks := entity.Tasks{}
	sql := `SELECT
						id, user_id, title, status, created, modified
					FROM tasks`
	if err := db.SelectContext(ctx, &tasks, sql); err != nil {
		return nil, err
	}
	return tasks, nil
}
