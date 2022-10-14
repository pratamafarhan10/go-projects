package routes

import (
	"net/http"

	"github.com/go-projects/01-to_do_list_app/internal/controllers"
	MiddlewareAuth "github.com/go-projects/01-to_do_list_app/internal/middleware/auth"
	"github.com/julienschmidt/httprouter"
)

func Router() *httprouter.Router {
	router := httprouter.New()
	auth := controllers.NewAuthController()
	user := controllers.NewUserController()

	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)
	router.GET("/tes", MiddlewareAuth.VerifyJWT(auth.Tes))
	router.GET("/logout", MiddlewareAuth.VerifyJWT(auth.Logout))
	router.GET("/user", MiddlewareAuth.VerifyJWT(user.GetUser))
	router.POST("/user/update", MiddlewareAuth.VerifyJWT(user.UpdateUser))
	router.Handler("GET", "/assets/*filepath", http.StripPrefix("/assets", http.FileServer(http.Dir("../assets"))))
	return router
}
