package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/tonybobo/go-chat/config"
	"github.com/tonybobo/go-chat/pkg/global/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _db *mongo.Database

func init() {


	username := config.GetConfig().Mongo.Username
	password := config.GetConfig().Mongo.Password
	database := config.GetConfig().Mongo.Database


	dsn := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.bkuvj3e.mongodb.net/?retryWrites=true&w=majority", username, password)

	client , err := mongo.NewClient(options.Client().ApplyURI(dsn))

	if err != nil {
		panic("Fail to connect to DB err:" + err.Error())
	}

	ctx , cancel := context.WithTimeout(context.Background() , 10*time.Second)

	defer cancel()
	err = client.Connect(ctx)
	if err != nil{
		panic(err)
	}
	log.Logger.Info("DB" , log.String("mongodb" , "connected"))
	
	_db = client.Database(database)
}

func GetDB() *mongo.Database  {
	return _db
}
