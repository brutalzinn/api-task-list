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
	route.Mount("/swagger", httpSwagger.WrapHandler)

	task_route.TaskRoute(route)
	login_route.LoginRoute(route)
	user_route.UserRoute(route)
	apikey_route.ApiKeyRoute(route)

	http.ListenAndServe(fmt.Sprintf(":%s", configs.GetServerPort()), route)
}
