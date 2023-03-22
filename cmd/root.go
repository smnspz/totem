/*
Copyright © 2023 s.rosani@anoki.it
*/
package cmd

import (
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "totem",
		Short: "CLI tool for Anoki's Totem",
		Long:  "Daily activity reporting tool for Anoki's totem",
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatalf("Cannot find home folder %v", err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigType("toml")
		viper.SetConfigName(".totemconfig")
	}

	viper.AutomaticEnv()
	viper.ReadInConfig()
}
