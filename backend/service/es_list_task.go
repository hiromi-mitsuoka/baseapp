package service

import (
	"context"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/hiromi-mitsuoka/baseapp/store"
)

type EsListTask struct {
	ES   store.ES
	Repo EsTaskLister
}

func (elt *EsListTask) EsListTask(ctx context.Context) (*esapi.Response, error) {
	ets, err := elt.Repo.EsListTask(ctx, elt.ES)
	if err != nil {
		return nil, fmt.Errorf("failed to list from es: %w", err)
	}
	return ets, nil
}
