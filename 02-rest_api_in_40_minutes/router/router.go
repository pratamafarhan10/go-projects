package router

import (
	"github.com/go-projects/02-rest_api_in_40_minutes/controller"
	"github.com/go-projects/02-rest_api_in_40_minutes/middleware/auth"
	"github.com/julienschmidt/httprouter"
)

func Router() *httprouter.Router {
	router := httprouter.New()
	router.GET("/getUser", auth.Authentication(controller.GetUser))

	return router
}
