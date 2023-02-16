package main

import (
	"fmt"
	"net/http"

	"github.com/brutalzinn/api-task-list/configs"
	apikey_route "github.com/brutalzinn/api-task-list/routes/apikey"
	login_route "github.com/brutalzinn/api-task-list/routes/login"
	task_route "github.com/brutalzinn/api-task-list/routes/task"
	user_route "github.com/brutalzinn/api-task-list/routes/user"

	_ "github.com/brutalzinn/api-task-list/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           github.com/brutalzinn/api-task-list
// @version         1.0
// @description     Swagger example

// @host      localhost:9000
// @BasePath  /api/v1

func main() {
	err := configs.Load()
	if err != nil {
		panic((err))
	}
	route := chi.NewRouter()
	route.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	route.Mount("/swagger", httpSwagger.WrapHandler)

	task_route.TaskRoute(route)
	login_route.LoginRoute(route)
	user_route.UserRoute(route)
	apikey_route.ApiKeyRoute(route)
	fmt.Printf("Api started %s", configs.GetServerPort())
	http.ListenAndServe(fmt.Sprintf(":%s", configs.GetServerPort()), route)
}
