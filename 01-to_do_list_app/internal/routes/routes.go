package routes

import (
	"github.com/go-projects/01-to_do_list_app/internal/controllers"
	MiddlewareAuth "github.com/go-projects/01-to_do_list_app/internal/middleware/auth"
	"github.com/julienschmidt/httprouter"
)

func Router() *httprouter.Router {
	router := httprouter.New()
	auth := controllers.NewAuthController()

	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)
	router.GET("/tes", MiddlewareAuth.VerifyJWT(auth.Tes))
	return router
}
