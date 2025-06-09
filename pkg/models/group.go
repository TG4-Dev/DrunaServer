package models

type Group struct {
	ID            int    `json:"groupID"`
	OwnerID       string `json:"ownerID"`
	Name          string `json:"name"`
	ConfirmedTime string `json:confirmedTime`
}
