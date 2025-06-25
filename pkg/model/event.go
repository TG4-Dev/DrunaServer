package model

import (
	"time"
)

type Event struct {
	ID        int       `json:"eventID"`
	UserID    int       `json:"userID"`
	StartTime time.Time `json:"startTime" binding:"required"`
	EndTime   time.Time `json:"endTime" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	Type      string    `json:"type"`
}
