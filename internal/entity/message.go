package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id" `
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeletedAt   bool               `json:"deletedAt" bson:"deletedAt"  `
	FromUserId  string             `json:"fromUserId" bson:"fromUserId"`
	ToUserId    string             `json:"toUserId" bson:"toUserId"`
	Content     string             `json:"content" bson:"content"`
	MessageType int16              `json:"messageType" bson:"messageType"`
	ContentType int16              `json:"contentType" bson:"contentType"`
	Image       string             `json:"image" bson:"image"`
	Url         string             `json:"url" bson:"url"`
}

type MessageRequest struct {
	MessageType int32  `json:"messageType"`
	Uid         string `json:"uid"`
	FriendUid   string `json:"friendUid"`
}
