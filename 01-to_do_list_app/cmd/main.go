package main

import (
	"log"
	"net/http"

	"github.com/go-projects/01-to_do_list_app/env"
	"github.com/go-projects/01-to_do_list_app/internal/routes"
)

func main() {
	env.SetEnv()
	log.Fatalln(http.ListenAndServe(":8080", routes.Router()))
}
