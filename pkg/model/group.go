package model

import (
	"time"
)

type Group struct {
	ID            int       `json:"groupID"`
	OwnerID       int       `json:"ownerID"`
	Name          string    `json:"name"`
	ConfirmedTime time.Time `json:"ConfirmedTime"`
}
