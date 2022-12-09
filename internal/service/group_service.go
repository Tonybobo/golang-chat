package service

import (
	"context"
	"errors"

	"time"

	"github.com/google/uuid"
	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/repository"
	"github.com/tonybobo/go-chat/pkg/global/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type groupService struct{}

var GroupSerivce = new(groupService)

func (g *groupService) GetGroups(uid string) ([]primitive.M , error) {
	db := repository.GetDB()
	ctx ,cancel := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancel()

	var queryUser *entity.User
	query := bson.D{{Key: "uid" , Value: uid}}
	if err := db.Collection("users").FindOne(ctx , query).Decode(&queryUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil , errors.New("No user found")
		}
		return nil , err
	} 
	

	var groups []primitive.M

	match := bson.D{
		{Key: "$match" , Value: bson.D{{Key: "userId" , Value: queryUser.Uid}}},
	}

	lookUp := bson.D{
		{
			Key: "$lookup" , Value: bson.D{
				{Key: "from" , Value: "groups"},
				{Key: "localField" , Value: "groupId"},
				{Key: "foreignField" , Value: "uid"},
				{Key: "as" , Value: "groupChat"},
			},
		},
	}

	group := bson.D{
		{Key: "$group" , Value: bson.D{
			{Key: "_id" , Value: "$userId"},
			{Key:"group" , Value: bson.D{{Key: "$push" ,Value:  bson.M{"$first" :"$$ROOT.groupChat"}}}},
		},
	},
	}

	
	result , err := db.Collection("groupMembers").Aggregate(ctx , mongo.Pipeline{match , lookUp , group})

	if err != nil {
		log.Logger.Error("aggregation" , log.Any("Error" , err))
	}

	if err = result.All(ctx , &groups); err != nil {
		log.Logger.Error("aggregation" , log.Any("Error" , err))
	}

	return groups , nil
}

func (g *groupService) SaveGroup(uid string, group *entity.GroupChat) error {
	db := repository.GetDB()
	ctx ,cancel := context.WithTimeout(context.Background() , 10*time.Second)
	defer cancel()
	var user entity.User
	query := bson.D{{Key:"uid" , Value: uid}}
	if err := db.Collection("users").FindOne(ctx , query).Decode(&user); err != nil {
		return err
	}
	group.ID = primitive.NewObjectID()
	group.UserId = user.Uid
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()
	group.Uid = uuid.New().String()
	db.Collection("groups").InsertOne(ctx , &group)

	groupMember := &entity.GroupMember{
		ID: primitive.NewObjectID(),
		UserId:  user.Uid,
		GroupId: group.Uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:    user.Username,
		Mute:    false,
	}
	_ , err := db.Collection("groupMembers").InsertOne(ctx ,&groupMember)
	log.Logger.Error("error" , log.Any("error" , err))
	return nil
}

// func (g *groupService) JoinGroup(userUid string , groupUid string) error {
// 	var user entity.User
// 	db := repository.GetDB()
// 	userResult := db.First(&user,"uid = ?" , userUid)
// 	if userResult.RowsAffected == 0 {
// 		return errors.New("no user found")
// 	}

// 	var group entity.GroupChat
// 	groupResult := db.First(&group , "uid = ?" , groupUid)
// 	if groupResult.RowsAffected == 0 {
// 		return errors.New("no group found")
// 	}

// 	var groupMember entity.GroupMember
// 	memberResult := db.First(&groupMember , "user_id = ? AND group_id = ?" , user.Id , group.ID )
// 	if memberResult.RowsAffected > 0 {
// 		return errors.New("user has been added in the group previously")
// 	}
// 	name := user.Name
// 	if name == ""{
// 		name = user.Username
// 	}

// 	insert := &entity.GroupMember{
// 		UserId: user.Id,
// 		GroupId: group.ID,
// 		Name: name,
// 		Mute: false,
// 	}
// 	db.Create(&insert)

// 	return nil
// }

// func (g *groupService) GetGroupUsers(uid string) (*[]entity.User , error) {
// 	var group entity.GroupChat
// 	db := repository.GetDB()
// 	result := db.First(&group , "uid = ? " , uid)
// 	if result.RowsAffected == 0 {
// 		return nil , errors.New("no group found")
// 	}

// 	var user *[]entity.User
// 	db.Raw("SELECT u.uid , u.avatar , u.username FROM group_chats AS g JOIN group_members as gm ON gm.group_id = g.id JOIN users as u ON u.id = gm.user_id WHERE g.id = ?" , group.ID).Scan(&user)

// 	return user , nil
// }