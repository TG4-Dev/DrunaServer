package model

import (
	"time"
)

type Event struct {
	ID        int       `json:"eventID" db:"id"`
	UserID    int       `json:"userID" db:"user_id"`
	GroupID   *int      `json:"groupID,omitempty" db:"group_id"`
	StartTime time.Time `json:"startTime" binding:"required" db:"start_time"`
	EndTime   time.Time `json:"endTime" binding:"required" db:"end_time"`
	Title     string    `json:"title" binding:"required" db:"title"`
	Type      string    `json:"type" db:"type"`
}

type EventFilter struct {
	DateFrom *time.Time
	DateTo   *time.Time
	Type     string
	Limit    int
	Offset   int
}

type EventListResponse struct {
	Events []Event `json:"events"`
	Total  int     `json:"total"`
	Limit  int     `json:"limit"`
	Offset int     `json:"offset"`
}
