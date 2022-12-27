package service

import (
	"context"
	"errors"
	"math"

	"time"

	"github.com/tonybobo/go-chat/internal/entity"
	"github.com/tonybobo/go-chat/internal/repository"
	"github.com/tonybobo/go-chat/pkg/common/constant"
	"github.com/tonybobo/go-chat/pkg/common/utils"
	"github.com/tonybobo/go-chat/pkg/global/log"
	"github.com/tonybobo/go-chat/pkg/protocol"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type messageStruct struct{}

var MessageService = new(messageStruct)

func (m *messageStruct) GetMessages(limit int, page int, request *entity.MessageRequest) ([]primitive.M, float64, error) {
	db := repository.GetDB()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var messages []primitive.M

	if request.MessageType == constant.MESSAGE_TYPE_USER {
		username := []string{
			request.Uid, request.FriendUid,
		}
		cursor, err := db.Collection("messages").Find(ctx, bson.D{
			{Key: "$or", Value: bson.A{
				bson.M{
					"fromUserId": bson.M{"$in": username},
				},
				bson.M{
					"toUserId": bson.M{"$in": username},
				},
			}},
		}, utils.NewMongoPagination(limit, page).GetPaginatedOpts("createdAt", 1))

		if err != nil {
			log.Logger.Error("db error", log.String("error: ", err.Error()))
			return nil, 0, err
		}

		count, err := db.Collection("messages").CountDocuments(ctx, bson.D{
			{Key: "$or", Value: bson.A{
				bson.M{
					"fromUserId": bson.M{"$in": username},
				},
				bson.M{
					"toUserId": bson.M{"$in": username},
				},
			}},
		},
		)

		totalPage := math.Ceil(float64(count) / float64(limit))

		if err != nil {
			log.Logger.Error("db error", log.String("error: ", err.Error()))
			return nil, 0, err
		}

		if err = cursor.All(ctx, &messages); err != nil {
			log.Logger.Error("cursor error", log.String("error: ", err.Error()))
			return nil, 0, err
		}
		return messages, totalPage, err
	}

	if request.MessageType == constant.MESSAGE_TYPE_GROUP {
		match := bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "toUserId", Value: request.FriendUid},
				{Key: "messageType", Value: constant.MESSAGE_TYPE_GROUP},
			}},
		}

		sort := bson.D{
			{Key: "$sort", Value: bson.D{{Key: "createdAt", Value: 1}}},
		}

		skip := bson.D{{Key: "$skip", Value: page*limit - limit}}

		limitStage := bson.D{{Key: "$limit", Value: limit}}

		lookUp := bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "users"},
				{Key: "localField", Value: "fromUserId"},
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
				{Key: "as", Value: "sender"},
			},
			},
		}

		project := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "from", Value: bson.M{"$first": "$sender"}},
				{Key: "content", Value: 1},
				{Key: "contentType", Value: 1},
				{Key: "url", Value: 1},
				{Key: "createdAt", Value: 1},
			}},
		}

		cursor, err := db.Collection("messages").Aggregate(ctx, mongo.Pipeline{match, sort, limitStage, skip, lookUp, project})

		if err != nil {
			log.Logger.Error("db error", log.String("error: ", err.Error()))
			return nil, 0, err
		}

		if err := cursor.All(ctx, &messages); err != nil {
			log.Logger.Error("cursor error", log.String("error: ", err.Error()))
			return nil, 0, err
		}

		count, err := db.Collection("messages").CountDocuments(ctx, bson.D{
			{Key: "toUserId", Value: request.FriendUid},
			{Key: "messageType", Value: constant.MESSAGE_TYPE_GROUP},
		},
		)

		if err != nil {
			log.Logger.Error("db error", log.String("error: ", err.Error()))
			return nil, 0, err
		}

		totalPage := math.Ceil(float64(count) / float64(limit))

		return messages, totalPage, nil

	}
	return nil, 0, errors.New("unsupported message type")
}

func (m *messageStruct) SaveMessage(message *protocol.Message) {
	db := repository.GetDB()
	var fromUser *entity.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.Collection("users").FindOne(ctx, bson.D{{Key: "uid", Value: message.From}}).Decode(&fromUser); err != nil {
		log.Logger.Error("message error", log.String("error", err.Error()))
		return
	}

	var recipient string

	if message.MessageType == constant.MESSAGE_TYPE_USER {
		var target *entity.User
		if err := db.Collection("users").FindOne(ctx, bson.D{{Key: "uid", Value: message.To}}).Decode(&target); err != nil {
			log.Logger.Error("message error", log.String("error", err.Error()))
			return
		}
		recipient = target.Uid
	}

	if message.MessageType == constant.MESSAGE_TYPE_GROUP {
		var target *entity.GroupChat
		if err := db.Collection("groups").FindOne(ctx, bson.D{{Key: "uid", Value: message.To}}).Decode(&target); err != nil {
			log.Logger.Error("message error", log.String("error", err.Error()))
			return
		}
		recipient = target.Uid
	}

	insert := entity.Message{
		ID:          primitive.NewObjectID(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		FromUserId:  message.From,
		ToUserId:    recipient,
		Content:     message.Content,
		ContentType: int16(message.ContentType),
		MessageType: int16(message.MessageType),
		Url:         message.Url,
	}

	db.Collection("messages").InsertOne(ctx, &insert)
}
