package main

import (
	"fmt"
	"net/http"

	"github.com/brutalzinn/api-task-list/configs"
	_ "github.com/brutalzinn/api-task-list/oauth"
	oauth_api_server "github.com/brutalzinn/api-task-list/oauth"

	_ "github.com/brutalzinn/api-task-list/docs"
	apikey_route "github.com/brutalzinn/api-task-list/routes/apikey"
	login_route "github.com/brutalzinn/api-task-list/routes/login"
	repo_route "github.com/brutalzinn/api-task-list/routes/repo"
	task_route "github.com/brutalzinn/api-task-list/routes/task"
	user_route "github.com/brutalzinn/api-task-list/routes/user"
	oauth_route "github.com/brutalzinn/api-task-list/routes/oauth"
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
	config := configs.GetConfig().API
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
	oauth_route.Register(route)
	oauth_api_server.InitOauthServer()
	fmt.Printf("API-AUTO-ASSISTANT STARTED WITH PORT %s", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", config.Port), route)
}
