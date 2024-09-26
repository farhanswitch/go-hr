package healthcheck

import (
	"github.com/go-chi/chi/v5"
)

func initRoutes(router *chi.Mux, controller healthCheckController) {
	router.Get("/health-check/db/now", controller.getNowCtrl)
}
