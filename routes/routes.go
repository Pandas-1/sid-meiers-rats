package routes 

import (
	"net/http"
	"rats/controller"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
		fs.ServeHTTP(w, r)
	})
	r.Post("/register", controller.Register)
	r.Post("/login", controller.Login)
	r.Get("/battle/ws", controller.BattleWS)

	r.Group(func(r chi.Router) {
		r.Use(controller.AuthMiddleware)
		villageRoutes(r)
		armyRoutes(r)
		shopRoutes(r)
		matchmakingRoutes(r)
		battleRoutes(r)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/static/login.html", http.StatusFound)
	})
	return r
}