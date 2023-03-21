package http

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/smnspz/totem/internal/domain"
)

func GetToken(user *domain.User, baseUrl *string) string {
	body, err := json.Marshal(map[string]string{
		"username": *user.Email,
		"password": *user.Password,
	})

	if err != nil {
		log.Fatalf("Failed to parse request payload: %v", err)
	}

	url := *baseUrl + "/jwt/login"

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("Failed to send auth request: %v", err)
	}
	defer resp.Body.Close()

	body, ioErr := io.ReadAll(resp.Body)
	if ioErr != nil {
		log.Fatalf("Failed to read response: %v", ioErr)
	}

	var retVal map[string]interface{}

	if err := json.Unmarshal(body, &retVal); err != nil {
		log.Fatalf("Failed to parse response: %v", err)
	}

	return retVal["token"].(string)
}
