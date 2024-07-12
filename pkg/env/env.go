package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv(path string) {
	// this function will load the .env file
	// into the environment variables only if
	// it's local environment
	env := os.Getenv("APP_ENV")
	log.Println("APP_ENV: ", env)
	if env == "" || env == "local" {
		var err error
		if path != "" {
			err = godotenv.Load(path)
		} else {
			err = godotenv.Load()
		}
		if err != nil {
			log.Printf("Error loading .env file %v", err)
		}
	}
}
