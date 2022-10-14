package models

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Db *mongo.Database
var UserCollection *mongo.Collection
var TodoListCollection *mongo.Collection

func init() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://admin:password@localhost:27017/todolistapp"))
	if err != nil {
		log.Println("ERROR CONNECTING TO MONGODB")
		log.Fatalln(err)
	}
	Client = client
	Db = client.Database("todolistapp")
	UserCollection = Db.Collection("users")
	TodoListCollection = Db.Collection("todolists")
}
