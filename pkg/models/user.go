package models

type User struct {
	ID           int           `gorm:"primary key;autoIncrement" json:"id"`
	Name         string        `gorm:"not null" json:"name"`
	Username     string        `gorm:"not null" json:"username"`
	Email        string        `gorm:"not null" json:"email"`
	PasswordHash string        `gorm:"not null" json:"passwordHash"`
	AvatarUrl    string        `json:"avatarURL"`
	Events       []Event       `gorm:"foreignKey:UserID"`
	OwnedGroups  []Group       `gorm:"foreignKey:OwnerID"`
	GroupMembers []GroupMember `gorm:"foreignKey:UserID"`
	Friends      []Friend      `gorm:"foreignKey:UserID"`
}
