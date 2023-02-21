package user_route

import (
	user_controller "github.com/brutalzinn/api-task-list/controllers/user"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Post("/users", user_controller.Create)
	route.Put("/users/{id}", user_controller.Update)
	route.Delete("/users/{id}", user_controller.Delete)
	route.Get("/users", user_controller.List)
	route.Get("/users/{id}", user_controller.Get)
}
