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

type FriendRequest struct {
	Uid string `json:"uid"`
	FriendUid string  `json:"friend_uid"`
}

func (f *UserFriend) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("CreatedAt", time.Now())
	return nil
}

func (f *UserFriend) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}