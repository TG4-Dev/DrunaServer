package models

import (
	"time"
)

type GroupMember struct {
	GroupID       int       `gorm:"not null" json:"eventID"`
	UserID        string    `gorm:"not null" json:"userID"`
	ConfirmedTime time.Time `json:"ConfirmedTime"`
	Group         Group     `gorm:"foreignKey:GroupID"`
	User          User      `gorm:"foreignKey:UserID"`
	//PRIMARY KEY(groupID, userID)
}
