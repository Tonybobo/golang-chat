package entity

import (
	"time"

	"gorm.io/gorm"
)

type Group struct {
	ID int32 `json:"id" gorm:"primaryKey"`
	Uid string `json:"uid" gorm:"type:varchar(250);not null;unique_index:idx_uid;comment:'uid'"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
	UserId int32 `json:"userId" gorm:"index;comment:'group owner id'"`
	Name string `json:"name" gorm:"type:varchar(250);comment:'group name'"`
	Notice string `json:"notice" gorm:"type:varchar(350); comment:'notice'"`
}