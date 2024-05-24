package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetStrapi() (string, string) {
	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); !exists {
		err := godotenv.Load("../.env")
		if err != nil {
			fmt.Println("Error loading .env file")
		}
	}

	strapiUrl := os.Getenv("STRAPI_URL")
	strapiToken := os.Getenv("STRAPI_TOKEN")
	return strapiUrl, strapiToken
}

func GetPort() string {
	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); !exists {
		if err := godotenv.Load("../.env"); err != nil {
			log.Fatal("error loading .env file:", err)
		}
	}
	port := os.Getenv("PORT")
	address := ":" + port
	return address
}
