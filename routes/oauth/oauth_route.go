package oauth_route

import (
	oauth_controller "github.com/brutalzinn/api-task-list/controllers/oauth"
	"github.com/brutalzinn/api-task-list/middlewares"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.HandleFunc("/oauth/login", oauth_controller.LoginHandler)
		r.HandleFunc("/oauth/auth", oauth_controller.AuthHandler)
		r.HandleFunc("/oauth/authorize", oauth_controller.Authorize)
		r.HandleFunc("/oauth/token", oauth_controller.Token)
	})
	route.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware)
		r.HandleFunc("/oauth/generate", oauth_controller.Generate)
		// r.HandleFunc("/oauth/delete", oauth_controller.Test)
	})
}
