package healthcheck

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func InitModule(db *sql.DB, router *chi.Mux) {

	repo := factoryHealthCheckRepository(db)
	service := factoryHealthCheckService(repo)
	controller := factoryHealthCheckController(service)

	initRoutes(router, controller)
}
