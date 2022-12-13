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
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
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

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params, user models.User) {
	res := models.UserResponse{
		Id:        user.Id.Hex(),
		Email:     user.Email,
		FirstName: user.Email,
		LastName:  user.LastName,
		Picture:   user.Picture,
	}

	sendSuccessResponse(w, res, http.StatusOK)
}

func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params, user models.User) {
	// Get user request
	req := models.UpdateUserRequest{
		Id:              user.Id,
		Email:           r.FormValue("email"),
		OldPassword:     r.FormValue("oldPassword"),
		Password:        r.FormValue("password"),
		PasswordConfirm: r.FormValue("passwordConfirm"),
		FirstName:       r.FormValue("firstname"),
		LastName:        r.FormValue("lastname"),
	}

	// Validate request
	err := validator.New().Struct(&req)
	if err != nil {
		split := strings.Split(err.Error(), "\n")
		sendErrorResponse(w, split, http.StatusBadRequest)
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
	if currentPicture != "" {
		err := os.Remove(`..` + currentPicture)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return "", fmt.Errorf("internal server error")
		}
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
