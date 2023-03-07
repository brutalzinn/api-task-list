package repo_route

import (
	repo_controller "github.com/brutalzinn/api-task-list/controllers/repo"
	"github.com/brutalzinn/api-task-list/middlewares"
	"github.com/brutalzinn/api-task-list/middlewares/hypermedia"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware)
		r.Use(middlewares.ApiKeyMiddleware)
		r.Use(createRepoHypermedia().Handler)
		r.Route("/repo", func(r chi.Router) {
			r.Post("/", repo_controller.Create)
			r.Patch("/{id}", repo_controller.Patch)
			r.Delete("/{id}", repo_controller.Delete)
			r.Get("/paginate", repo_controller.Paginate)
			r.Get("/{id}", repo_controller.Get)
		})
	})
}

func createRepoHypermedia() *hypermedia.HyperMedia {
	var links []hypermedia.HypermediaLink
	links = append(links, hypermedia.CreateHyperMedia("task_list", "/task/paginate?page=[page]&limit=[limit]&repo_id=%d&order=[DESC]", "GET"))
	links = append(links, hypermedia.CreateHyperMedia("delete", "/repo/%d", "DELETE"))
	links = append(links, hypermedia.CreateHyperMedia("update_one", "/repo/%d", "PATCH"))
	links = append(links, hypermedia.CreateHyperMedia("_self", "/repo/%d", "GET"))
	options := hypermedia.HyperMediaOptions{
		Links: links,
	}
	hypermediaMiddleware := hypermedia.New(options)
	return hypermediaMiddleware
}
