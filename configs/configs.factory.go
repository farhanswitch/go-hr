package configs

import (
	"log"

	"github.com/joho/godotenv"
)

func envConfigFactory(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("Cannot load environtment file. Error: %s", err.Error())
	}
}
