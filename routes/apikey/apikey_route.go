package apikey_route

import (
	apikey_controller "api-auto-assistant/controllers/apikey"
	"api-auto-assistant/middlewares"

	"github.com/go-chi/chi/v5"
)

func ApiKeyRoute(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware)
		r.Post("/apikey/generate", apikey_controller.Generate)
		r.Post("/apikey/revoke", apikey_controller.Revoke)
	})

}
