package models

type Friend struct {
	UserID      int    `json:"userID"`
	FrinedID    string `json:"friendID"`
	Status      string `json:"status"`
	RequestAt   string `json:requestAt`
	confirmedAt string `json:"confirmedAt"`
	//PRIMARY KEY(user_id, friend_id)
}
