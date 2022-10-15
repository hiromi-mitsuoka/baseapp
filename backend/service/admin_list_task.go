package service

import (
	"context"
	"fmt"

	"github.com/hiromi-mitsuoka/baseapp/entity"
	"github.com/hiromi-mitsuoka/baseapp/store"
)

type AdminListTask struct {
	DB   store.Queryer
	Repo AdminTaskLister
}

func (alt *AdminListTask) AdminListTask(ctx context.Context) (entity.Tasks, error) {
	ats, err := alt.Repo.AdminListTask(ctx, alt.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to list: %w", err)
	}
	return ats, nil
}
