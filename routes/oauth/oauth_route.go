package oauth_route

import (
	oauth_controller "github.com/brutalzinn/api-task-list/controllers/oauth"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.Get("/login", oauth_controller.Login)
		r.Post("/oauth/auth", oauth_controller.Authentication)
	})
}
