package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/hiromi-mitsuoka/baseapp/entity"
)

func (r *Repository) RegisterUser(ctx context.Context, db Execer, u *entity.User) error {
	u.Created = r.Clocker.Now()
	u.Modified = r.Clocker.Now()
	sql := `INSERT INTO users (
						name, password, role, created, modified
					) VALUES (?, ?, ?, ?, ?)`
	result, err := db.ExecContext(ctx, sql, u.Name, u.Password, u.Role, u.Created, u.Modified)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return fmt.Errorf("cannot create same name user: %w", ErrAlreadyEntry)
		}
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = entity.UserID(id)
	return nil
}

func (r *Repository) GetUser(
	ctx context.Context, db Queryer, name string,
) (*entity.User, error) {
	u := &entity.User{}
	sql := `SELECT
						id, name, password, role, created, modified
					FROM users
					WHERE name = ?`
	// TODO: このnameの入れ方はSQLインジェクション起きる？
	// TODO: 読んでみる → https://zenn.dev/tenntenn/articles/b452911b4200ff
	if err := db.GetContext(ctx, u, sql, name); err != nil {
		return nil, err
	}
	return u, nil
}
