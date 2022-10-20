package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/go-projects/01-to_do_list_app/internal/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodosController struct{}

func NewTodosController() TodosController {
	return TodosController{}
}

func (tc TodosController) CreateTodos(w http.ResponseWriter, r *http.Request, _ httprouter.Params, user models.User) {
	// Get user request
	req := models.TodoList{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Validate user request
	err = validator.New().Struct(&req)
	if err != nil {
		split := strings.Split(err.Error(), "\n")
		sendErrorResponse(w, split, http.StatusNotFound)
		return
	}

	req.Id = primitive.NewObjectID()
	req.Task.Id = primitive.NewObjectID()
	req.User_Id = user.Id

	err = req.InsertTodoList()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, "insert todolist successful", http.StatusCreated)
}

func (tc TodosController) GetTodoList(w http.ResponseWriter, r *http.Request, _ httprouter.Params, user models.User) {
	// Get query param if any
	date := r.FormValue("date")
	tl := models.TodoLists{}
	if date != "" {
		nd, err := time.Parse("01-02-2006", date)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		date = nd.Format(time.Layout)

		err = tl.GetTodoList(bson.M{"user_id": user.Id, "date": date}, bson.M{}, &tl)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		sendSuccessResponse(w, tl, http.StatusOK)
		return
	}
	todolists, err := tl.GetManyTodoLists(bson.M{"user_id": user.Id})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, todolists, http.StatusOK)
}
