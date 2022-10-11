package MiddlewareAuth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-projects/01-to_do_list_app/internal/controllers"
	"github.com/golang-jwt/jwt"
	"github.com/julienschmidt/httprouter"
)

func VerifyJWT(handler func(w http.ResponseWriter, r *http.Request, p httprouter.Params)) httprouter.Handle {
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
			log.Println(err.Error())
			log.Println(token.Valid)
			// http.Error(w, err.Error(), http.StatusUnauthorized)
			// return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			fmt.Println(claims)
		}

		if token.Valid {
			handler(w, r, p)
		} else {
			http.Error(w, "User not authenticated", http.StatusUnauthorized)
		}
	})
}
