package repository

import (
	"database/sql"
	"druna_server/pkg/model"
	"errors"
	"fmt"
	"time"

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

func (r *EventPostgres) UpdateEvent(userID int, event model.Event) error {
	query := fmt.Sprintf(`
		UPDATE %s SET start_time = $1, end_time = $2, title = $3, type = $4
		WHERE id = $5 AND user_id = $6`, eventsTable)
	result, err := r.db.Exec(query, event.StartTime, event.EndTime, event.Title, event.Type, event.ID, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("event not found or you are not the owner")
	}
	return nil
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

func (r *EventPostgres) HasOverlappingEvent(userID int, start, end time.Time, excludeID int) (bool, error) {
	query := fmt.Sprintf(`
		SELECT 1 FROM %s
		WHERE user_id = $1
		  AND id <> $2
		  AND start_time < $4
		  AND end_time > $3
		LIMIT 1`, eventsTable)
	var exists int
	err := r.db.Get(&exists, query, userID, excludeID, start, end)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *EventPostgres) GetEventList(userID int) ([]model.Event, error) {
	return r.GetEventListFiltered(userID, model.EventFilter{})
}

func (r *EventPostgres) GetEventListFiltered(userID int, filter model.EventFilter) ([]model.Event, error) {
	query := fmt.Sprintf("SELECT id, user_id, title, start_time, end_time, type FROM %s WHERE user_id = $1", eventsTable)
	args := []interface{}{userID}
	argIdx := 2

	if filter.DateFrom != nil {
		query += fmt.Sprintf(" AND end_time >= $%d", argIdx)
		args = append(args, *filter.DateFrom)
		argIdx++
	}
	if filter.DateTo != nil {
		query += fmt.Sprintf(" AND start_time <= $%d", argIdx)
		args = append(args, *filter.DateTo)
		argIdx++
	}
	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, filter.Type)
		argIdx++
	}

	query += " ORDER BY start_time ASC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, filter.Limit)
		argIdx++
	}
	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIdx)
		args = append(args, filter.Offset)
	}

	var events []model.Event
	if err := r.db.Select(&events, query, args...); err != nil {
		return nil, err
	}
	return events, nil
}

func (r *EventPostgres) CountEvents(userID int, filter model.EventFilter) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE user_id = $1", eventsTable)
	args := []interface{}{userID}
	argIdx := 2

	if filter.DateFrom != nil {
		query += fmt.Sprintf(" AND end_time >= $%d", argIdx)
		args = append(args, *filter.DateFrom)
		argIdx++
	}
	if filter.DateTo != nil {
		query += fmt.Sprintf(" AND start_time <= $%d", argIdx)
		args = append(args, *filter.DateTo)
		argIdx++
	}
	if filter.Type != "" {
		query += fmt.Sprintf(" AND type = $%d", argIdx)
		args = append(args, filter.Type)
	}

	var count int
	if err := r.db.Get(&count, query, args...); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *EventPostgres) GetEventsForUsers(userIDs []int, dateFrom, dateTo time.Time) (map[int][]model.Event, error) {
	if len(userIDs) == 0 {
		return map[int][]model.Event{}, nil
	}

	query, args, err := sqlx.In(fmt.Sprintf(`
		SELECT id, user_id, title, start_time, end_time, type FROM %s
		WHERE user_id IN (?)
		  AND end_time >= ?
		  AND start_time <= ?
		ORDER BY user_id, start_time`, eventsTable), userIDs, dateFrom, dateTo)
	if err != nil {
		return nil, err
	}
	query = r.db.Rebind(query)

	var events []model.Event
	if err := r.db.Select(&events, query, args...); err != nil {
		return nil, err
	}

	result := make(map[int][]model.Event)
	for _, event := range events {
		result[event.UserID] = append(result[event.UserID], event)
	}
	return result, nil
}
