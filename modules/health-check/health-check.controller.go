package healthcheck

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type healthCheckController struct {
	service healthCheckService
}

var controller healthCheckController

func (hcc healthCheckController) getNowCtrl(w http.ResponseWriter, r *http.Request) {
	res, err := hcc.service.GetNowSrvc()
	w.Header().Add("Content-Type", "application/json")
	if err != nil {
		log.Printf("Error getNowCtrl.\nError: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Something went wrong!"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println(res)
	strJson, _ := json.Marshal(res)
	w.Write([]byte(fmt.Sprintf(`{"data":%s}`, strJson)))
}
func factoryHealthCheckController(service healthCheckService) healthCheckController {
	if controller == (healthCheckController{}) {
		controller = healthCheckController{service}
	}
	return controller
}
