package connections

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func DbMysqlFactory() *sql.DB {
	if db == nil {
		dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
		if err != nil {
			log.Fatalln("Cannot read Database Port")
		}
		config := mysql.Config{
			User:   os.Getenv("DB_USER"),
			Passwd: os.Getenv("DB_PASS"),
			Net:    "tcp",
			Addr:   fmt.Sprintf("%s:%d", os.Getenv("DB_HOST"), dbPort),
			DBName: os.Getenv("DB_NAME"),
		}

		db, err = sql.Open("mysql", config.FormatDSN())
		if err != nil {
			log.Fatalf("Cannot connecto to Database.\nError: %s", err.Error())
		}
		pingErr := db.Ping()
		if pingErr != nil {
			log.Fatalf("Cannot ping the database. Error: %s", err.Error())
		}

		log.Println("Database connected!")
	}
	return db
}
