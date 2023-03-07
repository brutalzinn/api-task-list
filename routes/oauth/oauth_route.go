package oauth_route

import (
	oauth_controller "github.com/brutalzinn/api-task-list/controllers/oauth"
	"github.com/brutalzinn/api-task-list/middlewares"
	"github.com/brutalzinn/api-task-list/middlewares/hypermedia"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Route("/oauth", func(r chi.Router) {
		r.HandleFunc("/login", oauth_controller.LoginHandler)
		r.HandleFunc("/authorize", oauth_controller.Authorize)
		r.HandleFunc("/auth", oauth_controller.AuthHandler)
		r.HandleFunc("/token", oauth_controller.Token)
		r.Get("/test", oauth_controller.Test)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.JWTMiddleware)
			r.Use(createHyperMedia().Handler)
			r.Post("/generate", oauth_controller.Generate)
			r.Post("/regenerate/{id}", oauth_controller.Regenerate)
			r.Get("/list", oauth_controller.List)
		})

	})
}

func createHyperMedia() *hypermedia.HyperMedia {
	var links []hypermedia.HypermediaLink
	links = append(links, hypermedia.CreateHyperMedia("regenerate", "/oauth/regenerate/%d", "PATCH"))
	options := hypermedia.HyperMediaOptions{
		Links: links,
	}
	hypermediaMiddleware := hypermedia.New(options)
	return hypermediaMiddleware
}
