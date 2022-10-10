package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/go-projects/01-to_do_list_app/backend/internal/models"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (ac AuthController) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := models.User{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	alreadyTaken := req.EmailAlreadyTaken()
	if alreadyTaken {
		r := ErrorResponse{
			Errors: map[string][]string{
				"username": {"username already taken"},
			},
		}
		bs, err := json.Marshal(r)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bs)
		return
	}

	hp, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Password = string(hp)

	id, err := req.InsertUser()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	s := SuccessResponse{
		Data: models.UserResponse{
			Id:        id,
			Email:     req.Email,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Picture:   req.Picture,
			Role:      req.Role,
		},
	}

	bs, err := json.Marshal(s)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(bs)
}
