package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TodoLists struct {
	Id      primitive.ObjectID `bson:"_id" json:"_id"`
	Date    string             `bson:"date" json:"date"`
	User_Id primitive.ObjectID `bson:"user_id" json:"user_id"`
	Tasks   []Task             `bson:"tasks" json:"tasks"`
}

type Task struct {
	Id        primitive.ObjectID `bson:"_id" json:"_id"`
	Title     string             `bson:"title" json:"title" validate:"required"`
	Completed bool               `bson:"completed" json:"completed" validate:"boolean"`
	Time      string             `bson:"time" json:"time" validate:"required"`
}

type TodoList struct {
	Id      primitive.ObjectID `bson:"_id" json:"_id"`
	Date    string             `bson:"date" json:"date" validate:"required,datetime=01-02-2006"`
	User_Id primitive.ObjectID `bson:"user_id" json:"user_id"`
	Task    Task               `bson:"task" json:"task" validate:"required"`
}

func (t TodoLists) GetTodoList(filter bson.M, projection bson.M, dst any) error {
	res := TodoListCollection.FindOne(context.TODO(), filter, options.FindOne().SetProjection(projection))
	err := res.Decode(dst)

	return err
}

func (t TodoLists) GetManyTodoLists(filter bson.M) ([]TodoLists, error) {
	cur, err := TodoListCollection.Find(context.TODO(), filter)
	if err != nil {
		return []TodoLists{}, err
	}

	// err = cur.All(context.TODO(), &dst)
	todolists := []TodoLists{}

	for cur.Next(context.TODO()) {
		data := TodoLists{}

		err = cur.Decode(&data)
		if err != nil {
			return []TodoLists{}, err
		}

		todolists = append(todolists, data)
	}

	return todolists, nil
}

func (t TodoList) InsertTodoList() error {
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

	_, err = TodoListCollection.UpdateOne(context.TODO(), bson.M{"date": t.Date, "user_id": t.User_Id}, bson.M{"$push": bson.M{"tasks": bson.M{"_id": t.Task.Id, "title": t.Task.Title, "completed": t.Task.Completed, "time": t.Task.Time}}})

	return err
}

func (t TodoList) UpdateTodoList() error {
	res, err := TodoListCollection.UpdateOne(
		context.TODO(),
		bson.M{"date": t.Date, "user_id": t.User_Id, "tasks": bson.M{"$elemMatch": bson.M{"_id": t.Task.Id}}},
		bson.M{"$set": bson.M{"tasks.$.title": t.Task.Title, "tasks.$.completed": t.Task.Completed, "tasks.$.time": t.Task.Time}},
	)

	if res.MatchedCount < 1 {
		return mongo.ErrNoDocuments
	}

	return err
}

func (t TodoList) DeleteTodoList() error {
	res, err := TodoListCollection.UpdateOne(
		context.TODO(),
		bson.M{"date": t.Date, "user_id": t.User_Id},
		bson.M{"$pull": bson.M{"tasks": bson.M{"_id": t.Task.Id}}},
	)
	if err != nil {
		return err
	}

	if res.MatchedCount < 1 {
		return mongo.ErrNoDocuments
	}

	_, err = TodoListCollection.DeleteOne(context.TODO(), bson.M{"date": t.Date, "user_id": t.User_Id, "tasks": bson.M{"$size": 0}})

	return err
}
