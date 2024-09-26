package employee

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

func InitModule(db *sql.DB, router *chi.Mux, validate *validator.Validate) {
	repo := factoryEmployeeRepository(db)
	service := factoryEmployeeService(repo)
	controller := factoryEmployeeController(service, validate)
	initRoutes(router, controller)
}
