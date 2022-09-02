package service

import (
	"context"
	"fmt"

	"github.com/hiromi-mitsuoka/baseapp/store"
)

type Login struct {
	DB store.Queryer
	// NOTE: テストを考慮しインターフェースを噛ませることで，*store.Repository型・*auth.JWTerを直接参照しない
	Repo           UserGetter
	TokenGenerator TokenGenerator
}

func (l *Login) Login(ctx context.Context, name, pw string) (string, error) {
	u, err := l.Repo.GetUser(ctx, l.DB, name)
	if err != nil {
		return "", fmt.Errorf("failed to list: %w", err)
	}
	if err = u.ComparePassword(pw); err != nil {
		return "", fmt.Errorf("wrong password: %w", err)
	}
	jwt, err := l.TokenGenerator.GenerateToken(ctx, *u)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT: %w", err)
	}
	return string(jwt), nil
}
