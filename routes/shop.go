package routes

import (
    "rats/controller"
    "github.com/go-chi/chi/v5"
)

func shopRoutes(r chi.Router) {
    r.Get("/shop/buildings", controller.GetBuildings)
    r.Get("/shop/troops", controller.GetTroops)
}