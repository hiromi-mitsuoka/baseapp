//go:build tools

package main

import _ "github.com/matryer/moq"

// NOTE: matryer/moq は，golang/mock に比べ，モックの振る舞いを指定する際に型を意識して実装が可能
// 	     golang/mock は，振る舞いを設定する際に引数がanyのセッターを利用するため

// 参考記事（go generate）
// https://speakerdeck.com/yaegashi/go-generate-wan-quan-ru-men
// https://qiita.com/yaegashi/items/d1fd9f7d0c75b2bb7446

// NOTE: このファイルを定義しておくことで，go.modによるバージョン管理ができる
//       → ビルドタグを指定しない実アプリケーションのビルド時には無視される
//       （go run コマンドを使用すると，「常に実行タイミングで最新のバージョンのプログラムが実行されてしまう」心配があるため，上記を利用）
