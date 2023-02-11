package task_route

import (
	task_controller "api-auto-assistant/controllers/task"

	"github.com/go-chi/chi/v5"
)

func TaskRoute(route *chi.Mux) {
	route.Post("/task", task_controller.Create)
	route.Put("/task/{id}", task_controller.Update)
	route.Delete("/task/{id}", task_controller.Delete)
	route.Get("/task", task_controller.List)
	route.Get("/task/{id}", task_controller.Get)
}
