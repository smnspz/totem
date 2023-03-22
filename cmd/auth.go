/*
Copyright © 2023 Simone Rosani <s.rosani@anoki.it>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/smnspz/totem/internal/config"
	"github.com/smnspz/totem/internal/domain"
	http "github.com/smnspz/totem/internal/http"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		baseUrl := config.GetEnvVar("BASE_URL_DEV")
		emailRegexp := config.GetEnvVar("EMAIL_REGEXP")

		if isInteractive() {
			email = getEmail(emailRegexp)
			password = getPassword()
			saveConfigs(&domain.User{Email: &email, Password: &password})
		}

		token := http.GetToken(&domain.User{Email: &email, Password: &password}, baseUrl)
		fmt.Println("\nToken:", token)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.Flags().StringVarP(&email, "email", "u", "", "your anoki corporate email")
	authCmd.Flags().StringVarP(&password, "password", "p", "", "your anoki password")
	viper.BindPFlag("email", authCmd.Flags().Lookup("email"))
	viper.BindPFlag("password", authCmd.Flags().Lookup("password"))
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

func getEmail(emailRegexp *string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read stdin: %v", err)
	}
	match, _ := regexp.MatchString(*emailRegexp, email)
	if !match {
		fmt.Println("Your email must end with @anoki.it")
		getEmail(emailRegexp)
	}

	return strings.Trim(email, "\n")
}

func isInteractive() bool {
	return (email == "" && password == "") ||
		(email == "" && password != "") ||
		(email != "" && password == "")
}

func saveConfigs(user *domain.User) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nDo you want to save your credentials? (y/n) ")
	arg, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Error parsing stdin: %v", err)
	}
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Cannot find home folder %v", err)
	}
	totemPath := strings.Join([]string{home, ".totemconfig"}, "/")
	switch strings.Trim(arg, "\n") {
	case "y":
		viper.WriteConfigAs(totemPath)
		fmt.Println("You can find and edit your .totemconfig file under", totemPath)
	case "n":
		break
	default:
		viper.WriteConfigAs(totemPath)
		fmt.Println("You can find and edit your .totemconfig file under", totemPath)
	}
}
