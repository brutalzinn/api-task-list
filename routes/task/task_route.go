package task_route

import (
	task_controller "api-auto-assistant/controllers/task"
	"api-auto-assistant/middlewares"

	"github.com/go-chi/chi/v5"
)

func TaskRoute(route *chi.Mux) {
	route.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.Post("/task", task_controller.Create)
		r.Put("/task/{id}", task_controller.Update)
		r.Delete("/task/{id}", task_controller.Delete)
		r.Get("/task", task_controller.List)
		r.Get("/task/{id}", task_controller.Get)
	})

}
