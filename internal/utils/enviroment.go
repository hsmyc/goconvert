package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func GetStrapi() (string, string) {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	strapiUrl := os.Getenv("STRAPI_URL")
	strapiToken := os.Getenv("STRAPI_TOKEN")
	return strapiUrl, strapiToken
}
