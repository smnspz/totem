/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authorize using anoki's credentials",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		token := getToken(&User{"a.scocco@anoki.it", "totempass"})
		fmt.Println("Token: ", token)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}

func getEnvVar(envToGet string) *string {
	// path, _ := filepath.Abs(".")
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to retrieve env vars: %v\n", err)
	}
	envVar := os.Getenv(envToGet)
	return &envVar
}

type User struct {
	email    string
	password string
}

func getToken(user *User) string {

	baseUrl := getEnvVar("BASE_URL_DEV")

	body, err := json.Marshal(map[string]string{
		"username": user.email,
		"password": user.password,
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
