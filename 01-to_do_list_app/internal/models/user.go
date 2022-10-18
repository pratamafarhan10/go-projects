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
	Email          string             `bson:"email" json:"email" validate:"required,email"`
	Password       string             `bson:"password" json:"password" validate:"required,min=8"`
	FirstName      string             `bson:"firstname" json:"firstname" validate:"required"`
	LastName       string             `bson:"lastname" json:"lastname" validate:"required"`
	Picture        string             `bson:"picture" json:"picture" validate:"required"`
	IsVerified     bool               `bson:"isVerified" json:"isVerified"`
	Role           string             `bson:"role" json:"role" validate:"required"`
	Token          string             `bson:"token" json:"token"`
	ForgotPassword ForgotPassword     `bson:"forgotPassword" json:"forgotPassword"`
	Verification   Verification       `bson:"verification" json:"verification"`
}

type UpdateUserRequest struct {
	Id              primitive.ObjectID `bson:"_id" json:"_id" validate:"required"`
	Email           string             `bson:"email" json:"email" validate:"required,email"`
	OldPassword     string             `bson:"oldPassword" json:"oldPassword" validate:"required,min=8"`
	Password        string             `bson:"password" json:"password" validate:"required,min=8,eqfield=PasswordConfirm"`
	PasswordConfirm string             `bson:"passwordConfirm" json:"passwordConfirm" validate:"required,min=8"`
	FirstName       string             `bson:"firstname" json:"firstname" validate:"required"`
	LastName        string             `bson:"lastname" json:"lastname" validate:"required"`
	Picture         string             `bson:"picture" json:"picture"`
}

type UserResponse struct {
	Id        string `bson:"_id" json:"_id"`
	Email     string `bson:"email" json:"email"`
	FirstName string `bson:"firstname" json:"firstname"`
	LastName  string `bson:"lastname" json:"lastname"`
	Picture   string `bson:"picture" json:"picture"`
}

type ForgotPassword struct {
	Token   string `bson:"token" json:"token"`
	Expires string `bson:"expires" json:"expires"`
}

type Verification struct {
	Token   string `bson:"token" json:"token"`
	Expires string `bson:"expires" json:"expires"`
}

func (user User) GetUser(filter bson.M, projection bson.M, dst any) error {
	res := UserCollection.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))
	err := res.Decode(dst)

	return err
}

func (user User) InsertUser() (string, error) {
	res, err := UserCollection.InsertOne(context.Background(), user)

	oid, _ := res.InsertedID.(primitive.ObjectID)
	return oid.Hex(), err
}

func (user User) UpdateUser(filter, update bson.M) error {
	_, err := UserCollection.UpdateOne(context.Background(), filter, update)

	return err
}

func (user User) CheckToken(dst any) error {
	res := UserCollection.FindOne(context.Background(), bson.M{"email": user.Email, "token": user.Token})
	err := res.Decode(dst)

	return err
}

func (user User) EmailAlreadyTaken() bool {
	res := UserCollection.FindOne(context.Background(), bson.M{"email": user.Email})

	container := User{}
	err := res.Decode(&container)

	return err != mongo.ErrNoDocuments
}
