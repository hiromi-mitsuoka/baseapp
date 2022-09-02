### JWT(JSON Web Token)

#### JWTとは
- JSON Web Tokenの略称であり，属性情報（Claim）をJSONデータ構造で表現したトークンの仕様
- 仕様はRFC7519
- 署名・暗号化ができ，URL-safeである
- 発音は"ジョット"

```golang
// TODO: JWTを深掘りしてまとめる
```

#### 秘密鍵作成
```terminal
openssl genrsa 4096 > secret.pem
```

### 秘密鍵から公開鍵作成
```terminal
openssl rsa -pubout < secret.pem > public.pem
```


参考記事
- [JSON Web Token（JWT）の紹介とYahoo! JAPANにおけるJWTの活用](https://techblog.yahoo.co.jp/advent-calendar-2017/jwt/)
- [GoでJWTの具体的な実装](https://christina04.hatenablog.com/entry/2017/04/15/042646)