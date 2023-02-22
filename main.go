package main

import (
	"fmt"
	"net/http"

	"github.com/brutalzinn/api-task-list/configs"
	apikey_route "github.com/brutalzinn/api-task-list/routes/apikey"
	login_route "github.com/brutalzinn/api-task-list/routes/login"
	repo_route "github.com/brutalzinn/api-task-list/routes/repo"
	task_route "github.com/brutalzinn/api-task-list/routes/task"
	user_route "github.com/brutalzinn/api-task-list/routes/user"

	_ "github.com/brutalzinn/api-task-list/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           github.com/brutalzinn/api-task-list
// @version         1.0
// @description     API TASK LIST

// @host      localhost:9000
// @BasePath  /api/v1

func main() {
	err := configs.Load()
	if err != nil {
		panic(err)
	}
	route := chi.NewRouter()
	route.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	route.Use(middleware.Logger)
	route.Mount("/swagger", httpSwagger.WrapHandler)
	task_route.Register(route)
	login_route.Register(route)
	user_route.Register(route)
	apikey_route.Register(route)
	repo_route.Register(route)

	port := configs.GetServerPort()
	fmt.Printf("API-AUTO-ASSISTANT STARTED WITH PORT %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), route)
}
