package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-projects/01-to_do_list_app/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	UserModel models.User
}

func NewUserController() *UserController {
	return &UserController{
		UserModel: models.User{},
	}
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params, token *jwt.Token) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	email, ok := claims["email"].(string)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	res := models.UserResponse{}

	err := uc.UserModel.GetUser(bson.M{"email": email}, bson.M{"forgotPassword": 0, "token": 0, "password": 0}, &res)
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

	bs, err := json.Marshal(SuccessResponse{Data: res})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params, token *jwt.Token) {
	// Get User Request
	req := models.UpdateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println(req)

	// err = validator.New().Struct(&req)
	// if err != nil {

	// 	split := strings.Split(err.Error(), "\n")

	// 	res := ErrorResponse{
	// 		Errors: split,
	// 	}

	// 	bs, err := json.Marshal(res)
	// 	if err != nil {
	// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write(bs)
	// 	return
	// }

	// claims, ok := token.Claims.(jwt.MapClaims)
	// if !ok {
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	// email, ok := claims["email"].(string)
	// if !ok {
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	// user := models.User{}

	// err = uc.UserModel.GetUser(bson.M{"email": email}, bson.M{}, &user)
	// if err != nil {
	// 	if err == mongo.ErrNoDocuments {
	// 		res := ErrorResponse{
	// 			Errors: "user not found",
	// 		}

	// 		bs, err := json.Marshal(res)
	// 		if err != nil {
	// 			http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 			return
	// 		}

	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		w.Write(bs)
	// 		return
	// 	}
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	// // Compare old password
	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	// if err != nil {
	// 	res := ErrorResponse{
	// 		Errors: map[string][]string{
	// 			"oldPassword": {"old password is wrong"},
	// 		},
	// 	}

	// 	bs, err := json.Marshal(res)
	// 	if err != nil {
	// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write(bs)
	// 	return
	// }

	// // Generate new password
	// hp, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	// if err != nil {
	// 	http.Error(w, "Internal server error", http.StatusInternalServerError)
	// 	return
	// }

	// // Create new file
	// os.create

	// user.Email = req.Email
	// user.Password = string(hp)
	// user.FirstName = req.FirstName
	// user.LastName = req.LastName
}
