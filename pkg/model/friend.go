package model

import (
	"time"
)

type Friend struct {
	UserID      int       `json:"userID"`
	FriendID    int       `json:"friendID"`
	Status      string    `json:"status"`
	RequestAt   time.Time `json:"requestAt"`
	ConfirmedAt time.Time `json:"confirmedAt"`
	// User        User      `json:"user"`
	// Friend      User      `json:"friend"`
	//PRIMARY KEY(user_id, friend_id)
}

type FriendInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}
