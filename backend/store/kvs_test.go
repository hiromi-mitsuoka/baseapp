package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hiromi-mitsuoka/baseapp/entity"
	"github.com/hiromi-mitsuoka/baseapp/testutil"
)

func TestDVS_Save(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisTest(t)

	sut := &KVS{Cli: cli}
	key := "TestKVS_Save"
	uid := entity.UserID(1234)
	ctx := context.Background()
	// NOTE: Redisに保存したデータを削除しておく
	t.Cleanup(func() {
		cli.Del(ctx, key)
	})
	if err := sut.Save(ctx, key, uid); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}

func TestKVS_Load(t *testing.T) {
	t.Parallel()

	cli := testutil.OpenRedisTest(t)
	sut := &KVS{Cli: cli}

	// NOTE: 2つのテストケースの実行後の検証方法が異なるため，テーブル駆動テストではなく，t.Runメソッドを利用したサブテストで実装
	// https://qiita.com/marnie_ms4/items/d5233045a084cebeea14
	// サブテスト: テストに階層を作ることが可能．Runメソッドを使ってテスト内に小さなテストを記述
	t.Run("ok", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Load_ok"
		uid := entity.UserID(1234)
		ctx := context.Background()
		cli.Set(ctx, key, int64(uid), 30*time.Minute)
		// NOTE: テストする前にデータを削除しておく
		t.Cleanup(func() {
			cli.Del(ctx, key)
		})
		got, err := sut.Load(ctx, key)
		if err != nil {
			t.Fatalf("want no error, but got %v", err)
		}
		if got != uid {
			t.Errorf("want %d, but got %d", uid, got)
		}
	})

	t.Run("notFound", func(t *testing.T) {
		t.Parallel()

		key := "TestKVS_Load_notFound"
		ctx := context.Background()
		got, err := sut.Load(ctx, key)
		if err == nil || !errors.Is(err, ErrNotFound) {
			t.Errorf("want %v, but got %v(value = %d)", ErrNotFound, err, got)
		}
	})
}
