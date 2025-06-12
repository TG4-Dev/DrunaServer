package models

import (
	"time"
)

type Friend struct {
	UserID      int       `gorm:"not null" json:"userID"`
	FriendID    int       `gorm:"not null" json:"friendID"`
	User        User      `gorm:"foreignKey:UserID"`
	Friend      User      `gorm:"foreignKey:FriendID"`
	Status      string    `gorm:"not null" json:"status"`
	RequestAt   time.Time `gorm:"not null" json:"requestAt"`
	ConfirmedAt time.Time `json:"confirmedAt"`
	//PRIMARY KEY(user_id, friend_id)
}
