package apikey_route

import (
	apikey_controller "github.com/brutalzinn/api-task-list/controllers/apikey"
	"github.com/brutalzinn/api-task-list/middlewares"
	"github.com/brutalzinn/api-task-list/middlewares/hypermedia"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware)
		r.Post("/apikey/generate", apikey_controller.Generate)
		r.Post("/apikey/regenerate/{id}", apikey_controller.Regenerate)
		r.Delete("/apikey/delete/{id}", apikey_controller.Delete)
		r.Group(func(r chi.Router) {
			r.Use(createApiHyperMedia().Handler)
			r.Get("/apikey", apikey_controller.List)
		})
	})
}

func createApiHyperMedia() *hypermedia.HyperMedia {
	var links []hypermedia.HypermediaLink
	links = append(links, hypermedia.CreateHyperMedia("delete", "/apikey/delete/%s", "DELETE"))
	links = append(links, hypermedia.CreateHyperMedia("regenerate", "/apikey/regenerate/%s", "POST"))
	options := hypermedia.HyperMediaOptions{
		Links: links,
	}
	hypermediaMiddleware := hypermedia.New(options)
	return hypermediaMiddleware
}
