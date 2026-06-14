package routes

import (
    "rats/controller"
    "github.com/go-chi/chi/v5"
)

func matchmakingRoutes(r chi.Router) {
    r.Post("/matchmaking/find", controller.FindOpponent)
}