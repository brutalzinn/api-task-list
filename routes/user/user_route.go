package user_route

import (
	user_controller "github.com/brutalzinn/api-task-list/controllers/user"
	"github.com/brutalzinn/api-task-list/middlewares"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware)
		r.Route("/users", func(r chi.Router) {
			r.Put("/{id}", user_controller.Update)
			r.Delete("/{id}", user_controller.Delete)
			r.Get("/", user_controller.List)
			r.Get("/{id}", user_controller.Get)
		})
	})
}
