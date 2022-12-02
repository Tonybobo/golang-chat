package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id        int32          `json:"id" gorm:"primaryKey;AUTO_INCREMENT;comment:'id'"`
	Uid       string         `json:"uid" gorm:"type:varchar(150);not null;unique_index:idx_uid;comment:'uid'"`
	Username  string         `json:"username" form:"username" binding:"required" gorm:"unique;not null; comment:'username'"`
	Name      string         `json:"name" gorm:"comment:'Real Name'"`
	Password  string         `json:"password" form:"password" binding:"required" gorm:"type:varchar(150);not null;comment:'password'"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(150);comment:'avatar'"`
	Email     string         `json:"email" gorm:"type:varchar(80);column:email;comment:'email'"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}

type UserResponse struct {
	Uid      string `json:"uid" gorm:"type:varchar(150);not null;unique_index:idx_uid;comment:'uid'"`
	Username string `json:"username" form:"username" binding:"required" gorm:"unique;not null; comment:'username'"`
	Avatar   string `json:"avatar" gorm:"type:varchar(150);comment:'avatar'"`
	Email    string `json:"email" gorm:"type:varchar(80);column:email;comment:'email'"`
}

type Register struct {
	Username        string `json:"username" form:"username" binding:"required" gorm:"unique;not null; comment:'username'"`
	Password        string `json:"password" form:"password" binding:"required" gorm:"type:varchar(150);not null;comment:'password'"`
	PasswordConfirm string `json:"passwordConfirm" form:"password" gorm:"type:varchar(150);not null;comment:'password'"`
	Email           string `json:"email" gorm:"type:varchar(80);column:email;comment:'email'"`
}

type Login struct {
	Username string `json:"username" form:"username" binding:"required" gorm:"unique;not null; comment:'username'"`
	Password string `json:"password" form:"password" binding:"required" gorm:"type:varchar(150);not null;comment:'password'"`
}

type EditUser struct {
	Username string `json:"username" form:"username" binding:"required" gorm:"unique;not null; comment:'username'"`
	Name     string `json:"name" gorm:"comment:'Real Name'"`
	Email    string `json:"email" gorm:"type:varchar(80);column:email;comment:'email'"`
}

type SearchResponse struct {
	User  []User      `json:"user"`
	Group []GroupChat `json:"group"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("CreatedAt", time.Now())
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}

func FilteredResponse(user *User) UserResponse {
	return UserResponse{
		Uid:      user.Uid,
		Avatar:   user.Avatar,
		Username: user.Username,
		Email:    user.Email,
	}
}
