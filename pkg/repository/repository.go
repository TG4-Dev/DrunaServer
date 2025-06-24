package repository

import (
	"druna_server/pkg/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, passwordHash string) (model.User, error)
}

type User interface {
}

type Event interface {
	CreateEvent(user model.Event) (int, error)
	DeleteEvent(userID, eventID int) error
}

type Friendship interface {
	CreateFriendRequest(userID, friendID int) error
	AcceptFriendRequest(userID int, friendID int) error
	ExistsByUsername(username string) (int, error)
	GetFriendList(userID int) ([]model.FriendInfo, error)
	GetFriendRequestList(userID int) ([]model.FriendInfo, error)
	DeleteFriend(userID int, friendID int) error
}

type Group interface {
}

type Repository struct {
	Authorization
	User
	Event
	Friendship
	Group
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Event:         NewEventPostgres(db),
		Friendship:    NewFriendShipPostgres(db),
	}
}
