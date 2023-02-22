package task_route

import (
	task_controller "github.com/brutalzinn/api-task-list/controllers/task"
	"github.com/brutalzinn/api-task-list/middlewares"

	"github.com/go-chi/chi/v5"
)

func Register(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.Use(middlewares.JWTMiddleware)
		r.Use(middlewares.ApiKeyMiddleware)
		r.Post("/task", task_controller.Create)
		r.Put("/task/{id}", task_controller.Update)
		r.Delete("/task/{id}", task_controller.Delete)
		r.Get("/task", task_controller.List)
		r.Get("/task/paginate", task_controller.Paginate)
		r.Get("/task/{id}", task_controller.Get)
	})

}
