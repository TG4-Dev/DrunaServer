package service

import (
	"druna_server/pkg/model"
	"os"
	"testing"
	"time"
)

type mockAuthRepo struct {
	users map[string]model.User
}

func (m *mockAuthRepo) CreateUser(user model.User) (int, error) {
	user.ID = len(m.users) + 1
	m.users[user.Username] = user
	return user.ID, nil
}

func (m *mockAuthRepo) GetUserByUsername(username string) (model.User, error) {
	user, ok := m.users[username]
	if !ok {
		return model.User{}, os.ErrNotExist
	}
	return user, nil
}

func (m *mockAuthRepo) GetUserByTelegramID(telegramID int64) (model.User, error) {
	return model.User{}, os.ErrNotExist
}

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-tests-only")
	os.Exit(m.Run())
}

func TestHashAndVerifyPassword(t *testing.T) {
	hash, err := hashPassword("secret123")
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}
	if hash == "secret123" {
		t.Fatal("password must be hashed")
	}
}

func TestGenerateAndParseToken(t *testing.T) {
	svc := NewAuthService(&mockAuthRepo{})
	user := model.User{ID: 42, Username: "alice"}

	token, err := svc.GenerateToken(time.Hour, user)
	if err != nil {
		t.Fatalf("generate token failed: %v", err)
	}

	userID, username, err := svc.ParseToken(token)
	if err != nil {
		t.Fatalf("parse token failed: %v", err)
	}
	if userID != 42 || username != "alice" {
		t.Fatalf("unexpected claims: %d %s", userID, username)
	}
}

func TestGenerateAccessRefreshToken(t *testing.T) {
	hash, err := hashPassword("password")
	if err != nil {
		t.Fatalf("hash failed: %v", err)
	}

	repo := &mockAuthRepo{
		users: map[string]model.User{
			"bob": {ID: 1, Username: "bob", PasswordHash: hash},
		},
	}
	svc := NewAuthService(repo)

	access, refresh, err := svc.GenerateAccessRefreshToken("bob", "password")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if access == "" || refresh == "" {
		t.Fatal("expected non-empty tokens")
	}
}
