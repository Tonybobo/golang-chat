package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupChat struct {
	ID         primitive.ObjectID          `json:"_id" bson:"_id" `
	Uid       string         `json:"uid" bson:"uid" validate:"required"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
	DeletedAt bool `json:"deletedAt" bson:"deletedAt"  `

	UserId    string          `json:"userId" bson:"userId"`
	Name      string         `json:"name" bson:"name"`
	Notice    string         `json:"notice" bson:"notice"`
	Avatar    string             `json:"avatar" bson:"avatar"`
}

type GroupResponse struct {
	Uid       string         `json:"uid" bson:"uid" `
	GroupId		string `json:"groupId" bson:"groupId"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Name 	string	`json:"name" bson:"name"`
	Notice 	string `json:"notice" bson:"name"`
}


