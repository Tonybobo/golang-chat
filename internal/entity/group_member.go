package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupMember struct {
	ID  primitive.ObjectID          `json:"_id" bson:"_id" `
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
	DeletedAt bool `json:"deletedAt" bson:"deletedAt"  `
	UserId string `json:"userId" bson:"userId"`
	GroupId string `json:"groupId" bson:"userId"`
	Name string `json:"name" bson:"name"`
	Mute bool `json:"mute" bson:"mute"`
}

