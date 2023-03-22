package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/smnspz/totem/internal/domain"
)

func GetToken(user *domain.User, baseUrl *string) (string, error) {
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

	var (
		authError    domain.AuthError
		authResponse domain.AuthResponse
	)

	if resp.StatusCode != 200 {
		json.Unmarshal(body, &authError)
		return "", errors.New(authError.Messages[0].Text)
	}

	json.Unmarshal(body, &authResponse)
	return authResponse.Token, nil

}
