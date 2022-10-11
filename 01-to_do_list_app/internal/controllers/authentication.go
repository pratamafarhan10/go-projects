package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/go-projects/01-to_do_list_app/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

var SampleSecretKey = []byte("SecretYouShouldHide")

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

	err = validator.New().StructExcept(&req, "Picture")
	if err != nil {
		split := strings.Split(err.Error(), "\n")

		res := ErrorResponse{
			Errors: split,
		}

		bs, err := json.Marshal(res)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bs)
		return
	}

	alreadyTaken := req.EmailAlreadyTaken()
	if alreadyTaken {
		res := ErrorResponse{
			Errors: map[string][]string{
				"username": {"username already taken"},
			},
		}
		bs, err := json.Marshal(res)
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
	req.Id = primitive.NewObjectID()

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

func (ac AuthController) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := models.User{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = validator.New().StructExcept(&req, "FirstName", "LastName", "Picture", "Role")
	if err != nil {

		split := strings.Split(err.Error(), "\n")

		res := ErrorResponse{
			Errors: split,
		}

		bs, err := json.Marshal(res)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bs)
		return
	}

	// Get the user based on email
	user := models.User{}

	err = req.GetUser(bson.M{"email": req.Email}, bson.M{}, &user)
	if err == mongo.ErrNoDocuments {
		res := ErrorResponse{
			Errors: "user not found",
		}

		bs, err := json.Marshal(res)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bs)
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		res := ErrorResponse{
			Errors: "wrong password",
		}

		bs, err := json.Marshal(res)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(bs)
		return
	}

	token, err := ac.generateJWT(user.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	res := SuccessResponse{
		Data: map[string]string{
			"token": token,
		},
	}

	bs, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func (ac AuthController) generateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp":   time.Now().Add(10 * time.Minute),
		"email": email,
	})

	s, err := token.SignedString(SampleSecretKey)
	if err != nil {
		return "", err
	}
	fmt.Println(s)

	return s, nil
}

func (ac AuthController) Tes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintln(w, "Hello tes")
}
