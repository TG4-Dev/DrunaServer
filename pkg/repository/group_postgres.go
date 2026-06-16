package repository

import (
	"database/sql"
	"druna_server/pkg/model"
	"errors"
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
	tx, err := r.db.Beginx()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var id int
	query := fmt.Sprintf(
		"INSERT INTO %s (owner_id, name, confirmed_time) VALUES ($1, $2, $3) RETURNING id",
		groupTable,
	)
	if err := tx.QueryRow(query, input.OwnerID, input.Name, time.Now()).Scan(&id); err != nil {
		return 0, err
	}

	memberQuery := fmt.Sprintf(
		"INSERT INTO group_members (group_id, user_id, confirmed_time) VALUES ($1, $2, $3)",
	)
	if _, err := tx.Exec(memberQuery, id, input.OwnerID, time.Now()); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *GroupPostgres) ListGroups(userID int) ([]model.Group, error) {
	query := fmt.Sprintf(`
		SELECT DISTINCT g.id, g.owner_id, g.name, g.confirmed_time
		FROM %s g
		LEFT JOIN group_members gm ON g.id = gm.group_id
		WHERE g.owner_id = $1 OR gm.user_id = $1
		ORDER BY g.id`, groupTable)

	var groups []model.Group
	if err := r.db.Select(&groups, query, userID); err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *GroupPostgres) GetGroupDetails(groupID, userID int) (model.GroupDetails, error) {
	var details model.GroupDetails

	accessQuery := `
		SELECT 1 FROM groups g
		LEFT JOIN group_members gm ON g.id = gm.group_id
		WHERE g.id = $1 AND (g.owner_id = $2 OR gm.user_id = $2)
		LIMIT 1`
	var access int
	if err := r.db.Get(&access, accessQuery, groupID, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return details, errors.New("group not found or access denied")
		}
		return details, err
	}

	groupQuery := fmt.Sprintf(
		"SELECT id, owner_id, name, confirmed_time FROM %s WHERE id = $1",
		groupTable,
	)
	if err := r.db.Get(&details.Group, groupQuery, groupID); err != nil {
		return details, err
	}

	membersQuery := `
		SELECT u.id, u.name, u.username
		FROM group_members gm
		JOIN users u ON gm.user_id = u.id
		WHERE gm.group_id = $1
		ORDER BY u.username`
	if err := r.db.Select(&details.Members, membersQuery, groupID); err != nil {
		return details, err
	}

	return details, nil
}

func (r *GroupPostgres) AddGroupMember(groupID, ownerID, memberID int) error {
	var owner int
	ownerQuery := fmt.Sprintf("SELECT owner_id FROM %s WHERE id = $1", groupTable)
	if err := r.db.Get(&owner, ownerQuery, groupID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("group not found")
		}
		return err
	}
	if owner != ownerID {
		return errors.New("only group owner can add members")
	}

	query := `INSERT INTO group_members (group_id, user_id, confirmed_time) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, groupID, memberID, time.Now())
	return err
}

func (r *GroupPostgres) IsGroupMember(groupID, userID int) (bool, error) {
	var exists int
	query := `SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2 LIMIT 1`
	err := r.db.Get(&exists, query, groupID, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}
