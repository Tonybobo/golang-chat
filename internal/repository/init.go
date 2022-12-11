package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/tonybobo/go-chat/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _db *mongo.Database

func init() {

	username := config.GetConfig().Mongo.Username
	password := config.GetConfig().Mongo.Password
	database := config.GetConfig().Mongo.Database

	dsn := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.bkuvj3e.mongodb.net/?retryWrites=true&w=majority", username, password)

	client, err := mongo.NewClient(options.Client().ApplyURI(dsn))

	if err != nil {
		panic("Fail to connect to DB err:" + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	_db = client.Database(database)

	index := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "username", Value: "text"}},
		},
	}

	index2 := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "name", Value: "text"}},
		},
	}
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)

	_db.Collection("users").Indexes().CreateMany(
		context.TODO(),
		index,
		opts,
	)

	_db.Collection("groups").Indexes().CreateMany(
		ctx,
		index2,
		opts,
	)

}

func GetDB() *mongo.Database {
	return _db
}
