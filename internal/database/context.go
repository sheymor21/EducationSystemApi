package database

import (
	"context"
	mongo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const dbName = "qualificationsDb"

var mc MongoConfig

type MongoConfig struct {
	DbUri    string
	Username string
	Password string
}

type MongoClient struct {
	Users *mongo.Collection
}

func SetMongoConfig(data MongoConfig) {
	mc.DbUri = data.DbUri
	mc.Username = data.Username
	mc.Password = data.Password
}

func GetDatabaseConnection() (*MongoClient, *mongo.Client) {
	var auth options.Credential
	{
		auth.Password = mc.Password
		auth.Username = mc.Username
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mc.DbUri).SetAuth(auth))
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}
	db := client.Database(dbName)

	MongoEngine := &MongoClient{
		Users: db.Collection("Users"),
	}

	return MongoEngine, client
}

func CloseConnection(client *mongo.Client, ctx context.Context) {
	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
