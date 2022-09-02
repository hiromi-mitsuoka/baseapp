package fixture

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/hiromi-mitsuoka/baseapp/entity"
)

// NOTE: フィクスチャ（ダミーの事前データ）を作成するテストヘルパー
//       アプリケーションのコードベースが増えることによる事前データのメンテの手間を解消するため，
//       テストヘルパーとしてダミーデータの生成関数を用意
func User(u *entity.User) *entity.User {
	result := &entity.User{
		ID:       entity.UserID(rand.Int()),
		Name:     "testuser" + strconv.Itoa(rand.Int())[:5],
		Password: "password",
		Role:     "admin",
		Created:  time.Now(),
		Modified: time.Now(),
	}
	if u == nil {
		return result
	}
	if u.ID != 0 {
		result.ID = u.ID
	}
	if u.Password != "" {
		result.Password = u.Password
	}
	if u.Role != "" {
		result.Role = u.Role
	}
	if !u.Created.IsZero() {
		result.Created = u.Created
	}
	if !u.Modified.IsZero() {
		result.Modified = u.Modified
	}
	return result
}
