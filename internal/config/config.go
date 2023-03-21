package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvVar(envToGet string) *string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to retrieve env vars: %v\n", err)
	}
	envVar := os.Getenv(envToGet)
	return &envVar
}
