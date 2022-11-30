package entity

import (
	"time"

	"gorm.io/gorm"
)

type GroupMember struct {
	ID int32 `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	UserId int32 `json:"userId" gorm:"index;comment:'group member user id'"`
	GroupId int32 `json:"groupId" gorm:"index;comment:'group id'"`
	Name string `json:"name" gorm:"type:varchar(350);comment:'group name'"`
	Mute bool `json:"mute" gorm:"comment:'member are muted by Owner'"`
}