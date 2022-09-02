# backend

## エンドポイント

## 動作確認コマンド

### ユーザー登録してアクセストークン発行

管理者
```terminal
curl -X POST localhost:18000/register -d '{"name": "admin_user", "password": "test", "role": "admin"}'
```

一般ユーザー
```terminal
curl -X POST localhost:18000/register -d '{"name": "normal_user", "password": "testtest", "role": "user"}'
```

### ログイン
```terminal
curl -XPOST localhost:18000/login -d '{"user_name": "admin_user", "password": "test"}'

// {"access_token":"eyJh......................
// ユーザー毎にaccess_tokenが異なる
```

### タスク登録

```terminal
export TOKEN=eyJh......................
curl -XPOST -H "Authorization: Bearer $TOKEN" localhost:18000/tasks -d @./handler/testdata/add_task/ok_req.json.golden
```

### ユーザー自身のタスク取得
```terminal
export TOKEN=eyJh......................
curl -XGET -H "Authorization: Bearer $TOKEN" localhost:18000/tasks | jq
```

### 管理者アクセス
```terminal
curl -XGET -H "Authorization: Bearer $TOKEN" localhost:18000/admin

// {"message": "admin only"}
```



