package store

import (
	"bytes"
	"context"
	"encoding/json"
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

// TODO: このファイルをesにクエリを投げるファイルにして，ES構築はrepository.goに移すか検討する
func NewES(ctx context.Context, cfg *config.Config) (*ES, *esapi.Response, error) {
	escfg := elasticsearch.Config{
		Addresses: []string{
			// NOTE: Use http for connecting without security enabled
			// https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/connecting.html#connecting-without-security
			// NOTE: dockerのnetworkの機能で，コンテナ名で指定可能
			// https://qiita.com/dyoshikawa/items/05d627b962da35f7d5b6
			// https://discuss.elastic.co/t/logstash-elasticsearch/206159/2
			fmt.Sprintf("http://es:%d", cfg.ESPort01),
			fmt.Sprintf("http://es:%d", cfg.ESPort02),
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

func (r *Repository) EsListTask(
	ctx context.Context, escli ES,
) (*esapi.Response, error) {
	// search
	// https://developer.okta.com/blog/2021/04/23/elasticsearch-go-developers-guide
	var buffer bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"name": "test",
			},
		},
	}
	json.NewEncoder(&buffer).Encode(query)
	response, _ := escli.Cli.Search(
		escli.Cli.Search.WithIndex("user-index"),
		escli.Cli.Search.WithBody(&buffer),
	)
	return response, nil
}
