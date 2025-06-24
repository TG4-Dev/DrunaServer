package model

import (
	"time"
)

type GroupMember struct {
	GroupID       int       `json:"groupID" binding:"required"`
	UserID        int       `json:"userID" binding:"required"`
	ConfirmedTime time.Time `json:"confirmedTime"`
	Group         Group     `json:"group,omitempty"`
	User          User      `json:"user,omitempty"`
	// PRIMARY KEY(groupID, userID)
}
