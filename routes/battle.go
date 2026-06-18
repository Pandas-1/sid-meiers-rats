package routes

import (
    "rats/controller"
    "github.com/go-chi/chi/v5"
)

func battleRoutes(r chi.Router) {
    r.Post("/battle/start", controller.StartBattle)
    r.Get("/battle/ws", controller.BattleWS)
}