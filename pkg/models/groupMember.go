package models

type GroupMember struct {
	GroupID int    `json:"eventID"`
	UserID  string `json:"userID"`
	//PRIMARY KEY(groupID, userID)
	ConfirmedTime string `json:"confirmedTime"`
}
