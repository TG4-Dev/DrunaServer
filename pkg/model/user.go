package model

type User struct {
	ID           int           `json:"id"`
	Name         string        `json:"name" binding:"required"`
	Username     string        `json:"username" binding:"required"`
	Email        string        `json:"email" binding:"required"`
	PasswordHash string        `json:"passwordHash" binding:"required"`
	AvatarURL    string        `json:"avatarURL"`
	Events       []Event       `json:"events"`
	OwnedGroups  []Group       `json:"ownedGroups"`
	GroupMembers []GroupMember `json:"groupMembers"`
	Friends      []Friend      `json:"friends"`
}
