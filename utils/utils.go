package utils

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Error while loading `.env` file")
	}

	return os.Getenv(key)
}

func ValidateEmail(email string) bool {
	if match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email); !match {
		return false
	}

	return true
}

func StringContainsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
