package service

import (
	"context"
	"errors"
	"fmt"
	"reflect"

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

func (g *groupService) SaveGroup(uid string, group *entity.GroupChat) error {
	db := repository.GetDB()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user entity.User
	query := bson.D{{Key: "uid", Value: uid}}
	if err := db.Collection("users").FindOne(ctx, query).Decode(&user); err != nil {
		return err
	}
	group.ID = primitive.NewObjectID()
	group.UserId = user.Uid
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()
	group.Uid = uuid.New().String()
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
	_, err := db.Collection("groupMembers").InsertOne(ctx, &groupMember)
	log.Logger.Error("error", log.Any("error", err))
	return nil
}

func (g *groupService) JoinGroup(userUid string, groupUid string) error {
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

	var groupMember *entity.GroupMember

	if err := db.Collection("groupMembers").FindOne(ctx, bson.D{
		{Key: "$and", Value: bson.A{
			bson.M{"userId": user.Uid},
			bson.M{"groupId": group.Uid},
		}},
	}).Decode(&groupMember); err != nil {
		log.Logger.Error("error", log.Any("error :", err))
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

	return nil
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
