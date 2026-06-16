package repository

import (
	"druna_server/pkg/model"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user model.User) (int, error) {
	var id int
	query := fmt.Sprintf(
		"INSERT INTO %s (name, username, email, password_hash, telegram_id) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		usersTable,
	)
	row := r.db.QueryRow(query, user.Name, user.Username, user.Email, user.PasswordHash, user.TelegramID)

	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	query := fmt.Sprintf(
		"SELECT id, name, username, email, password_hash, avatar_url, telegram_id FROM %s WHERE username=$1",
		usersTable,
	)
	err := r.db.Get(&user, query, username)
	return user, err
}

func (r *AuthPostgres) GetUserByTelegramID(telegramID int64) (model.User, error) {
	var user model.User
	query := fmt.Sprintf(
		"SELECT id, name, username, email, password_hash, avatar_url, telegram_id FROM %s WHERE telegram_id=$1",
		usersTable,
	)
	err := r.db.Get(&user, query, telegramID)
	return user, err
}

func (r *AuthPostgres) SearchUsers(prefix string) ([]model.FriendInfo, error) {
	trimmed := strings.TrimSpace(prefix)
	if trimmed == "" {
		return []model.FriendInfo{}, nil
	}

	var users []model.FriendInfo
	query := fmt.Sprintf(`
		SELECT id, name, username FROM %s
		WHERE username ILIKE $1
		ORDER BY username
		LIMIT 20`, usersTable)
	if err := r.db.Select(&users, query, trimmed+"%"); err != nil {
		return nil, err
	}
	return users, nil
}
