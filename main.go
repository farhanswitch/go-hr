package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	"github.com/farhanswitch/go-hr/configs"
	"github.com/farhanswitch/go-hr/connections"
	employee "github.com/farhanswitch/go-hr/modules/employee"
	healthcheck "github.com/farhanswitch/go-hr/modules/health-check"
)

const port string = ":8282"

type AppConfig struct {
	DB        *sql.DB
	Validator *validator.Validate
}

var appConfig AppConfig

func main() {

	initModules()
	initAppConfigs()
	defer appConfig.DB.Close()
	router := chi.NewRouter()
	plugin(router)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {

		res, err := appConfig.DB.Query("CALL sp_get_all_employee(?,?,?,?,?)", "45", 0, 10, "first_name", "ASC")
		defer res.Close()

		if err != nil {
			fmt.Println(err)
		}
		for res.Next() {
			var id int32
			var first, last string
			err := res.Scan(&id, &first, &last)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(id, first, last)
		}
		for res.NextResultSet() {
			if res.Next() {
				var a interface{}
				err = res.Scan(&a)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println(a)
			}
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"hello"}`))
	})
	initInternalService(router, appConfig.DB)
	log.Printf("Server listening on port %s", port)
	http.ListenAndServe(port, router)

}

func plugin(router *chi.Mux) {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
}

func initModules() {
	configs.InitModule("./env/local.env")
}
func initAppConfigs() {
	if appConfig == (AppConfig{}) {
		appConfig = AppConfig{
			DB:        connections.DbMysqlFactory(),
			Validator: validator.New(),
		}
	}

}
func initInternalService(router *chi.Mux, db *sql.DB) {
	healthcheck.InitModule(db, router)
	employee.InitModule(db, router, appConfig.Validator)
}
