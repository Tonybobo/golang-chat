package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tonybobo/go-chat/config"
	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/repository"
	"github.com/tonybobo/go-chat/pkg/common/utils"
	"github.com/tonybobo/go-chat/pkg/global/log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userService struct{}

var UserService = new(userService)

func (u *userService) Register(register *entity.Register) (user *entity.User, err error) {
	var newUser entity.User
	if register.Password != register.PasswordConfirm {
		return nil, errors.New("password not match")
	}
	db := repository.GetDB().Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := bson.D{{Key: "username", Value: register.Username}}
	count, err := db.CountDocuments(ctx, query)
	if err != nil {
		log.Logger.Error("Fail to Count Document", log.Any("Error", err))
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
	newUser.Avatar = config.GetConfig().GCP.DefaultAvatar

	db.InsertOne(ctx, &newUser)
	return &newUser, nil
}

func (u *userService) Login(login *entity.Login) (*entity.User, bool) {
	db := repository.GetDB().Collection("users")
	var queryUser *entity.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := bson.D{{Key: "username", Value: login.Username}}
	if err := db.FindOne(ctx, query).Decode(&queryUser); err != nil {
		return nil, false
	}

	if err := utils.VerifyPassword(queryUser.Password, login.Password); err != nil {
		return nil, false
	} else {
		return queryUser, true
	}
}

func (u *userService) UploadUserAvatar(c *gin.Context)(*entity.User , error) {
	var queryUser *entity.User
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Logger.Error("Parse Form error", log.Any("error", err))
		return nil, err
	}

	username := c.Request.PostFormValue("username")

	db := repository.GetDB().Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "username", Value: username}}

	if err := db.FindOne(ctx, filter).Decode(&queryUser); err != nil {
		log.Logger.Error("error", log.Any("error", err))
		return nil, err
	}


	f, uploadedFile, _ := c.Request.FormFile("avatar")
	
	defer f.Close()


	if strings.Split(queryUser.Avatar, "avatar/")[1] == uploadedFile.Filename  {
		return nil , errors.New("same filename")
	}

	if err := utils.Uploader.DeleteImage(queryUser.Avatar); err != nil {
			return nil , err
	}
	avatar, err := utils.Uploader.UploadImage(f, "avatar/"+uploadedFile.Filename)

	if err != nil {
		return nil, err
	}

	var updatedUser *entity.User 

	update :=  bson.D{{Key: "$set", Value: bson.D{
			{Key: "avatar", Value: config.GetConfig().GCP.URL + avatar},
	}}}

	if err := db.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedUser); err != nil {
		log.Logger.Error("error", log.Any("error", err))
		return nil, err
	}

	return updatedUser , nil 
}

func (u *userService) EditUserDetail(c *gin.Context) (*entity.User, error) {
	var queryUser *entity.User
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Logger.Error("Parse Form error", log.Any("error", err))
		return nil, err
	}

	username := c.Request.PostFormValue("username")
	name := c.Request.PostFormValue("name")
	email := c.Request.PostFormValue("email")

	db := repository.GetDB().Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "username", Value: username}}


	if err := db.FindOne(ctx, filter).Decode(&queryUser); err != nil {
		log.Logger.Error("error", log.Any("error", err))
		return nil, err
	}

	update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: name},
				{Key: "email", Value: email},
			}},
		}


	var updatedUser *entity.User

	if err := db.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedUser); err != nil {
		log.Logger.Error("error", log.Any("error", err))
		return nil, err
	}

	return updatedUser, nil

}

func (u *userService) GetUserDetails(uid string) *entity.User {
	var user *entity.User
	db := repository.GetDB().Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := bson.D{{Key: "uid", Value: uid}}
	if err := db.FindOne(ctx, query).Decode(&user); err != nil {
		return nil
	}
	return user
}

func (u *userService) GetUsersOrGroupBy(name string) *entity.SearchResponse {
	var queryUser []entity.User
	db := repository.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	users, err := db.Collection("users").Find(ctx, bson.D{{Key: "username", Value: bson.D{{Key: "$regex", Value: name}, {Key: "$options", Value: "i"}}}})
	if err != nil {
		log.Logger.Error("error using text search", log.Any("error", err))
	}

	if err := users.All(ctx, &queryUser); err != nil {
		log.Logger.Error("error using text search", log.Any("error", err))
	}
	var queryGroup []entity.GroupChat

	groups, err := db.Collection("groups").Find(ctx, bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: name}, {Key: "$options", Value: "i"}}}})

	if err != nil {
		log.Logger.Error("error using text search", log.Any("error", err))
	}

	if err := groups.All(ctx, &queryGroup); err != nil {
		log.Logger.Error("error using text search", log.Any("error", err))
	}

	return &entity.SearchResponse{
		User:  queryUser,
		Group: queryGroup,
	}
}

func (u *userService) AddFriend(request *entity.FriendRequest) error {
	var user *entity.User
	db := repository.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.Collection("users").FindOne(ctx, bson.D{{Key: "uid", Value: request.Uid}}).Decode(&user); err != nil {
		return err
	}

	var friend *entity.User
	if err := db.Collection("users").FindOne(ctx, bson.D{{Key: "uid", Value: request.FriendUid}}).Decode(&friend); err != nil {
		return err
	}

	count, err := db.Collection("userFriend").CountDocuments(ctx, bson.D{
		{Key: "$and", Value: bson.A{bson.M{"userId": user.Uid}, bson.M{"friendId": friend.Uid}}},
	})

	if err != nil {
		log.Logger.Error("error ", log.Any("error :", err))
		return err
	}
	if count > 0 {
		log.Logger.Error("error ", log.String("error :", "already friend"))
		return errors.New("user has been added to your friend list")
	}

	addFriend := &entity.UserFriend{
		Id:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserId:    user.Uid,
		FriendId:  friend.Uid,
	}
	db.Collection("userFriend").InsertOne(ctx, &addFriend)

	return nil
}

func (u *userService) GetFriends(uid string) ([]primitive.M, error) {
	db := repository.GetDB()
	var user *entity.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.Collection("users").FindOne(ctx, bson.D{{Key: "uid", Value: uid}}).Decode(&user); err != nil {
		log.Logger.Error("error", log.Any("Get Friends", err))
		return nil, err
	}

	match := bson.D{{Key: "$match", Value: bson.D{{Key: "userId", Value: user.Uid}}}}
	lookUp := bson.D{
		{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "users"},
				{Key: "localField", Value: "friendId"},
				{Key: "foreignField", Value: "uid"},
				{Key: "pipeline", Value: bson.A{
					bson.M{"$project": bson.D{
						{Key: "_id", Value: 0},
						{Key: "createdAt", Value: 0},
						{Key: "updatedAt", Value: 0},
						{Key: "deletedAt", Value: 0},
						{Key: "password", Value: 0},
					}},
				}},
				{Key: "as", Value: "friends"},
			},
		},
	}

	group := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$userId"},
			{Key: "friends", Value: bson.D{{Key: "$push", Value: bson.M{"$first": "$$ROOT.friends"}}}},
		},
		},
	}

	project := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "friends", Value: "$friends"},
		}},
	}

	match1 := bson.D{{Key: "$match", Value: bson.D{{Key: "friendId", Value: user.Uid}}}}
	lookUp1 := bson.D{
		{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "users"},
				{Key: "localField", Value: "userId"},
				{Key: "foreignField", Value: "uid"},
				{Key: "pipeline", Value: bson.A{
					bson.M{"$project": bson.D{
						{Key: "_id", Value: 0},
						{Key: "createdAt", Value: 0},
						{Key: "updatedAt", Value: 0},
						{Key: "deletedAt", Value: 0},
						{Key: "password", Value: 0},
					}},
				}},
				{Key: "as", Value: "friends"},
			},
		},
	}

	group1 := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$userId"},
			{Key: "friends", Value: bson.D{{Key: "$push", Value: bson.M{"$first": "$$ROOT.friends"}}}},
		},
		},
	}

	result, err := db.Collection("userFriend").Aggregate(ctx, mongo.Pipeline{match, lookUp, group, project})

	if err != nil {
		log.Logger.Error("aggregation", log.Any("Error", err))
	}

	var allFriends []primitive.M

	for result.Next(ctx) {
		var res bson.M
		if err := result.Decode(&res); err != nil {
			log.Logger.Error("error", log.Any("error :", err))
		}
		fmt.Println(res)
		if res != nil {
			for _, value := range res["friends"].(primitive.A) {
				fmt.Println(value)

				allFriends = append(allFriends, value.(primitive.M))
			}
		}

	}

	result1, err := db.Collection("userFriend").Aggregate(ctx, mongo.Pipeline{match1, lookUp1, group1, project})

	if err != nil {
		log.Logger.Error("aggregation", log.Any("Error", err))
	}

	for result1.Next(ctx) {
		var res bson.M
		if err := result1.Decode(&res); err != nil {
			log.Logger.Error("error", log.Any("error :", err))
		}
		fmt.Println(res)

		if res != nil {
			for _, value := range res["friends"].(primitive.A) {
				fmt.Println(value)
				res := reflect.TypeOf(value)
				fmt.Print(res)
				allFriends = append(allFriends, value.(primitive.M))
			}
		}

	}

	return allFriends, nil
}
