package store

import (
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/hiromi-mitsuoka/baseapp/config"
)

type ES struct {
	Cli *elasticsearch.Client
}

func NewES(cfg *config.Config) (*ES, *esapi.Response, error) {
	escfg := elasticsearch.Config{
		Addresses: []string{
			// https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/connecting.html#connecting-without-security
			// NOTE: Use http for connecting without security enabled
			// TODO: http://localhost:9200で接続したい．現状「docker network inspect baseapp_default」or 「curl http://localhost:9201/_nodes/http\?pretty\=true」でIPアドレスを確認している
			fmt.Sprintf("http://172.29.0.5:%d", cfg.ESPort01),
			fmt.Sprintf("http://172.29.0.5:%d", cfg.ESPort02),
		},
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
	return &ES{Cli: cli}, res, nil
}
