package login_route

import (
	user_controller "api-auto-assistant/controllers/user"

	"github.com/go-chi/chi/v5"
)

func LoginRoute(route *chi.Mux) {
	route.Post("/login", user_controller.Login)
}
