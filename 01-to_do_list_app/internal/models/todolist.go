package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TodoList struct {
	Id      primitive.ObjectID `bson:"_id" json:"_id"`
	Date    time.Time          `bson:"date" json:"date" validate:"required,datetime=2006-02-01"`
	User_Id primitive.ObjectID `bson:"user_id" json:"user_id" validate:"required"`
	Todos   []Todos            `bson:"todos" json:"todos" validate:"required"`
}

type Todos struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	Title     string             `bson:"title" json:"title" validate:"required"`
	Completed bool               `bson:"completed" json:"completed" validate:"required"`
	Time      time.Time          `bson:"date" json:"date"`
}

func (t TodoList) GetTodoList(filter bson.M, projection bson.M, dst any) error {
	res := TodoListCollection.FindOne(context.Background(), filter, options.FindOne().SetProjection(projection))
	err := res.Decode(dst)

	return err
}

func (t TodoList) InsertTodoList() (string, error) {
	// Check if the data has been inserted before
	n, err := TodoListCollection.CountDocuments(context.Background(), bson.M{"date": t.Date, "user_id": t.User_Id})
	if err != nil {
		return "", err
	}

	if n < 1 {
		res, err := TodoListCollection.InsertOne(context.Background(), t)

		oid, _ := res.InsertedID.(primitive.ObjectID)
		return oid.Hex(), err
	}

	for _, val := range t.Todos {
		_, err = TodoListCollection.UpdateOne(context.Background(), bson.M{"date": t.Date, "user_id": t.User_Id}, bson.M{"$push": val})
		if err != nil {
			return "", err
		}
	}

	return "", nil
}