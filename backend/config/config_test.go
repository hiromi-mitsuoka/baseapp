package config

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	wantPort := 3333
	// fmt.Sprint()でstring型に
	t.Setenv("PORT", fmt.Sprint(wantPort))

	got, err := New()
	if err != nil {
		// https://simple-minds-think-alike.moritamorie.com/entry/go-testing-error-fatal
		// t.Fatalf() : FailNow→ 対象の関数のテストに失敗した記録を残し、後続のテストは実行しない。
		t.Fatalf("cannot create config: %v", err)
	}
	if got.Port != wantPort {
		// https://simple-minds-think-alike.moritamorie.com/entry/go-testing-error-fatal
		// t.Errorf() : Fail→ 対象の関数のテストに失敗した記録を残すが、後続のテストは実行する。
		t.Errorf("want %d, but %d", wantPort, got.Port)
	}
	wantEnv := "dev"
	if got.Env != wantEnv {
		t.Errorf("want %s, but %s", wantEnv, got.Env)
	}
}
