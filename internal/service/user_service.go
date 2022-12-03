package service

import (
	"errors"
	"fmt"

	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/repository"
	"github.com/tonybobo/go-chat/pkg/common/utils"

	"github.com/google/uuid"
)

type userService struct{}

var UserService = new(userService)

func (u *userService) Register(register *entity.Register) (user *entity.User, err error) {
	var newUser entity.User
	if register.Password != register.PasswordConfirm {
		return nil, errors.New("password not match")
	}
	db := repository.GetDB()
	var count int64
	db.Model(user).Where("username", register.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("Username is taken.")
	}
	hashPassword, err := utils.HashPassword(register.Password)
	if err != nil {
		return nil, err
	}
	newUser.Username = register.Username
	newUser.Email = register.Email
	newUser.Password = hashPassword
	newUser.Uid = uuid.New().String()

	db.Create(&newUser)
	return &newUser, nil
}

func (u *userService) Login(login *entity.Login) (*entity.User, bool) {
	db := repository.GetDB()
	var queryUser *entity.User
	db.First(&queryUser, "username = ?", login.Username)
	if err := utils.VerifyPassword(queryUser.Password, login.Password); err != nil {
		return nil, false
	} else {
		return queryUser, true
	}
}

func (u *userService) EditUserDetail(user *entity.EditUser) error {
	var queryUser *entity.User
	db := repository.GetDB()
	result := db.First(&queryUser, "username= ?", user.Username)
	if result.RowsAffected == 0 {
		return errors.New("no user with this username")
	}
	queryUser.Name = user.Name
	queryUser.Email = user.Email
	db.Save(queryUser)
	return nil
}

func (u *userService) GetUserDetails(uid string) *entity.User {
	var user *entity.User
	db := repository.GetDB()
	db.Select("uid", "username", "avatar", "name").First(&user, "uid = ?", uid)
	return user
}

func (u *userService) GetUsersOrGroupBy(name string) *entity.SearchResponse {
	var queryUser []entity.User
	db := repository.GetDB()
	db.Raw("SELECT uid , username , name , avatar FROM users WHERE name LIKE ?", fmt.Sprintf("%%%s%%", name)).Scan(&queryUser)
	var queryGroup []entity.GroupChat
	db.Raw("SELECT uid , name FROM group_chats WHERE name LIKE ?", fmt.Sprintf("%%%s%%", name)).Scan(&queryGroup)

	return &entity.SearchResponse{
		User:  queryUser,
		Group: queryGroup,
	}
}

func (u *userService) AddFriend(request *entity.FriendRequest) error {
	var user *entity.User
	db := repository.GetDB()
	result := db.First(&user , "uid = ?" , request.Uid)
	if result.RowsAffected == 0 {
		return errors.New("no user found")
	}

	var friend *entity.User
	result2:= db.First(&friend , "uid = ?" , request.FriendUid)
	if result2.RowsAffected ==  0 {
		return errors.New("no user found")
	}

	var userFriend *entity.UserFriend
	result3 := db.First(&userFriend , "user_id = ? AND friend_id = ?" , user.Id , friend.Id)

	if result3.RowsAffected > 0 {
		return errors.New("user has been added to your friend list")
	}

	addFriend := &entity.UserFriend{
		UserId: user.Id,
		FriendId: friend.Id,
	}

	db.Create(&addFriend)
	return nil
}

func (u *userService) GetFriends (uid string) (*[]entity.User , error) {
	db := repository.GetDB()
	var user *entity.User 
	result := db.First(&user , "uid = ?" , uid)
	if result.RowsAffected == 0 {
		return nil , errors.New("no user found")
	} 

	var friends *[]entity.User

	db.Raw("SELECT u.uid , u.username , u.avatar FROM user_friends as uf JOIN users as u ON uf.friend_id = u.id WHERE uf.user_id = ?" , user.Id).Scan(&friends)
	return friends , nil
}
