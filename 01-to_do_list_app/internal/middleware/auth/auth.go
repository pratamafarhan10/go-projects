package MiddlewareAuth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-projects/01-to_do_list_app/internal/controllers"
	"github.com/go-projects/01-to_do_list_app/internal/models"
	"github.com/golang-jwt/jwt"
	"github.com/julienschmidt/httprouter"
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
			return controllers.SampleSecretKey, nil
		}))

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		user := models.User{}
		if ok && token.Valid {
			email := claims["email"].(string)
			user.Token = tokenString
			user.Email = email

			err = user.CheckToken(&user)
			if err == nil {
				http.Error(w, "User not authenticated", http.StatusUnauthorized)
				return
			} else if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}

		if token.Valid {
			handler(w, r, p, user)
		} else {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
		}
	})
}
