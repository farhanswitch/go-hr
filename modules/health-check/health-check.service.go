package healthcheck

import (
	"log"
)

type healthCheckService struct {
	repo healthCheckRepo
}

var service healthCheckService

func (hcs healthCheckService) GetNowSrvc() (DbNow, error) {
	now, err := hcs.repo.GetNow()
	if err != nil {
		log.Printf("Error when getNow from DB.\nError: %s\n", err.Error())
		return DbNow{}, err
	}
	return now, nil
}
func factoryHealthCheckService(repo healthCheckRepo) healthCheckService {
	if service == (healthCheckService{}) {
		service = healthCheckService{repo}
	}
	return service
}
