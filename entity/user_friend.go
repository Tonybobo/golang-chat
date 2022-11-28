package entity

import (
	"time"

	"gorm.io/gorm"
)

type UserFriend struct {
	ID int32 `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt      `json:"deletedAt"`
	UserId int32 `json:"userId" gorm:"index;comment:'userId'"`
	FriendId int32 `json:"friendId" gorm:"index;comment:'friendId'"`
}