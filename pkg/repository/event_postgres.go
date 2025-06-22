package repository

import (
	"druna_server/pkg/model"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type EventPostgres struct {
	db *sqlx.DB
}

func NewEventPostgres(db *sqlx.DB) *EventPostgres {
	return &EventPostgres{db: db}
}

func (r *EventPostgres) CreateEvent(event model.Event) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_id, start_time, end_time, title, type) VALUES ($1, $2, $3, $4, $5) RETURNING id", eventsTable)

	row := r.db.QueryRow(query,
		event.UserID,
		event.StartTime,
		event.EndTime,
		event.Title,
		event.Type)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
