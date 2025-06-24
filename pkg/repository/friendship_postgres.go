package repository

import (
	"druna_server/pkg/model"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type FriendShipPostgres struct {
	db *sqlx.DB
}

func NewFriendShipPostgres(db *sqlx.DB) *FriendShipPostgres {
	return &FriendShipPostgres{db: db}
}

func (r *FriendShipPostgres) ExistsByUsername(username string) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE username = $1", usersTable)

	row := r.db.QueryRow(query, username)
	if err := row.Scan(&id); err != nil {
		fmt.Println(err)
		return 0, err
	}
	return id, nil
}

func (r *FriendShipPostgres) CreateFriendRequest(userID int, friendID int) error {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_id, friend_id, status, request_at, confirmed_at)VALUES ($1, $2, $3, $4, NULL) RETURNING user_id", friendsTable)

	row := r.db.QueryRow(query,
		userID,
		friendID,
		"pending",
		time.Now())
	if err := row.Scan(&id); err != nil {
		return err
	}
	return nil
}

func (r *FriendShipPostgres) AcceptFriendRequest(userID int, friendID int) error {
	var id int
	query := fmt.Sprintf("UPDATE %s SET status = 'accepted' WHERE friend_id = $1 AND user_id = $2 RETURNING friend_id", friendsTable)

	row := r.db.QueryRow(query,
		userID,
		friendID)
	if err := row.Scan(&id); err != nil {
		return err
	}
	return nil
}

func (r *FriendShipPostgres) RejectFriendRequest(userID int, friendID int) error {
	var id int
	query := fmt.Sprintf("UPDATE %s SET status = 'rejected' WHERE friend_id = $1 AND user_id = $2 RETURNING friend_id", friendsTable)

	row := r.db.QueryRow(query,
		userID,
		friendID)
	if err := row.Scan(&id); err != nil {
		return err
	}
	return nil
}

func (r *FriendShipPostgres) GetFriendList(userID int) ([]model.FriendInfo, error) {
	var friends []model.FriendInfo

	query := `
        SELECT u.id, u.name, u.username 
        FROM friends f
        JOIN users u ON 
            (f.friend_id = u.id AND f.user_id = $1 AND f.status = 'accepted') OR
            (f.user_id = u.id AND f.friend_id = $1 AND f.status = 'accepted')`

	err := r.db.Select(&friends, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get friend list: %w", err)
	}

	return friends, nil
}

func (r *FriendShipPostgres) GetFriendRequestList(userID int) ([]model.FriendInfo, error) {
	var friends []model.FriendInfo

	query := `
        SELECT u.id, u.name, u.username 
        FROM friends f
        JOIN users u ON 
            (f.friend_id = u.id AND f.user_id = $1 AND f.status = 'pending') OR
            (f.user_id = u.id AND f.friend_id = $1 AND f.status = 'pending')`

	err := r.db.Select(&friends, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get friend list: %w", err)
	}

	return friends, nil
}

func (r *FriendShipPostgres) DeleteFriend(userID int, friendID int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE (user_id = $1 AND friend_id = $2) OR (user_id = $2 AND friend_id = $1)", friendsTable)

	result, err := r.db.Exec(query, userID, friendID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no friendship found to delete")
	}

	return nil
}
