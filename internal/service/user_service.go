package service

import (
	"context"
	"errors"
	"time"

	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/repository"
	"github.com/tonybobo/go-chat/pkg/common/utils"
	"github.com/tonybobo/go-chat/pkg/global/log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userService struct{}

var UserService = new(userService)

func (u *userService) Register(register *entity.Register) (user *entity.User, err error) {
	var newUser entity.User
	if register.Password != register.PasswordConfirm {
		return nil, errors.New("password not match")
	}
	db := repository.GetDB().Collection("users")
	ctx , cancel := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancel()

	query := bson.D{{Key: "username" , Value: register.Username}}
	count , err := db.CountDocuments(ctx , query)
	if err != nil {
		log.Logger.Error("Fail to Count Document" , log.Any("Error" , err))
	}
	if count > 0 {
		return nil, errors.New("Username is taken.")
	}
	hashPassword, err := utils.HashPassword(register.Password)
	if err != nil {
		return nil, err
	}
	newUser.Id = primitive.NewObjectID()
	newUser.CreatedAt = time.Now()
	newUser.UpdatedAt = time.Now()
	newUser.Username = register.Username
	newUser.Email = register.Email
	newUser.Password = hashPassword
	newUser.Uid = uuid.New().String()

	db.InsertOne(ctx, &newUser)
	return &newUser, nil
}

func (u *userService) Login(login *entity.Login) (*entity.User, bool) {
	db := repository.GetDB().Collection("users")
	var queryUser *entity.User
	ctx , cancel := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancel()

	query := bson.D{{Key: "username" , Value: login.Username}}
	if err:= db.FindOne(ctx , query).Decode(&queryUser); err != nil{
		return nil , false
	}

	
	if err := utils.VerifyPassword(queryUser.Password, login.Password); err != nil {
		return nil, false
	} else {
		return queryUser, true
	}
}

func (u *userService) EditUserDetail(user *entity.EditUser) error {
	var queryUser *entity.User
	db := repository.GetDB().Collection("users")
	ctx , cancel := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "username" , Value: user.Username}}
	update := bson.D{
		{Key:"$set" , Value: bson.D{
			{Key: "name" , Value: user.Name  },
			{Key: "email" , Value: user.Email  },
		}},
	}
	
	if err:= db.FindOneAndUpdate(ctx , filter , update).Decode(&queryUser); err != nil{
		return err
	}
	return nil
}

func (u *userService) GetUserDetails(uid string) *entity.User {
	var user *entity.User
	db := repository.GetDB().Collection("users")
	ctx , cancel := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancel()

	query := bson.D{{Key: "uid" , Value: uid}}
	if err:= db.FindOne(ctx , query).Decode(&user); err != nil{
		return nil
	}
	return user
}

// func (u *userService) GetUsersOrGroupBy(name string) *entity.SearchResponse {
// 	var queryUser []entity.User
// 	db := repository.GetDB()
// 	db.Raw("SELECT uid , username , name , avatar FROM users WHERE name LIKE ?", fmt.Sprintf("%%%s%%", name)).Scan(&queryUser)
// 	var queryGroup []entity.GroupChat
// 	db.Raw("SELECT uid , name FROM group_chats WHERE name LIKE ?", fmt.Sprintf("%%%s%%", name)).Scan(&queryGroup)

// 	return &entity.SearchResponse{
// 		User:  queryUser,
// 		Group: queryGroup,
// 	}
// }

// func (u *userService) AddFriend(request *entity.FriendRequest) error {
// 	var user *entity.User
// 	db := repository.GetDB()
// 	result := db.First(&user , "uid = ?" , request.Uid)
// 	if result.RowsAffected == 0 {
// 		return errors.New("no user found")
// 	}

// 	var friend *entity.User
// 	result2:= db.First(&friend , "uid = ?" , request.FriendUid)
// 	if result2.RowsAffected ==  0 {
// 		return errors.New("no user found")
// 	}

// 	var userFriend *entity.UserFriend
// 	result3 := db.First(&userFriend , "user_id = ? AND friend_id = ?" , user.Id , friend.Id)

// 	if result3.RowsAffected > 0 {
// 		return errors.New("user has been added to your friend list")
// 	}

// 	addFriend := &entity.UserFriend{
// 		UserId: user.Id,
// 		FriendId: friend.Id,
// 	}

// 	db.Create(&addFriend)
// 	return nil
// }

// func (u *userService) GetFriends (uid string) (*[]entity.User , error) {
// 	db := repository.GetDB()
// 	var user *entity.User 
// 	result := db.First(&user , "uid = ?" , uid)
// 	if result.RowsAffected == 0 {
// 		return nil , errors.New("no user found")
// 	} 

// 	var friends *[]entity.User

// 	db.Raw("SELECT u.uid , u.username , u.avatar FROM user_friends as uf JOIN users as u ON uf.friend_id = u.id WHERE uf.user_id = ?" , user.Id).Scan(&friends)
// 	return friends , nil
// }
