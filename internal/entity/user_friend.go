package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserFriend struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id" `
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeletedAt bool               `json:"deletedAt" bson:"deletedAt"  `
	UserId    string             `json:"userId" bson:"userId" validate:"required"`
	FriendId  string             `json:"friendId" bson:"friendId" validate:"required"`
}

type FriendRequest struct {
	Uid       string `json:"uid" bson:"uid"`
	FriendUid string `json:"friend_uid" bson:"friend_uid" `
}
