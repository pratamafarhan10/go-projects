package main

import (
	"log"
	"net/http"

	"github.com/go-projects/01-to_do_list_app/backend/internal/routes"
)

func main() {
	log.Fatalln(http.ListenAndServe(":8080", routes.Router()))
}
