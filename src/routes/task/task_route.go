package task_route

import (
	task_controller "github.com/brutalzinn/api-task-list/controllers/task"
	"github.com/brutalzinn/api-task-list/middlewares"
	"github.com/brutalzinn/api-task-list/middlewares/hypermedia"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware)
		r.Use(middlewares.ApiKeyMiddleware)
		r.Use(createHyperMedia().Handler)
		r.Route("/task", func(r chi.Router) {
			r.Post("/", task_controller.Create)
			r.Patch("/", task_controller.Patch)
			r.Put("/", task_controller.Put)
			r.Delete("/{id}", task_controller.Delete)
			r.Get("/paginate", task_controller.Paginate)
			r.Get("/{id}", task_controller.Get)
		})
	})
}

func createHyperMedia() *hypermedia.HyperMedia {
	var links []hypermedia.HypermediaLink
	links = append(links, hypermedia.CreateHyperMedia("delete", "/task/%d", "DELETE"))
	links = append(links, hypermedia.CreateHyperMedia("update_one", "/task/%d", "PATCH"))
	links = append(links, hypermedia.CreateHyperMedia("update_all", "/task/%d", "PUT"))
	links = append(links, hypermedia.CreateHyperMedia("_self", "/task/%d", "GET"))
	options := hypermedia.HyperMediaOptions{
		Links: links,
	}
	hypermediaMiddleware := hypermedia.New(options)
	return hypermediaMiddleware
}
