package repository

import (
	"druna_server/pkg/model"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type GroupPostgres struct {
	db *sqlx.DB
}

func NewGroupPostgres(db *sqlx.DB) *GroupPostgres {
	return &GroupPostgres{db: db}
}

func (r *GroupPostgres) CreateGroup(input model.Group) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (owner_id, name, confirmed_time) VALUES ($1, $2, $3) RETURNING id", groupTable)
	row := r.db.QueryRow(query,
		input.OwnerID,
		input.Name,
		time.Now())
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return 1, nil
}
