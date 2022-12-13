package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-projects/02-rest_api_in_40_minutes/model"
	"github.com/julienschmidt/httprouter"
)

func GetUser(w http.ResponseWriter, _ *http.Request, _ httprouter.Params, user model.User) {
	b, err := json.Marshal(user)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
