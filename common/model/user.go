package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserId   int64  `gorm:"not null;primarykey;autoIncrement"`
	Username string `gorm:"type:varchar(24);not null;uniqueIndex"`
	Password []byte `gorm:"type:VARBINARY(60);not null"`

	FollowingCount int64
	FollowerCount  int64

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt gorm.DeletedAt

	// 前User关注后User的多对多关系用follow表来表示; 前User使用UserId作为主键，在follow表中叫做FollowerId; 后User使用UserId作为主键，在follow表中的主键叫做FollowedId
	Following []User `gorm:"many2many:follow;foreignKey:UserId;joinForeignKey:FollowerId;References:UserId;JoinReferences:FollowedID"`
	Followers []User `gorm:"many2many:follow;foreignKey:UserId;joinForeignKey:FollowedId;References:UserId;JoinReferences:FollowerID"`
}

type Follow struct {
	//ID         int64     `gorm:"not null;primaryKey;autoIncrement"`
	FollowerId int64     `gorm:"not null;primaryKey"`
	FollowedId int64     `gorm:"not null;primaryKey"`
	CreatedAt  time.Time `gorm:"not null;"`
	//UpdatedAt  time.Time `gorm:"not null;"`

	Follower User `gorm:"foreignKey:FollowerId;References:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Followed User `gorm:"foreignKey:FollowedId;References:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
