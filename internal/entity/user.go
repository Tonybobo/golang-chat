package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id" `
	Uid       string             `json:"uid" bson:"uid" `
	Username  string             `json:"username" bson:"username" validate:"required, min=4"`
	Name      string             `json:"name" bson:"name"`
	Password  string             `json:"password" bson:"password" validate:"required,min=4"`
	Avatar    string             `json:"avatar" bson:"avatar"`
	Email     string             `json:"email" bson:"email" validate:"required"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeletedAt bool               `json:"deletedAt" bson:"deletedAt"  `
}

type UserResponse struct {
	Uid      string `json:"uid" bson:"uid"`
	Username string `json:"username" bson:"username" validate:"required"`
	Name     string `json:"name" bson:"name"`
	Avatar   string `json:"avatar" bson:"avatar" `
	Email    string `json:"email"  bson:"avatar" validate:"required"`
}

type Register struct {
	Username        string `json:"username" bson:"username" validate:"required, min=4"`
	Password        string `json:"password" bson:"password" validate:"required,min=4"`
	PasswordConfirm string `json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
	Email           string `json:"email" bson:"email" binding:"required"`
}

type Login struct {
	Username string `json:"username" bson:"username" validate:"required, min=4"`
	Password string `json:"password" bson:"password" validate:"required,min=4"`
}

type EditUser struct {
	Username string `json:"username" bson:"username" validate:"required"`
	Name     string `json:"name" bson:"name"`
	Avatar   string `json:"avatar" bson:"avatar" `
	Email    string `json:"email"  bson:"avatar" validate:"required"`
}

type SearchResponse struct {
	User  []User      `json:"user"`
	Group []GroupChat `json:"group"`
}

func FilteredResponse(user *User) *UserResponse {
	return &UserResponse{
		Uid:      user.Uid,
		Name:     user.Name,
		Avatar:   user.Avatar,
		Username: user.Username,
		Email:    user.Email,
	}
}
