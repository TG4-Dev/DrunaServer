package model

type User struct {
	ID           int           `json:"id"`
	Name         string        `json:"name" binding:"required"`
	Username     string        `json:"username" binding:"required"`
	Email        string        `json:"email" binding:"required"`
	PasswordHash string        `json:"passwordHash" binding:"required"`
	AvatarUrl    string        `json:"avatarURL"`
	Events       []Event       `gorm:"foreignKey:UserID"`
	OwnedGroups  []Group       `gorm:"foreignKey:OwnerID"`
	GroupMembers []GroupMember `gorm:"foreignKey:UserID"`
	Friends      []Friend      `gorm:"foreignKey:UserID"`
}
