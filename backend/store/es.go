package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/hiromi-mitsuoka/baseapp/config"
)

type ES struct {
	Cli *elasticsearch.Client
	Idx *esapi.Indices
}

func NewES(ctx context.Context, cfg *config.Config) (*ES, *esapi.Response, error) {
	escfg := elasticsearch.Config{
		Addresses: []string{
			// https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/connecting.html#connecting-without-security
			// NOTE: Use http for connecting without security enabled
			// TODO: http://localhost:9200で接続したい．現状「docker network inspect baseapp_default」or 「curl http://localhost:9201/_nodes/http\?pretty\=true」でIPアドレスを確認している
			fmt.Sprintf("http://172.29.0.5:%d", cfg.ESPort01),
			fmt.Sprintf("http://172.29.0.5:%d", cfg.ESPort02),
		},
		// https://github.com/elastic/go-elasticsearch#usage
		// TODO: To configure other HTTP settings
	}

	cli, err := elasticsearch.NewClient(escfg)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error creating the client: %s", err))
	}
	// TODO: apiコンテナが先に立ち上がってしまうため，リトライ処理かsleepを入れる
	res, err := cli.Info()
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error getting response: %s", err))
	}

	// TODO: 冗長にならないよう修正したい
	err = checkIndex(ctx, cli, "user-index", user_mapping)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error creating index: %s", err))
	}
	err = checkIndex(ctx, cli, "task-index", task_mapping)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error creating index: %s", err))
	}

	return &ES{Cli: cli}, res, nil
}

func checkIndex(ctx context.Context, cli *elasticsearch.Client, idx_name string, mapping any) error {
	// https://github.com/elastic/go-elasticsearch/blob/v8.4.0/esapi/api.indices.exists.go
	_, err := cli.Indices.Exists([]string{idx_name})
	if err != nil {
		_, err := cli.Indices.Create(idx_name)
		if err != nil {
			return err
		}
	}
	return nil
}

// mappings
// TODO: mappingを指定したい & 確認終わったら，created, modifiedを除外
var user_mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"users":{
			"properties":{
				"name":{
					"type":"keyword"
				},
				"role":{
					"type":"keyword"
				},
				"created":{
					"type":"date"
				},
				"modified":{
					"type":"date"
				}
			}
		}
	}
}`

var task_mapping = ``