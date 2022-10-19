package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/go-projects/01-to_do_list_app/internal/models"
	"github.com/julienschmidt/httprouter"
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

	for _, val := range req.Todos {
		val.Id = primitive.NewObjectID()
	}

}
