package main

import (
	"log"
	"net/http"

	"github.com/go-projects/02-rest_api_in_40_minutes/router"
)

func main() {
	log.Fatalln(http.ListenAndServe(":8080", router.Router()))
}
