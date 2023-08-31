package model

import "time"

type User struct {
	UserId    int64     `gorm:"not null;primarykey;autoIncrement"`
	Username  string    `gorm:"type:varchar(24);not null;uniqueIndex"`
	Password  []byte    `gorm:"type:VARBINARY(60);not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
