package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/repository"
)

type groupService struct{}

var GroupSerivce = new(groupService)

func (g *groupService) SaveGroup(uid string, group *entity.GroupChat) error {
	db := repository.GetDB()
	var user entity.User
	db.AutoMigrate(group)
	result := db.Find(&user, "uid = ?", uid)
	if result.RowsAffected == 0 {
		return errors.New("no user with this uid exist")
	}
	group.UserId = user.Id
	group.Uid = uuid.New().String()
	db.Save(&group)

	groupMember := entity.GroupMember{
		UserId:  user.Id,
		GroupId: group.ID,
		Name:    user.Username,
		Mute:    false,
	}
	db.AutoMigrate(groupMember)
	db.Save(&groupMember)
	return nil
}
