/*
Copyright © 2023 Simone Rosani <s.rosani@anoki.it>
*/
package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/smnspz/totem/internal/config"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	email    string
	password string
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authorize using anoki's credentials",
	Long: `Launch the auth subcommand to login interactively, 
or use the --username and --password flags

For example:
totem auth -u your.name@anoki.it -p yourpass

`,
	Run: func(cmd *cobra.Command, args []string) {
		if isInteractive() {
			email = getEmail()
			password = getPassword()
		}
		token := getToken(&User{&email, &password})
		fmt.Println("\nToken: ", token)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringVarP(&email, "email", "u", "", "your anoki corporate email")
	authCmd.Flags().StringVarP(&password, "password", "p", "", "your anoki password")
}

type User struct {
	email    *string
	password *string
}

func getToken(user *User) string {
	baseUrl := config.GetEnvVar("BASE_URL_DEV")

	body, err := json.Marshal(map[string]string{
		"username": *user.email,
		"password": *user.password,
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

func getPassword() string {
	tty, err := os.Open("/dev/tty")
	if err != nil {
		panic(err)
	}
	defer tty.Close()
	fmt.Print("Type your password: ")
	pwd, err := term.ReadPassword(int(tty.Fd()))
	if err != nil {
		panic(err)
	}
	return string(pwd)
}

func getEmail() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter email: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.Trim(username, "\n")
}

func isInteractive() bool {
	return (email == "" && password == "") || (email == "" && password != "") || (email != "" && password == "")
}
