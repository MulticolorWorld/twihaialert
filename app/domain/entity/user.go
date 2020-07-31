package entity

import "time"

type User struct {
	ID              int
	CreatedAt       *time.Time
	LastLogin       *time.Time
	LastNotify      *time.Time
	LastNotifyCount int
	DMNotification  int `gorm:"Column:dm_notification"`
}

func NewUser() *User {
	now := time.Now()

	u := new(User)
	u.ID = 0
	u.CreatedAt = &now
	u.LastLogin = &now
	u.LastNotify = nil
	u.LastNotifyCount = 0
	u.DMNotification = 0
	return u
}

func (u User) TableName() string {
	return "user"
}
