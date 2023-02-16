package main

import (
	"api-auto-assistant/configs"
	apikey_route "api-auto-assistant/routes/apikey"
	login_route "api-auto-assistant/routes/login"
	task_route "api-auto-assistant/routes/task"
	user_route "api-auto-assistant/routes/user"
	"fmt"
	"net/http"

	_ "api-auto-assistant/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           API-AUTO-ASSISTANT
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
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	route.Mount("/swagger", httpSwagger.WrapHandler)

	task_route.TaskRoute(route)
	login_route.LoginRoute(route)
	user_route.UserRoute(route)
	apikey_route.ApiKeyRoute(route)
	print("Api started")
	http.ListenAndServe(fmt.Sprintf(":%s", configs.GetServerPort()), route)
}
