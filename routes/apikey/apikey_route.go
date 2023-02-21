package apikey_route

import (
	apikey_controller "github.com/brutalzinn/api-task-list/controllers/apikey"
	"github.com/brutalzinn/api-task-list/middlewares"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware)
		r.Post("/apikey/generate", apikey_controller.Generate)
		r.Get("/apikey", apikey_controller.List)
		r.Post("/apikey/revoke/{id}", apikey_controller.Revoke)
	})

}
