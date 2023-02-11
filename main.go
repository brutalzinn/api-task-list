package main

import (
	"api-auto-assistant/configs"
	task_route "api-auto-assistant/routes/task"
	user_route "api-auto-assistant/routes/user"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	err := configs.Load()
	if err != nil {
		panic((err))
	}
	route := chi.NewRouter()
	task_route.TaskRoute(route)
	user_route.UserRoute(route)
	http.ListenAndServe(fmt.Sprintf(":%s", configs.GetServerPort()), route)

}
