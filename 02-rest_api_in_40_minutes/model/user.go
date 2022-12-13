package model

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection

type User struct {
	Id       primitive.ObjectID `json:"_id" bson:"_id"`
	FullName string             `json:"fullName" bson:"fullName"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
}

func init() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://adminAccount1:password@localhost:27017/restapi40minutes"))
	if err != nil {
		log.Println("ERROR CONNECTING TO MONGODB")
		log.Fatalln(err)
	}

	UserCollection = client.Database("restapi40minutes").Collection("user")
}

func (u User) GetUser(filter bson.M, dst *User) error {
	res := UserCollection.FindOne(context.Background(), filter)

	err := res.Decode(&dst)

	return err
}
