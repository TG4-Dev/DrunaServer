package models

import (
	"time"
)

type Group struct {
	ID            int       `gorm:"primary key;autoIncrement" json:"groupID"`
	OwnerID       string    `gorm:"not null" json:"ownerID"`
	Name          string    `gorm:"not null" json:"name"`
	ConfirmedTime time.Time `json:"ConfirmedTime"`
}
