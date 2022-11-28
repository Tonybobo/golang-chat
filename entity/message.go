package entity

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID int32 `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	FromUserId int32 `json:"fromUserId" gorm:"index"`
	ToUserId int32 `json:"toUserId" gorm:"index;comment:'one or group id'"`
	Content string `json:"content" gorm:"type:varchar(2500)"`
	MessageType int16 `json:"messageType" gorm:"comment:'one or group message'"`
	ContentType int16 `json:"contentType" gorm:"comment:'text , file , pic , audio , video ...'"`
	Image string `json:"image" gorm:"type:text; comment:'image'" `
	Url string `json:"url" gorm:"type:varchar(250);comment:'file url'"`
}