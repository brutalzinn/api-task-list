package apikey_route

import (
	apikey_controller "github.com/brutalzinn/api-task-list/controllers/apikey"
	"github.com/brutalzinn/api-task-list/middlewares"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware)
		r.Get("/apikey", apikey_controller.List)
		r.Post("/apikey/generate", apikey_controller.Generate)
		r.Post("/apikey/regenerate/{id}", apikey_controller.Regenerate)
		r.Delete("/apikey/delete/{id}", apikey_controller.Delete)
	})

}
