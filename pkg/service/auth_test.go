package service

import (
	"druna_server/pkg/model"
	"errors"
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

func (m *mockAuthRepo) GetUserByID(userID int) (model.User, error) {
	for _, user := range m.users {
		if user.ID == userID {
			return user, nil
		}
	}
	return model.User{}, os.ErrNotExist
}

func (m *mockAuthRepo) UpdateUserProfile(userID int, name, avatarURL string) error {
	for username, user := range m.users {
		if user.ID == userID {
			if name != "" {
				user.Name = name
			}
			if avatarURL != "" {
				user.AvatarURL = avatarURL
			}
			m.users[username] = user
			return nil
		}
	}
	return os.ErrNotExist
}

func (m *mockAuthRepo) SearchUsers(prefix string) ([]model.FriendInfo, error) {
	return []model.FriendInfo{}, nil
}

type mockTokenRepo struct {
	revoked map[string]bool
}

func (m *mockTokenRepo) RevokeToken(jti string, expiresAt time.Time) error {
	if m.revoked == nil {
		m.revoked = map[string]bool{}
	}
	m.revoked[jti] = true
	return nil
}

func (m *mockTokenRepo) IsTokenRevoked(jti string) (bool, error) {
	return m.revoked != nil && m.revoked[jti], nil
}

func (m *mockTokenRepo) Ping() error { return nil }

func (m *mockTokenRepo) PurgeExpiredTokens() (int64, error) { return 0, nil }

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-tests-only")
	os.Exit(m.Run())
}

func TestCreateUserPasswordValidation(t *testing.T) {
	svc := NewAuthService(&mockAuthRepo{users: map[string]model.User{}}, &mockTokenRepo{})
	_, err := svc.CreateUser(model.User{
		Name:     "Short",
		Username: "short",
		Email:    "short@test.local",
		Password: "abc",
	})
	if !errors.Is(err, ErrPasswordTooShort) {
		t.Fatalf("expected ErrPasswordTooShort, got %v", err)
	}
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

func TestGenerateAndParseAccessToken(t *testing.T) {
	hash, _ := hashPassword("password")
	repo := &mockAuthRepo{
		users: map[string]model.User{
			"alice": {ID: 42, Username: "alice", PasswordHash: hash},
		},
	}
	svc := NewAuthService(repo, &mockTokenRepo{})

	access, refresh, err := svc.GenerateAccessRefreshToken("alice", "password")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	userID, username, err := svc.ParseAccessToken(access)
	if err != nil {
		t.Fatalf("parse access token failed: %v", err)
	}
	if userID != 42 || username != "alice" {
		t.Fatalf("unexpected claims: %d %s", userID, username)
	}

	_, _, err = svc.ParseAccessToken(refresh)
	if err == nil {
		t.Fatal("refresh token must not parse as access token")
	}
}

func TestRenewTokenRotatesRefresh(t *testing.T) {
	hash, _ := hashPassword("password")
	repo := &mockAuthRepo{
		users: map[string]model.User{
			"bob": {ID: 1, Username: "bob", PasswordHash: hash},
		},
	}
	tokenRepo := &mockTokenRepo{}
	svc := NewAuthService(repo, tokenRepo)

	_, refresh, err := svc.GenerateAccessRefreshToken("bob", "password")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	newAccess, newRefresh, err := svc.RenewToken(refresh)
	if err != nil {
		t.Fatalf("renew failed: %v", err)
	}
	if newAccess == "" || newRefresh == "" {
		t.Fatal("expected new tokens")
	}

	_, _, err = svc.RenewToken(refresh)
	if err == nil {
		t.Fatal("old refresh token should be revoked")
	}
}
