package oauth_route

import (
	oauth_controller "github.com/brutalzinn/api-task-list/controllers/oauth"
	"github.com/brutalzinn/api-task-list/middlewares"
	"github.com/brutalzinn/api-task-list/middlewares/hypermedia"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.HandleFunc("/oauth/login", oauth_controller.LoginHandler)
		r.HandleFunc("/oauth/auth", oauth_controller.AuthHandler)
		r.HandleFunc("/oauth/authorize", oauth_controller.Authorize)
		r.HandleFunc("/oauth/token", oauth_controller.Token)
		r.HandleFunc("/oauth/test", oauth_controller.Test)
		r.Group(func(r chi.Router) {
			r.Use(middlewares.JWTMiddleware)
			r.Use(createHyperMedia().Handler)
			r.Post("/oauth/generate", oauth_controller.Generate)
			r.Post("/oauth/regenerate/{id}", oauth_controller.RegenerateSecret)
			r.Patch("/oauth/update/{id}", oauth_controller.Update)
			r.Get("/oauth/list", oauth_controller.List)
		})
	})
}

func createHyperMedia() *hypermedia.HyperMedia {
	var links []hypermedia.HypermediaLink
	links = append(links, hypermedia.CreateHyperMedia("generate", "/oauth/generate/%d", "DELETE"))
	links = append(links, hypermedia.CreateHyperMedia("regenerate", "/oauth/regenerate/%d", "PATCH"))
	links = append(links, hypermedia.CreateHyperMedia("update", "/oauth/update/%d", "PUT"))
	links = append(links, hypermedia.CreateHyperMedia("list", "/oauth/list", "GET"))
	links = append(links, hypermedia.CreateHyperMedia("self", "/oauth/&d", "GET"))
	options := hypermedia.HyperMediaOptions{
		Links: links,
	}
	hypermediaMiddleware := hypermedia.New(options)
	return hypermediaMiddleware
}
