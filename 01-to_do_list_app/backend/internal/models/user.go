package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id             primitive.ObjectID `bson:"_id" json:"_id"`
	Email          string             `bson:"email" json:"email"`
	Password       string             `bson:"password" json:"password"`
	FirstName      string             `bson:"firstname" json:"firstname"`
	LastName       string             `bson:"lastname" json:"lastname"`
	Picture        string             `bson:"picture" json:"picture"`
	Role           string             `bson:"role" json:"role"`
	ForgotPassword ForgotPassword     `bson:"forgotPassword" json:"forgotPassword"`
}

type UserResponse struct {
	Id        string `bson:"_id" json:"_id"`
	Email     string `bson:"email" json:"email"`
	FirstName string `bson:"firstname" json:"firstname"`
	LastName  string `bson:"lastname" json:"lastname"`
	Picture   string `bson:"picture" json:"picture"`
	Role      string `bson:"role" json:"role"`
}

type ForgotPassword struct {
	Token   string `bson:"token" json:"token"`
	Expires string `bson:"expires" json:"expires"`
}

func (user User) GetUser(projection bson.M, dst any) error {
	res := UserCollection.FindOne(context.Background(), bson.M{"_id": user.Id}, options.FindOne().SetProjection(projection))
	err := res.Decode(dst)

	return err
}

func (user User) InsertUser() (string, error) {
	res, err := UserCollection.InsertOne(context.Background(), user)

	oid, _ := res.InsertedID.(primitive.ObjectID)
	return oid.Hex(), err
}

func (user User) UpdateUser(update bson.M) error {
	_, err := UserCollection.UpdateOne(context.Background(), bson.M{"_id": user.Id}, update)

	return err
}

func (user User) EmailAlreadyTaken() bool {
	res := UserCollection.FindOne(context.Background(), bson.M{"email": user.Email})

	container := User{}
	err := res.Decode(&container)

	return err != mongo.ErrNoDocuments
}
