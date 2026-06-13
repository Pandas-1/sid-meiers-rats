package routes

import (
    "rats/controller"
    "github.com/go-chi/chi/v5"
)

func villageRoutes(r chi.Router) {
    r.Get("/village", controller.GetVillage)
    r.Post("/village/place", controller.PlaceBuilding)
    r.Put("/village/move", controller.MoveBuilding)
    r.Put("/village/upgrade", controller.UpgradeBuilding)
}