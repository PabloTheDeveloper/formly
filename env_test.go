package ksat

import (
	"testing"
)

func TestNewEnv(t *testing.T) {
	env, err := NewEnv(LocalSqlite)
	if err != nil {
		t.Fatal(err)
	}
	if env == nil {
		t.Fatal("supposed to have an env variable")
	}
	if _, err := NewEnv(1); err == nil {
		t.Fatal("supposed to fail")
	}
}
