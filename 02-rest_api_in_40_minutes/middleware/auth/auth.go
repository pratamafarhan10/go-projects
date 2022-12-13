package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-projects/02-rest_api_in_40_minutes/model"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Authentication(handle func(http.ResponseWriter, *http.Request, httprouter.Params, model.User)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var userRequest model.User

		err := json.NewDecoder(r.Body).Decode(&userRequest)
		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}

		var user model.User
		err = userRequest.GetUser(bson.M{"email": userRequest.Email}, &user)
		fmt.Println(user)
		if err == mongo.ErrNoDocuments {
			log.Println("auth 1", err)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("user not found"))
			return
		}
		if err != nil {
			log.Println("auth 2", err)
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}

		if user.Password != userRequest.Password {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("user unauthorized"))
			return
		} else {
			handle(w, r, p, user)
		}
	}
}
