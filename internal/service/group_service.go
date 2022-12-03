package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/repository"
)

type groupService struct{}

var GroupSerivce = new(groupService)

func (g *groupService) GetGroups(uid string) ([]entity.GroupResponse , error) {
	db := repository.GetDB()

	var queryUser *entity.User
	result := db.First(&queryUser , "uid = ?" , uid)
	if result.RowsAffected == 0 {
		return nil , errors.New("no user found")
	}

	var groups []entity.GroupResponse

	db.Raw("SELECT g.id AS group_id , g.uid , g.created_at , g.name , g.notice FROM group_members AS gm LEFT JOIN group_chats as g ON gm.group_id = g.id WHERE gm.user_id = ? " , queryUser.Id).Scan(&groups)

	return groups , nil
}

func (g *groupService) SaveGroup(uid string, group *entity.GroupChat) error {
	db := repository.GetDB()
	var user entity.User
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
	db.Save(&groupMember)
	return nil
}

func (g *groupService) JoinGroup(userUid string , groupUid string) error {
	var user entity.User
	db := repository.GetDB()
	userResult := db.First(&user,"uid = ?" , userUid)
	if userResult.RowsAffected == 0 {
		return errors.New("no user found")
	}

	var group entity.GroupChat 
	groupResult := db.First(&group , "uid = ?" , groupUid)
	if groupResult.RowsAffected == 0 {
		return errors.New("no group found")
	}

	var groupMember entity.GroupMember
	memberResult := db.First(&groupMember , "user_id = ? AND group_id = ?" , user.Id , group.ID )
	if memberResult.RowsAffected > 0 {
		return errors.New("user has been added in the group previously")
	}
	name := user.Name
	if name == ""{
		name = user.Username
	}

	insert := &entity.GroupMember{
		UserId: user.Id,
		GroupId: group.ID,
		Name: name,
		Mute: false,
	}
	db.Create(&insert)

	return nil
}

func (g *groupService) GetGroupUsers(uid string) (*[]entity.User , error) {
	var group entity.GroupChat
	db := repository.GetDB()
	result := db.First(&group , "uid = ? " , uid)
	if result.RowsAffected == 0 {
		return nil , errors.New("no group found")
	} 

	var user *[]entity.User 
	db.Raw("SELECT u.uid , u.avatar , u.username FROM group_chats AS g JOIN group_members as gm ON gm.group_id = g.id JOIN users as u ON u.id = gm.user_id WHERE g.id = ?" , group.ID).Scan(&user)

	return user , nil
}