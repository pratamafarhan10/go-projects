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
	Todos   []Todo             `bson:"todos" json:"todos" validate:"required"`
}

type Todo struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	Title     string             `bson:"title" json:"title" validate:"required"`
	Completed bool               `bson:"completed" json:"completed" validate:"required"`
	Time      time.Time          `bson:"time" json:"time"`
}

type TodosRequest struct {
	Date    time.Time          `bson:"date" json:"date" validate:"required,datetime=2006-02-01"`
	User_Id primitive.ObjectID `bson:"user_id" json:"user_id" validate:"required"`
	Todo    Todo               `bson:"todos" json:"todos" validate:"required"`
}

type TodoLists struct {
	Data []TodoList
}

func (t TodoList) GetTodoList(filter bson.M, projection bson.M, dst any) error {
	res := TodoListCollection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection))
	err := res.Decode(dst)

	return err
}

func (t TodoList) GetManyTodoLists(filter bson.M, dst TodoLists) error {
	cur, err := TodoListCollection.Find(context.TODO(), filter)
	if err != nil {
		return err
	}

	err = cur.All(context.TODO(), &dst)
	return err
}

func (t TodoList) InsertTodoList(tr TodosRequest) error {
	// Check if the data has been inserted before
	n, err := TodoListCollection.CountDocuments(context.TODO(), bson.M{"date": t.Date, "user_id": t.User_Id})
	if err != nil {
		return err
	}

	if n < 1 {
		_, err := TodoListCollection.InsertOne(context.TODO(), bson.M{"_id": t.Id, "date": t.Date, "user_id": t.User_Id})
		if err != nil {
			return err
		}
	}

	_, err = TodoListCollection.UpdateOne(context.TODO(), bson.M{"date": t.Date, "user_id": t.User_Id}, bson.M{"$push": bson.M{"todos": bson.M{"_id": tr.Todo.Id, "title": tr.Todo.Title, "completed": tr.Todo.Completed, "time": tr.Todo.Time}}})

	return err
}

func (t TodoList) UpdateTodoList(data TodosRequest) error {
	_, err := TodoListCollection.UpdateOne(
		context.TODO(),
		bson.M{"date": data.Date, "user_id": data.User_Id, "todos": bson.M{"$elemMatch": bson.M{"_id": data.Todo.Id}}},
		bson.M{"$set": bson.M{"todos.$.title": data.Todo.Title, "todos.$.completed": data.Todo.Completed, "todos.$.time": data.Todo.Time}},
	)

	return err
}

func (t TodoList) DeleteTodoList(data TodosRequest) error {
	_, err := TodoListCollection.UpdateOne(
		context.TODO(),
		bson.M{"date": data.Date, "user_id": data.User_Id},
		bson.M{"$pull": bson.M{"todos": bson.M{"_id": data.Todo.Id}}},
	)

	return err
}
