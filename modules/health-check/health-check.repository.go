package healthcheck

import (
	"database/sql"
	"time"
)

type healthCheckRepo struct {
	db *sql.DB
}

var repo healthCheckRepo

type DbNow struct {
	Now time.Time `json:"now"`
}

func (hcr healthCheckRepo) GetNow() (DbNow, error) {
	var now DbNow
	var strTime string
	err := hcr.db.QueryRow("SELECT CURRENT_TIMESTAMP(6);").Scan(&strTime)
	if err != nil {
		return DbNow{}, err
	}
	var timeLayout string = "2006-01-02 15:04:05"
	timeNow, err := time.Parse(timeLayout, strTime)
	if err != nil {
		return DbNow{}, err
	}
	now.Now = timeNow
	return now, nil
}
func factoryHealthCheckRepository(db *sql.DB) healthCheckRepo {
	if repo == (healthCheckRepo{}) {
		repo = healthCheckRepo{db}
	}
	return repo
}
