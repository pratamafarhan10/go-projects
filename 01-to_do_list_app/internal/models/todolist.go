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
	Time      time.Time          `bson:"time" json:"time"`
}

type UpdateTodosRequest struct {
	Date      time.Time          `bson:"date" json:"date" validate:"required,datetime=2006-02-01"`
	User_Id   primitive.ObjectID `bson:"user_id" json:"user_id" validate:"required"`
	TodoId    primitive.ObjectID `bson:"todo_id" json:"todo_id" validate:"required"`
	Title     string             `bson:"title" json:"title" validate:"required"`
	Completed bool               `bson:"completed" json:"completed" validate:"required"`
	Time      time.Time          `bson:"time" json:"time"`
}

type TodoLists struct {
	Data []TodoList
}

func (t TodoList) GetTodoList(filter bson.M, projection bson.M, dst any) error {
	res := TodoListCollection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection))
	err := res.Decode(dst)

	return err
}

func (t TodoList) GetTodoLists(filter bson.M, dst TodoLists) error {
	cur, err := TodoListCollection.Find(context.TODO(), filter)
	if err != nil {
		return err
	}

	err = cur.All(context.TODO(), &dst)
	return err
}

func (t TodoList) InsertTodoList() (string, error) {
	// Check if the data has been inserted before
	n, err := TodoListCollection.CountDocuments(context.TODO(), bson.M{"date": t.Date, "user_id": t.User_Id})
	if err != nil {
		return "", err
	}

	if n < 1 {
		res, err := TodoListCollection.InsertOne(context.TODO(), t)

		oid, _ := res.InsertedID.(primitive.ObjectID)
		return oid.Hex(), err
	}

	for _, val := range t.Todos {
		_, err = TodoListCollection.UpdateOne(context.TODO(), bson.M{"date": t.Date, "user_id": t.User_Id}, bson.M{"$push": val})
		if err != nil {
			return "", err
		}
	}

	return "", nil
}

func (t TodoList) UpdateTodoList(data UpdateTodosRequest) error {
	_, err := TodoListCollection.UpdateOne(
		context.TODO(),
		bson.M{"date": data.Date, "user_id": data.User_Id, "todos": bson.M{"$elemMatch": bson.M{"_id": data.TodoId}}},
		bson.M{"$set": bson.M{"todos.$.title": data.Title, "todos.$.completed": data.Completed, "todos.$.time": data.Time}},
	)

	return err
}
