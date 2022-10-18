package MiddlewareAuth

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-projects/01-to_do_list_app/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func VerifyJWT(handler func(w http.ResponseWriter, r *http.Request, p httprouter.Params, user models.User)) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if r.Header["Authorization"] == nil {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
			return
		}

		tokenString := strings.Fields(r.Header["Authorization"][0])[1]
		token, err := jwt.Parse(tokenString, jwt.Keyfunc(func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		}))

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		user := models.User{}
		if ok && token.Valid {
			user.Email = claims["email"].(string)

			err = user.GetUser(bson.M{"email": user.Email}, bson.M{}, &user)
			if err == mongo.ErrNoDocuments {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}

			if user.Token != tokenString {
				http.Error(w, "User not authenticated", http.StatusUnauthorized)
				return
			}

			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
		handler(w, r, p, user)
	})
}
