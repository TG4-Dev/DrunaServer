package model

import (
	"time"
)

type Event struct {
	ID        int       `gorm:"primary key;autoIncrement" json:"eventID"`
	UserID    string    `gorm:"not null" json:"userID" binding:"required"`
	StartTime time.Time `gorm:"not null" json:"startTime" binding:"required"`
	EndTime   time.Time `gorm:"not null" json:"endTime binding:"required"`
	Title     string    `json:"title" binding:"required"`
	Type      string    `json:"type"`
}
