package routes 

import (
	"net/http"
	"rats/controller"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Post("/register", controller.Register)
	r.Post("/login", controller.Login)

	r.Group(func(r chi.Router) {
		r.Use(controller.AuthMiddleware)
	})

	return r
}