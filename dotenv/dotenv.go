package dotenv

import (
	"github.com/joho/godotenv"
	"log"
)

func Config() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return err
}
