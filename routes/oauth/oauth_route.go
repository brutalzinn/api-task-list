package oauth_route

import (
	oauth_controller "github.com/brutalzinn/api-task-list/controllers/oauth"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.Get("/oauth/login", oauth_controller.LoginForm)
		r.Post("/oauth/token", oauth_controller.Token)
		r.Post("/oauth/authorize", oauth_controller.Authorize)
	})
}
