package employee

import (
	"github.com/go-chi/chi/v5"
)

func initRoutes(router *chi.Mux, controller employeeController) {
	router.Route("/employees", func(r chi.Router) {
		r.Get("/get-all", controller.getAllCtrl)
		r.Get("/get-item/{employeeId}", controller.getItemCtrl)
		r.Post("/add", controller.createCtrl)
		r.Patch("/update", controller.updateCtrl)
	})

}
