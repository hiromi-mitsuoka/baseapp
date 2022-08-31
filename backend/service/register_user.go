package service

import (
	"context"
	"fmt"

	"github.com/hiromi-mitsuoka/baseapp/entity"
	"github.com/hiromi-mitsuoka/baseapp/store"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUser struct {
	DB store.Execer
	// store層とのやりとり可能なメソッドを，./interface.goに記載したものに限定し，疎結合を実装する
	Repo UserRegister
}

func (r *RegisterUser) RegisterUser(
	ctx context.Context, name, password, role string,
) (*entity.User, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// https://nu50218.dev/posts/fmt-errorf-format/
		// %w で一個だけエラーを指定してあげると、Unwrap() error を実装した error を返す
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}
	u := &entity.User{
		Name:     name,
		Password: string(pw),
		Role:     role,
	}

	// NOTE: RDBMSとのやりとりは，Repository層に委譲
	if err := r.Repo.RegisterUser(ctx, r.DB, u); err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}
	return u, nil
}
