package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id          int32      `json:"id" gorm:"primaryKey;AUTO_INCREMENT;comment:'id'"`
	Uid         string     `json:"uid" gorm:"type:varchar(150);not null;unique_index:idx_uid;comment:'uid'"`
	Username    string     `json:"username" form:"username" binding:"required" gorm:"unique;not null; comment:'username'"`
	Password    string     `json:"password" form:"password" binding:"required" gorm:"type:varchar(150);not null;comment:'password'"`
	Avatar      string     `json:"avatar" gorm:"type:varchar(150);comment:'avatar'"`
	Email       string     `json:"email" gorm:"type:varchar(80);column:email;comment:'email'"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt      `json:"deletedAt"`
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
