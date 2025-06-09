package models

type Event struct {
	ID        int    `json:"eventID"`
	UserID    string `json:"userID"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Title     string `json:title`
	Type      string `json:type`
}
