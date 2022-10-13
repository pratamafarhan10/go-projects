package controllers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator"
	"github.com/go-projects/01-to_do_list_app/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
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
		sendErrorResponse(w, "user not found", http.StatusNotFound)
		return
	}

	sendSuccessResponse(w, res, http.StatusOK)
}

func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params, token *jwt.Token) {
	// Get User Request
	id, err := primitive.ObjectIDFromHex(r.FormValue("_id"))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req := models.UpdateUserRequest{
		Id:              id,
		Email:           r.FormValue("email"),
		OldPassword:     r.FormValue("oldPassword"),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("passwordConfirm"),
		FirstName:       r.FormValue("firstname"),
		LastName:        r.FormValue("lastname"),
	}

	// Validate request
	err = validator.New().Struct(&req)
	if err != nil {
		split := strings.Split(err.Error(), "\n")
		sendErrorResponse(w, split, http.StatusBadRequest)
		return
	}

	// Get email from token
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

	// Get current user
	user := models.User{}
	err = uc.UserModel.GetUser(bson.M{"email": email}, bson.M{}, &user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			sendErrorResponse(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	req.Picture = user.Picture

	// Compare old password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		data := map[string][]string{
			"oldPassword": {"old password is wrong"},
		}
		sendErrorResponse(w, data, http.StatusBadRequest)
		return
	}

	// Generate new password
	hp, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Store picture if any
	f, h, err := r.FormFile("picture")
	if err != http.ErrMissingFile {
		fileName, err := uc.StorePicture(w, f, h, user.Id.Hex(), user.Picture)
		if err != nil {
			return
		}
		req.Picture = fileName
	}

	err = user.UpdateUser(bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"email": req.Email, "password": string(hp), "firstname": req.FirstName, "lastname": req.LastName, "picture": req.Picture}})

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := models.UserResponse{
		Id:        user.Id.Hex(),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Picture:   req.Picture,
		Role:      user.Role,
	}

	sendSuccessResponse(w, data, http.StatusOK)
}

func (uc UserController) StorePicture(w http.ResponseWriter, f multipart.File, h *multipart.FileHeader, userId, currentPicture string) (string, error) {
	defer f.Close()
	if h.Size > 5*1000000 {
		data := map[string][]string{
			"picture": {"picture size is too big"},
		}
		sendErrorResponse(w, data, http.StatusBadRequest)
		return "", fmt.Errorf("picture size is too big")
	}

	imgType := strings.Split(h.Header["Content-Type"][0], "/")[1]

	if imgType != "jpeg" && imgType != "jpg" && imgType != "png" {
		data := map[string][]string{
			"picture": {"accepted file type: jpeg, jpg, png"},
		}
		sendErrorResponse(w, data, http.StatusBadRequest)
		return "", fmt.Errorf("accepted file type: jpeg, jpg, png")
	}

	// Picture
	filename := `/assets/` + userId + `.` + imgType
	// Delete current picture
	err := os.Remove(`..` + currentPicture)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return "", fmt.Errorf("internal server error")
	}

	nf, err := os.Create(`..` + filename)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return "", fmt.Errorf("internal server error")
	}
	defer nf.Close()
	io.Copy(nf, f)

	return filename, nil
}
