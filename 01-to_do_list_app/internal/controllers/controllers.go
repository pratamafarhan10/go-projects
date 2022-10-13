package controllers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Errors interface{}
}

type SuccessResponse struct {
	Data interface{}
}

func sendErrorResponse(w http.ResponseWriter, data interface{}, status int) {
	res := ErrorResponse{
		Errors: data,
	}

	bs, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(bs)
}

func sendSuccessResponse(w http.ResponseWriter, data interface{}, status int) {
	s := SuccessResponse{
		Data: data,
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
