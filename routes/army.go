package routes

import (
    "rats/controller"
    "github.com/go-chi/chi/v5"
)

func armyRoutes(r chi.Router) {
    r.Get("/army", controller.GetArmy)
    r.Post("/army/train", controller.TrainArmy)
    r.Put("/army/upgrade", controller.LevelUpTroop)
}