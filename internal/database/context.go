package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

var mc MongoConfig
var m sync.Once
var mongoContext MongoContext

type MongoConfig struct {
	DbName   string
	DbUri    string
	Username string
	Password string
}

type MongoContext struct {
	Student  *mongo.Collection
	Teachers *mongo.Collection
	Marks    *mongo.Collection
	Client   *mongo.Client
}

func GetMongoContext() *MongoContext {
	return &mongoContext
}

func SetMongoConfig(data MongoConfig) {
	mc.DbName = data.DbName
	mc.DbUri = data.DbUri
	mc.Username = data.Username
	mc.Password = data.Password
}

func Run() {
	m.Do(func() {
		var auth options.Credential
		{
			auth.Password = mc.Password
			auth.Username = mc.Username
		}

		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mc.DbUri).SetAuth(auth))
		if err != nil {
			log.Fatal(err)
		}

		err = client.Ping(context.TODO(), nil)
		if err != nil {
			panic(err.Error())
		}

		db := client.Database(mc.DbName)
		mongoContext = MongoContext{
			Student:  db.Collection("Student"),
			Teachers: db.Collection("Teachers"),
			Marks:    db.Collection("Marks"),
			Client:   client,
		}
	})
}

func CloseConnection(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
}
