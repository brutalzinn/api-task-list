package repo_route

import (
	repo_controller "github.com/brutalzinn/api-task-list/controllers/repo"
	"github.com/brutalzinn/api-task-list/middlewares"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware)
		r.Use(middlewares.ApiKeyMiddleware)
		r.Post("/repo", repo_controller.Create)
		r.Put("/repo/{id}", repo_controller.Update)
		r.Delete("/repo/{id}", repo_controller.Delete)
		r.Get("/repo", repo_controller.List)
		r.Get("/repo/paginate", repo_controller.Paginate)
		r.Get("/repo/{id}", repo_controller.Get)
	})

}
