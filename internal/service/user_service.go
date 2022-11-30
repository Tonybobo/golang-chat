package service

import (
	"errors"
	"time"

	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/repository"
	"github.com/tonybobo/go-chat/pkg/common/utils"

	"github.com/google/uuid"
)

type userService struct{}

var UserService = new(userService)

func (u *userService) Register(user *entity.User) error {
	if user.Password != user.PasswordConfirm {
		return errors.New("password not match")
	}
	db := repository.GetDB()
	db.AutoMigrate(user)
	var count int64
	db.Model(user).Where("username", user.Username).Count(&count)
	if count > 0 {
		return errors.New("Username is taken.")
	}
	user.CreatedAt = time.Now()
	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashPassword
	user.Uid = uuid.New().String()
	db.Create(&user)
	return nil
}

func (u *userService) Login(user *entity.User) error {
	return nil
}
