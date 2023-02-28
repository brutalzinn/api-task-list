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
		r.Post("/repo", repo_controller.Create)
		r.Patch("/repo/{id}", repo_controller.Patch)
		r.Delete("/repo/{id}", repo_controller.Delete)
		r.Group(func(r chi.Router) {
			r.Use(createRepoHypermedia().Handler)
			r.Get("/repo/paginate", repo_controller.Paginate)
			r.Get("/repo/{id}", repo_controller.Get)
		})
	})
}

func createRepoHypermedia() *hypermedia.HyperMedia {
	var links []hypermedia.HypermediaLink
	links = append(links, hypermedia.CreateHyperMedia("task_list", "/task/paginate?page=[page]&limit=[limit]&repo_id=%d&order=[DESC]", "GET"))
	links = append(links, hypermedia.CreateHyperMedia("delete", "/repo/%d", "DELETE"))
	links = append(links, hypermedia.CreateHyperMedia("update_one", "/repo/%d", "PATCH"))
	links = append(links, hypermedia.CreateHyperMedia("self", "/repo/%d", "GET"))
	options := hypermedia.HyperMediaOptions{
		Links: links,
	}
	hypermediaMiddleware := hypermedia.New(options)
	return hypermediaMiddleware
}
