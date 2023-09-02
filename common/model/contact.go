package model

import (
	"gorm.io/gorm"
	"time"
)

type Message struct {
	MessageId int64  `gorm:"not null;primaryKey"`
	FromId    int64  `gorm:"index"`
	ToUserId  int64  `gorm:"index"`
	Content   string `gorm:"not null"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt gorm.DeletedAt
}

type Friend struct {
	gorm.Model
	UserId   int64 `gorm:"index"`
	FriendId int64 `gorm:"index"`
}
