package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

type TokenPostgres struct {
	db *sqlx.DB
}

func NewTokenPostgres(db *sqlx.DB) *TokenPostgres {
	return &TokenPostgres{db: db}
}

func (r *TokenPostgres) RevokeToken(jti string, expiresAt time.Time) error {
	query := `INSERT INTO revoked_tokens (jti, revoked_at, expires_at) VALUES ($1, $2, $3)
		ON CONFLICT (jti) DO NOTHING`
	_, err := r.db.Exec(query, jti, time.Now(), expiresAt)
	return err
}

func (r *TokenPostgres) IsTokenRevoked(jti string) (bool, error) {
	var exists int
	query := `SELECT 1 FROM revoked_tokens WHERE jti = $1 LIMIT 1`
	err := r.db.Get(&exists, query, jti)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *TokenPostgres) Ping() error {
	return r.db.Ping()
}
