package entity

import "time"

type User struct {
	ID              int        `gorm:"primary_key"`
	CreatedAt       *time.Time `gorm:"type:datetime;not null"`
	LastLogin       *time.Time `gorm:"type:datetime;not null"`
	LastNotify      *time.Time `gorm:"type:datetime"`
	LastNotifyCount int        `gorm:"not null"`
	DMNotification  int        `gorm:"Column:dm_notification;not null"`
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
