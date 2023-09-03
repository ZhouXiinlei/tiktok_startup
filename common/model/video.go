package model

import (
	"gorm.io/gorm"
	"time"
)

type Video struct {
	VideoId       int64  `gorm:"not null;primaryKey"`
	AuthorId      int64  `gorm:"not null;index"`
	Title         string `gorm:"not null;index"`
	PlayUrl       string `gorm:"not null"`
	CoverUrl      string `gorm:"not null"`
	FavoriteCount int64
	CommentCount  int64

	// has many not used
	//Comments  []Comment  `gorm:"foreignKey:VideoId;References:VideoId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	//Favorites []Favorite `gorm:"foreignKey:VideoId;References:VideoId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Author User `gorm:"foreignKey:AuthorId;References:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt gorm.DeletedAt
}

type Comment struct {
	CommentId int64  `gorm:"not null;primaryKey"`
	UserId    int64  `gorm:"not null;index"`
	VideoId   int64  `gorm:"not null;index"`
	Content   string `gorm:"not null"`

	Video Video `gorm:"joinForeignKey:VideoId;joinReferences:VideoId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	User  User  `gorm:"joinForeignKey:UserId;joinReferences:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt gorm.DeletedAt
}

type Favorite struct {
	FavoriteId int64 `gorm:"not null;index;primaryKey"`
	UserId     int64 `gorm:"not null;index"`
	VideoId    int64 `gorm:"not null;index"`

	User  User  `gorm:"joinForeignKey:UserId;joinReferences:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Video Video `gorm:"joinForeignKey:VideoId;joinReferences:VideoId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	DeletedAt gorm.DeletedAt
}
