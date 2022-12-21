package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tonybobo/go-chat/config"
	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/repository"
	"github.com/tonybobo/go-chat/pkg/common/utils"
	"github.com/tonybobo/go-chat/pkg/global/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type groupService struct{}

var GroupSerivce = new(groupService)

func (g *groupService) GetGroups(uid string) ([]primitive.M, error) {
	db := repository.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var queryUser *entity.User
	query := bson.D{{Key: "uid", Value: uid}}
	if err := db.Collection("users").FindOne(ctx, query).Decode(&queryUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("No user found")
		}
		return nil, err
	}

	var groups []primitive.M

	match := bson.D{
		{Key: "$match", Value: bson.D{{Key: "userId", Value: queryUser.Uid}}},
	}

	lookUp := bson.D{
		{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "groups"},
				{Key: "localField", Value: "groupId"},
				{Key: "foreignField", Value: "uid"},
				{Key: "as", Value: "groupChat"},
			},
		},
	}

	group := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$userId"},
			{Key: "group", Value: bson.D{{Key: "$push", Value: bson.M{"$first": "$$ROOT.groupChat"}}}},
		},
		},
	}

	result, err := db.Collection("groupMembers").Aggregate(ctx, mongo.Pipeline{match, lookUp, group})

	if err != nil {
		log.Logger.Error("aggregation", log.Any("Error", err))
	}

	if err = result.All(ctx, &groups); err != nil {
		log.Logger.Error("aggregation", log.Any("Error", err))
	}

	return groups, nil
}

func (g *groupService) SaveGroup(uid string, group *entity.GroupChat) (error) {
	db := repository.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user entity.User
	query := bson.D{{Key: "uid", Value: uid}}
	if err := db.Collection("users").FindOne(ctx, query).Decode(&user); err != nil {
		return err
	}

	count , err := db.Collection("groups").CountDocuments(ctx , bson.D{{Key: "name" , Value: group.Name}})

	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("name has been taken. Please think of a new name")
	}

	group.ID = primitive.NewObjectID()
	group.UserId = user.Uid
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()
	group.Uid = uuid.New().String()
	group.Avatar = config.GetConfig().GCP.DefaultGroupAvatar
	db.Collection("groups").InsertOne(ctx, &group)

	groupMember := &entity.GroupMember{
		ID:        primitive.NewObjectID(),
		UserId:    user.Uid,
		GroupId:   group.Uid,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      user.Username,
		Mute:      false,
	}
	_, err = db.Collection("groupMembers").InsertOne(ctx, &groupMember)
	log.Logger.Error("error", log.Any("error", err))
	return nil
}

func (g *groupService) JoinGroup(userUid string, groupUid string) (*entity.GroupChat ,error) {
	var user *entity.User
	db := repository.GetDB()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := db.Collection("users").FindOne(ctx, bson.D{{Key: "uid", Value: userUid}}).Decode(&user); err != nil {
		log.Logger.Error("error", log.Any("error :", err))
	}

	var group *entity.GroupChat

	if err := db.Collection("groups").FindOne(ctx, bson.D{{Key: "uid", Value: groupUid}}).Decode(&group); err != nil {
		log.Logger.Error("error", log.Any("error :", err))
	}

	count , err := db.Collection("groupMembers").CountDocuments(ctx, bson.D{
		{Key: "$and", Value: bson.A{
			bson.M{"userId": user.Uid},
			bson.M{"groupId": group.Uid},
		}},
	})

	if err !=  nil {
		return nil , err
	}

	if count > 0 {
		return nil , errors.New("already a group member ")
	}

	name := user.Name
	if name == "" {
		name = user.Username
	}

	insert := &entity.GroupMember{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserId:    user.Uid,
		GroupId:   group.Uid,
		Name:      name,
		Mute:      false,
	}
	db.Collection("groupMembers").InsertOne(ctx, &insert)

	return group ,  nil
}

func (g *groupService) GetGroupUsers(uid string) ([]primitive.M, error) {
	var group *entity.GroupChat
	db := repository.GetDB()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.Collection("groups").FindOne(ctx, bson.D{{Key: "uid", Value: uid}}).Decode(&group); err != nil {
		return nil, err
	}

	match := bson.D{
		{Key: "$match", Value: bson.D{{Key: "groupId", Value: group.Uid}}},
	}

	lookUp := bson.D{
		{Key: "$lookup", Value: bson.D{
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
			{Key: "as", Value: "members"},
		}},
	}

	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$groupId"},
			{Key: "members", Value: bson.D{{Key: "$push", Value: bson.M{"$first": "$$ROOT.members"}}}},
		},
		},
	}

	result, err := db.Collection("groupMembers").Aggregate(ctx, mongo.Pipeline{match, lookUp, groupStage})

	if err != nil {
		log.Logger.Error("aggregation", log.Any("Error", err))
	}

	var members []primitive.M

	for result.Next(ctx) {
		var res bson.M
		if err := result.Decode(&res); err != nil {
			log.Logger.Error("aggregation Decoding", log.Any("Error", err))
		}
		if res != nil {
			for _, value := range res["members"].(primitive.A) {
				fmt.Println(value)
				res := reflect.TypeOf(value)
				fmt.Print(res)
				members = append(members, value.(primitive.M))
			}
		}
	}

	return members, nil
}

func (g *groupService) UploadGroupAvatar(c *gin.Context) (*entity.GroupChat, error){
	var queryGroup *entity.GroupChat
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Logger.Error("Parse Form error", log.Any("error", err))
		return nil, err
	}
	groupId := c.Param("uid")
	db := repository.GetDB().Collection("groups")
	ctx , cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.FindOne(ctx , bson.D{{Key: "uid" , Value: groupId}}).Decode(&queryGroup); err != nil {
		log.Logger.Error("error", log.Any("error", err))
		return nil, err
	}

	f , uploadedFile , _ := c.Request.FormFile("avatar")

	defer f.Close()

	if strings.Split(queryGroup.Avatar, "avatar/group/")[1] == uploadedFile.Filename  {
		return nil , errors.New("same image")
	}


	if queryGroup.Avatar != config.GetConfig().GCP.DefaultGroupAvatar {
		if err := utils.Uploader.DeleteImage(queryGroup.Avatar); err != nil {
			return nil , err
		}
	}

	avatar, err := utils.Uploader.UploadImage(f, "avatar/group/"+uploadedFile.Filename)

	if err != nil {
		return nil, err
	}

	var updatedGroup *entity.GroupChat 

	update :=  bson.D{{Key: "$set", Value: bson.D{
			{Key: "avatar", Value: config.GetConfig().GCP.URL + avatar},
	}}}

	if err := db.FindOneAndUpdate(ctx, bson.D{{Key: "uid" , Value: queryGroup.Uid }}, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedGroup); err != nil {
		log.Logger.Error("error", log.Any("error", err))
		return nil, err
	}

	return updatedGroup , nil 
}

func (g* groupService) EditGroupDetail (c *gin.Context) (*entity.GroupChat , error ){
	
	uid := c.Param("uid")
	name := c.Request.PostFormValue("name")
	notice := c.Request.PostFormValue("notice")

	db := repository.GetDB().Collection("groups")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "uid", Value: uid}}

	update := bson.D{{
		Key: "$set" , Value: bson.D{
			{Key: "name" , Value: name},
			{Key:"notice" , Value: notice},
		},
	}}
	
	var updatedGroup *entity.GroupChat

	if err := db.FindOneAndUpdate(ctx , filter , update , options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updatedGroup); err != nil {
		log.Logger.Error("error", log.Any("error", err))
		return nil, err
	}
	
	return updatedGroup , nil
}
