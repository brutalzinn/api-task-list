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
		r.Post("/task", task_controller.Create)
		r.Patch("/task", task_controller.Patch)
		r.Put("/task", task_controller.Put)
		r.Delete("/task/{id}", task_controller.Delete)
		r.Group(func(r chi.Router) {
			r.Use(createTaskHyperMedia().Handler)
			r.Get("/task/paginate", task_controller.Paginate)
			r.Get("/task/{id}", task_controller.Get)
		})
	})
}

func createTaskHyperMedia() *hypermedia.HyperMedia {
	var links []hypermedia.HypermediaLink
	links = append(links, hypermedia.CreateHyperMedia("delete", "/task/%d", "DELETE"))
	links = append(links, hypermedia.CreateHyperMedia("update_one", "/task/%d", "PATCH"))
	links = append(links, hypermedia.CreateHyperMedia("update_all", "/task/%d", "PUT"))
	links = append(links, hypermedia.CreateHyperMedia("self", "/task/%d", "GET"))
	options := hypermedia.HyperMediaOptions{
		Links: links,
	}
	hypermediaMiddleware := hypermedia.New(options)
	return hypermediaMiddleware
}
