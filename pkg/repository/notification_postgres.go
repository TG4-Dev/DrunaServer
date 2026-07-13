package repository

import (
	"github.com/jmoiron/sqlx"
)

type NotificationPostgres struct {
	db *sqlx.DB
}

func NewNotificationPostgres(db *sqlx.DB) *NotificationPostgres {
	return &NotificationPostgres{db: db}
}

func (r *NotificationPostgres) Enqueue(userID int, notificationType string, payload string) error {
	query := `INSERT INTO notification_outbox (user_id, type, payload) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, userID, notificationType, payload)
	return err
}
