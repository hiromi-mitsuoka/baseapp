package store

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/hiromi-mitsuoka/baseapp/clock"
	"github.com/hiromi-mitsuoka/baseapp/entity"
	"github.com/hiromi-mitsuoka/baseapp/testutil"
	"github.com/hiromi-mitsuoka/baseapp/testutil/fixture"
	"github.com/jmoiron/sqlx"
)

// テストヘルパー関数
func prepareUser(ctx context.Context, t *testing.T, db Execer) entity.UserID {
	t.Helper()
	u := fixture.User(nil)
	result, err := db.ExecContext(
		ctx,
		`INSERT INTO users (name, password, role, created, modified)
		 VALUES (?, ?, ?, ?, ?);`,
		u.Name, u.Password, u.Role, u.Created, u.Modified,
	)
	if err != nil {
		t.Fatalf("insert user: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("got user_id: %v", err)
	}
	return entity.UserID(id)
}

func prepareTasks(ctx context.Context, t *testing.T, con Execer) (entity.UserID, entity.Tasks) {
	t.Helper()

	// 一度綺麗にする
	if _, err := con.ExecContext(ctx, "DELETE FROM tasks;"); err != nil {
		t.Logf("failed to initialize tasks: %v", err)
	}

	userID := prepareUser(ctx, t, con)
	otherUserID := prepareUser(ctx, t, con)
	c := clock.FixedClocker{}
	wants := entity.Tasks{
		{
			UserID:   userID,
			Title:    "want task 1",
			Status:   "todo",
			Created:  c.Now(),
			Modified: c.Now(),
		},
		{
			UserID:   userID,
			Title:    "want task 2",
			Status:   "done",
			Created:  c.Now(),
			Modified: c.Now(),
		},
	}
	tasks := entity.Tasks{
		wants[0],
		{
			UserID:   otherUserID,
			Title:    "not want task",
			Status:   "todo",
			Created:  c.Now(),
			Modified: c.Now(),
		},
		wants[1],
	}

	result, err := con.ExecContext(
		ctx,
		`INSERT INTO tasks
			(user_id, title, status, created, modified)
		VALUES
			(?, ?, ?, ? ,?),
			(?, ?, ?, ? ,?),
			(?, ?, ?, ? ,?);`,
		wants[0].UserID, wants[0].Title, wants[0].Status, wants[0].Created, wants[0].Modified,
		wants[0].UserID, wants[1].Title, wants[1].Status, wants[1].Created, wants[1].Modified,
		wants[0].UserID, wants[2].Title, wants[2].Status, wants[2].Created, wants[2].Modified,
	)
	if err != nil {
		t.Fatal(err)
	}
	// NOTE: 複数INSERTの場合，LstInsertIdで取得するのは最初のレコードID
	id, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	tasks[0].ID = entity.TaskID(id)
	tasks[1].ID = entity.TaskID(id + 1)
	tasks[2].ID = entity.TaskID(id + 2)
	return userID, wants
}

// NOTE: モックを利用したテスト
func TestRepository_AddTask(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	c := clock.FixedClocker{}
	var wantID int64 = 20
	okTask := &entity.Task{
		Title:    "ok task",
		Status:   "todo",
		Created:  c.Now(),
		Modified: c.Now(),
	}

	// github.com/DATA-DOG/go-sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectExec(
		// NOTE: エスケープが必要
		`INSERT INTO tasks \(title, status, created, modified\) VALUES \(\?, \?, \?, \?\)`,
	).WithArgs(okTask.Title, okTask.Status, c.Now(), c.Now()).
		WillReturnResult(sqlmock.NewResult(wantID, 1))

	// NOTE: mock作成時に初期化したdbを引数にして，xdbを初期化
	xdb := sqlx.NewDb(db, "mysql")
	r := &Repository{Clocker: c}
	// TODO: mock利用したdbを使って初期化したxdbを利用しているから，RDBMSは用いていない？？
	if err := r.AddTask(ctx, xdb, okTask); err != nil {
		t.Errorf("want no error, but got %v", err)
	}
}

// NOTE: RDBMSを利用したテスト
func TestRepository_ListTasks(t *testing.T) {
	ctx := context.Background()
	// entity.Taskを作成する他のテストケースと混ざるとテスト結果が異なりフェイルする
	// そのため，トランザクションを貼ることで，このテストケースの中だけのテーブル状態にする
	tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
	// このテストケースが完了したら元に戻す
	t.Cleanup(func() { _ = tx.Rollback() })
	if err != nil {
		t.Fatal(err)
	}

	wantUserID, wants := prepareTasks(ctx, t, tx)

	sut := &Repository{}
	gots, err := sut.ListTasks(ctx, tx, wantUserID)
	if err != nil {
		t.Fatalf("unexected error: %v", err)
	}
	if d := cmp.Diff(gots, wants); len(d) != 0 {
		t.Errorf("differs: (-got +want)\n%s", d)
	}
}
