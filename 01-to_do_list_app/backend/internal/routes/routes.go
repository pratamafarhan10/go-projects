package routes

import (
	"github.com/go-projects/01-to_do_list_app/backend/internal/controllers"
	"github.com/julienschmidt/httprouter"
)

func Router() *httprouter.Router {
	router := httprouter.New()
	auth := controllers.NewAuthController()

	router.POST("/register", auth.Register)
	return router
}
