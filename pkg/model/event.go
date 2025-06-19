package model

import (
	"time"
)

type Event struct {
	ID        int       `gorm:"primary key;autoIncrement" json:"eventID"`
	UserID    string    `gorm:"not null" json:"userID"`
	User      User      `gorm:"foreignKey:UserID"`
	StartTime time.Time `gorm:"not null" json:"startTime"`
	EndTime   time.Time `gorm:"not null" json:"endTime"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
}
