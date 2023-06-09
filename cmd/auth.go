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
		reader := bufio.NewReader(os.Stdin)
		baseUrl := config.GetEnvVar("BASE_URL_DEV")
		emailRegexp := config.GetEnvVar("EMAIL_REGEXP")

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				if isInteractive() {
					email = readEmail(reader, emailRegexp)
					password = readPassword()
					saveConfigs(reader, &domain.User{Email: &email, Password: &password})
				}
			} else {
				log.Fatalf("Failed to parse config file %v", err)
			}
		}

		email = viper.GetString("email")
		password = viper.GetString("password")

		token, err := http.GetToken(&domain.User{Email: &email, Password: &password}, baseUrl)
		if err != nil {
			log.Fatalln(err)
		}
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

func readPassword() string {
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

func readEmail(reader *bufio.Reader, emailRegexp *string) string {
	fmt.Print("Enter email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read stdin: %v", err)
	}
	match, _ := regexp.MatchString(*emailRegexp, email)
	if !match {
		fmt.Println("Your email must end with @anoki.it")
		readEmail(reader, emailRegexp)
	}

	return strings.Trim(email, "\n")
}

func isInteractive() bool {
	return (email == "" && password == "") ||
		(email == "" && password != "") ||
		(email != "" && password == "")
}

func saveConfigs(reader *bufio.Reader, user *domain.User) {
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
