# エンドポイント

| HTTPメソッド | パス | 概要 | アクセストークンの有無 |
| :--- | :--- | :--- | :--- |
| POST | /register | 新しいユーザーを登録 | No |
| POST | /login | 登録済みユーザー情報でアクセストークンを取得 | No |
| POST | /tasks | タスク登録 | Yes |
| GET | /tasks | トークンに紐づくユーザーのタスク取得 | Yes |
| PUT | /tasks/{taskID} | トークンに紐づくユーザーのタスク更新 | Yes |
| DELETE | /tasks/{taskID} | トークンに紐づくユーザーのタスク削除 | Yes |
|  |  |  |  |
| GET | /admin | 管理者ユーザーのみアクセス | Yes |

# 動作確認コマンド

## api （backend-nginxの81番経由）

### health-check
```terminal
curl localhost:81/health

// {"status": "OK"}
```

### ユーザー登録してアクセストークン発行

管理者
```terminal
curl -X POST localhost:81/register -d '{"name": "admin_user", "password": "test", "role": "admin"}'
```

一般ユーザー
```terminal
curl -X POST localhost:81/register -d '{"name": "normal_user", "password": "testtest", "role": "user"}'
```

### ログイン
```terminal
curl -XPOST localhost:81/login -d '{"user_name": "admin_user", "password": "test"}'

// {"access_token":"eyJh......................
// ユーザー毎・発行毎にaccess_tokenが異なる
```

### タスク登録

```terminal
export TOKEN=eyJh......................
curl -XPOST -H "Authorization: Bearer $TOKEN" localhost:81/tasks -d @./handler/testdata/add_task/ok_req.json.golden
```

### ユーザー自身のタスク取得
```terminal
export TOKEN=eyJh......................
curl -XGET -H "Authorization: Bearer $TOKEN" localhost:81/tasks | jq
```

### タスクの編集（id=1）
```terminal
export TOKEN=eyJh......................
curl -XPUT -H "Authorization: Bearer $TOKEN" localhost:81/tasks/1 -d @./handler/testdata/update_task/ok_req.json.golden
```

### タスクの削除（id=1）
```terminal
export TOKEN=eyJh......................
curl -XDELETE -H "Authorization: Bearer $TOKEN" localhost:81/tasks/1
```

### 管理者アクセス
```terminal
curl -XGET -H "Authorization: Bearer $TOKEN" localhost:81/admin

// {"message": "admin only"}
```


## ElasticSearch


### 基本情報取得
```terminal
curl -XGET "http://localhost:9201"
```

### indexの一覧取得
```terminal
curl -XGET "http://localhost:9201/_aliases" | jq
```

### indexのドキュメント数を取得
```terminal
curl -XGET "http://localhost:9201/_cat/count/<index_name>"
```

### mappingの確認
```terminal
curl -XGET "http://localhost:9201/<index-name>/_mapping?pretty"
```

**参考記事**
[【Elasticsearch】よく使うコマンド一覧](https://qiita.com/mug-cup/items/ba5dd0a14838e83e69ac)

## docker

### docker環境内のnetwork確認
```terminal
docker network inspect baseapp_default
```