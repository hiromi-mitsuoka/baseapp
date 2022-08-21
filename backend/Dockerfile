# Hot-reload : https://zenn.dev/hsaki/books/golang-httpserver-internal/viewer/serverstart
FROM golang:1.18-alpine as build-reflex

# https://qiita.com/asakaguchi/items/484ba262965ef3823f61
# apk : Alpine Linuxの独自パッケージ管理システム
RUN apk update && \
    apk upgrade && \
    apk add bash git && \
    rm -rf /var/cache/apk/*

RUN go install github.com/cespare/reflex@latest


FROM golang:1.18-alpine

# https://docs.docker.jp/engine/reference/builder.html#from
# COPY --from=<名前> 命令で使用し，構築したイメージを参照できる．
COPY --from=build-reflex /go/bin/reflex /go/bin/reflex

RUN apk update && \
    apk upgrade && \
    # https://pkgs.alpinelinux.org/package/v3.3/main/x86/build-base
    # TODO: mysql使用予定
    apk add bash git build-base jq mysql-client && \
    rm -rf /var/cache/apk/*

# https://qiita.com/maemori/items/d2098fa49a1232ad5a56#goosgoarch%E4%B8%80%E8%A6%A7%E3%83%AA%E3%82%B9%E3%83%88
# https://www.getto.systems/entry/2020/06/11/013316#cgo_enabled
# https://christina04.hatenablog.com/entry/2017/08/02/041620
ENV GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

WORKDIR /app

COPY go.mod .
COPY go.sum .
# https://go.dev/ref/mod#go-mod-download
RUN go mod download

COPY . .

RUN go build -o /backend-api main.go