package routes

import (
	"net/http"

	"github.com/go-projects/01-to_do_list_app/internal/controllers"
	MiddlewareAuth "github.com/go-projects/01-to_do_list_app/internal/middleware/auth"
	"github.com/julienschmidt/httprouter"
)

func RouterV1() *httprouter.Router {
	router := httprouter.New()
	auth := controllers.NewAuthController()
	user := controllers.NewUserController()
	todo := controllers.NewTodosController()

	router.POST("api/v1/register", auth.Register)
	router.GET("api/v1//verify/:token", auth.VerifyEmail)
	router.POST("api/v1//forgotpassword", auth.SendForgotPassword)
	router.GET("api/v1//forgotpassword/:token", auth.CheckForgotPasswordTokenValidity)
	router.POST("api/v1//forgotpassword/:token/update", auth.UpdatePassword)
	router.POST("api/v1//login", auth.Login)
	router.GET("api/v1//tes", MiddlewareAuth.VerifyJWT(auth.Tes))
	router.GET("api/v1//logout", MiddlewareAuth.VerifyJWT(auth.Logout))
	router.GET("api/v1//user", MiddlewareAuth.VerifyJWT(user.GetUser))
	router.POST("api/v1//user/update", MiddlewareAuth.VerifyJWT(user.UpdateUser))
	router.POST("api/v1/todos/create", MiddlewareAuth.VerifyJWT(todo.CreateTodos))
	router.POST("api/v1/todos/update", MiddlewareAuth.VerifyJWT(todo.UpdateTodoList))
	router.DELETE("api/v1/todos/delete", MiddlewareAuth.VerifyJWT(todo.DeleteTodoList))
	router.GET("api/v1/todos", MiddlewareAuth.VerifyJWT(todo.GetTodoList))
	router.Handler("GET", "api/v1//assets/*filepath", http.StripPrefix("api/v1//assets", http.FileServer(http.Dir("../assets"))))
	return router
}
