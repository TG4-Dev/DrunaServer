package repository

import (
	"druna_server/pkg/model"
	"errors"
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

func (r *EventPostgres) DeleteEvent(userID, eventID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", eventsTable)
	result, err := r.db.Exec(query, eventID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("event not found or you are not the owner")
	}
	return nil
}
