package model

import "time"

type SignUpDoc struct {
	Name         string `json:"name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

type SignInDoc struct {
	Username     string `json:"username"`
	PasswordHash string `json:"passwordHash"`
}

type EventDoc struct {
	ID        int       `json:"eventID"`
	StartTime time.Time `json:"startTime" binding:"required"`
	EndTime   time.Time `json:"endTime" binding:"required"`
	Title     string    `json:"title" binding:"required"`
	Type      string    `json:"type"`
}

type DeleteEventDoc struct {
	ID int `json:"eventID"`
}

type AddEventDoc struct {
	ID int `json:"Id"`
}
