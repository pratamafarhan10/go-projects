package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/go-projects/01-to_do_list_app/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type AuthController struct{}

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "To-Do List ltd"

func NewAuthController() *AuthController {
	return &AuthController{}
}

func (ac AuthController) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := models.User{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error decoding body request", http.StatusInternalServerError)
		return
	}

	err = validator.New().StructExcept(&req, "Picture")
	if err != nil {
		split := strings.Split(err.Error(), "\n")
		sendErrorResponse(w, split, http.StatusNotFound)
		return
	}

	alreadyTaken := req.EmailAlreadyTaken()
	if alreadyTaken {
		data := map[string][]string{
			"email": {"email already taken"},
		}
		sendErrorResponse(w, data, http.StatusNotFound)
		return
	}

	hp, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.Password = string(hp)
	req.Id = primitive.NewObjectID()
	req.IsVerified = false

	req.Verification.Token = uuid.New().String()
	req.Verification.Expires = time.Now().Add(24 * 7 * time.Hour).Format(time.Layout)

	id, err := req.InsertUser()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ac.SendEmail(req, "To-Do List App User Verification", `
	<a href="http://localhost:8080/api/v1/verify/`+req.Verification.Token+`">Verify your email</a>
	`)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := models.UserResponse{
		Id:        id,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Picture:   req.Picture,
	}

	sendSuccessResponse(w, data, http.StatusCreated)
}

func (ac AuthController) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := models.User{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error decoding body request", http.StatusInternalServerError)
		return
	}

	err = validator.New().StructExcept(&req, "FirstName", "LastName", "Picture", "Role")
	if err != nil {
		split := strings.Split(err.Error(), "\n")
		sendErrorResponse(w, split, http.StatusNotFound)
		return
	}

	// Get the user based on email
	user := models.User{}

	err = req.GetUser(bson.M{"email": req.Email}, bson.M{}, &user)
	if err == mongo.ErrNoDocuments {
		sendErrorResponse(w, "user not found", http.StatusNotFound)
		return
	}

	if !user.IsVerified {
		http.Error(w, "User not verified", http.StatusUnauthorized)
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		sendErrorResponse(w, "wrong password", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := ac.generateJWT(user.Email)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Update user token
	err = user.UpdateUser(bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"token": token}})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Sending back response
	data := map[string]string{
		"token": token,
	}

	sendSuccessResponse(w, data, http.StatusOK)
}

func (ac AuthController) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params, user models.User) {
	err := user.UpdateUser(bson.M{"token": user.Token, "email": user.Email}, bson.M{"$set": bson.M{"token": ""}})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Logout success"))
}

func (ac AuthController) generateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"email": email,
	})

	s, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return s, nil
}

func (ac AuthController) SendEmail(user models.User, subject, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "this.mahanran@gmail.com")
	mailer.SetHeader("To", user.Email)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(CONFIG_SMTP_HOST, CONFIG_SMTP_PORT, os.Getenv("CONFIG_AUTH_EMAIL"), os.Getenv("CONFIG_AUTH_PASSWORD"))

	err := dialer.DialAndSend(mailer)
	return err
}

func (ac AuthController) VerifyEmail(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := models.User{}
	user.Verification.Token = p.ByName("token")

	err := user.GetUser(bson.M{"verification.token": user.Verification.Token}, bson.M{}, &user)
	if err == mongo.ErrNoDocuments {
		sendErrorResponse(w, "user not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if t, _ := time.Parse(time.Layout, user.Verification.Expires); time.Until(t) < 0 {
		user.Verification.Token = uuid.New().String()
		user.Verification.Expires = time.Now().Add(24 * 7 * time.Hour).Format(time.Layout)

		user.UpdateUser(bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"verification.token": user.Verification.Token, "verification.expires": user.Verification.Expires}})

		ac.SendEmail(user, "To-Do List App User Verification", `
		<a href="http://localhost:8080/api/v1/verify/`+user.Verification.Token+`">Verify your email</a>
		`)
		sendErrorResponse(w, "token expires", http.StatusNotFound)

		return
	}

	err = user.UpdateUser(bson.M{"email": user.Email, "_id": user.Id}, bson.M{"$set": bson.M{"isVerified": true, "verification.token": "", "verification.expires": ""}})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User verification succeed"))
}

func (ac AuthController) SendForgotPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := models.User{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error decoding body request", http.StatusInternalServerError)
		return
	}

	err = validator.New().StructExcept(&req, "Password", "FirstName", "LastName", "Picture", "Role")
	if err != nil {
		split := strings.Split(err.Error(), "\n")
		sendErrorResponse(w, split, http.StatusNotFound)
		return
	}

	err = req.GetUser(bson.M{"email": req.Email}, bson.M{}, &req)
	if err == mongo.ErrNoDocuments {
		sendErrorResponse(w, "user not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	req.ForgotPassword.Token = uuid.New().String()
	req.ForgotPassword.Expires = time.Now().Add(time.Hour).Format(time.Layout)

	err = req.UpdateUser(bson.M{"_id": req.Id}, bson.M{"$set": bson.M{"forgotPassword.token": req.ForgotPassword.Token, "forgotPassword.expires": req.ForgotPassword.Expires}})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = ac.SendEmail(req, "Forgot Password", `
	<a href="http://localhost:8080/api/v1/forgotpassword/`+req.ForgotPassword.Token+`">Change password</a>
	`)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	sendSuccessResponse(w, "set forgot password token successfull", http.StatusCreated)
}

func (ac AuthController) CheckForgotPasswordTokenValidity(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fpToken := p.ByName("token")

	user := models.User{}

	err := user.GetUser(bson.M{"forgotPassword.token": fpToken}, bson.M{}, &user)
	if err == mongo.ErrNoDocuments {
		sendErrorResponse(w, "token not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if t, _ := time.Parse(time.Layout, user.ForgotPassword.Expires); time.Until(t) < 0 {
		sendErrorResponse(w, "token expired", http.StatusBadRequest)
		return
	}
	sendSuccessResponse(w, "token is valid", http.StatusOK)
}

func (ac AuthController) UpdatePassword(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fpToken := p.ByName("token")

	// Get user request
	req := models.User{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error decoding body request", http.StatusInternalServerError)
		return
	}

	// Validate user request
	err = validator.New().StructExcept(&req, "Email", "FirstName", "LastName", "Picture", "Role")
	if err != nil {
		split := strings.Split(err.Error(), "\n")
		sendErrorResponse(w, split, http.StatusNotFound)
		return
	}

	// Get current user
	user := models.User{}
	err = user.GetUser(bson.M{"forgotPassword.token": fpToken}, bson.M{}, &user)
	if err == mongo.ErrNoDocuments {
		sendErrorResponse(w, "user not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if token expired
	if t, _ := time.Parse(time.Layout, user.ForgotPassword.Expires); time.Until(t) < 0 {
		sendErrorResponse(w, "token expired", http.StatusBadRequest)
		return
	}

	hp, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Update password and reset token
	err = req.UpdateUser(bson.M{"_id": user.Id}, bson.M{"$set": bson.M{"password": string(hp), "forgotPassword.token": "", "forgotPassword.expires": ""}})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send success response
	sendSuccessResponse(w, "reset password succeed", http.StatusOK)
}

func (ac AuthController) Tes(w http.ResponseWriter, r *http.Request, _ httprouter.Params, _ models.User) {
	fmt.Fprintln(w, "Hello tes")
}
