.PHONY: test
test:
# https://qiita.com/tamanobi/items/6aca6e3cb2879f43f3ad#race-condition%E3%82%92%E7%99%BA%E8%A6%8B%E3%81%99%E3%82%8B--race-%E3%82%AA%E3%83%97%E3%82%B7%E3%83%A7%E3%83%B3
# -race : race conditionを発見する
# https://future-architect.github.io/articles/20210812a/
# shulle : テストとベンチマークの実行順序を制御する．順番をランダムにすることで、前に書いたテストの実行結果に依存する状態を検知し、回避できることが最大の目的
	go test -race -shuffle=on

.PHONY: generate
# https://speakerdeck.com/yaegashi/go-generate-wan-quan-ru-men
# https://qiita.com/yaegashi/items/d1fd9f7d0c75b2bb7446
generate:
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run