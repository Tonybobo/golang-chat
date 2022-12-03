package entity

import (
	"time"

	"gorm.io/gorm"
)

type GroupChat struct {
	ID        int32          `json:"id" gorm:"primaryKey"`
	Uid       string         `json:"uid" gorm:"type:varchar(250);not null;unique_index:idx_uid;comment:'uid'"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	UserId    int32          `json:"userId" gorm:"index;comment:'group owner id'"`
	Name      string         `json:"name" gorm:"type:varchar(250);comment:'group name'"`
	Notice    string         `json:"notice" gorm:"type:varchar(350); comment:'notice'"`
}

type GroupResponse struct {
	Uid       string         `json:"uid"`
	GroupId		int32 `json:"groupId"`
	CreatedAt time.Time `json:"createdAt"`
	Name 	string	`json:"name"`
	Notice 	string `json:"notice"`
}

func (g *GroupChat) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.SetColumn("CreatedAt", time.Now())
	return nil
}

func (g *GroupChat) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("UpdatedAt", time.Now())
	return nil
}
