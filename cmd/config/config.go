package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	homedir "github.com/mitchellh/go-homedir"
)

var (
	ApiUrl   = "" // URL for accessing the Dynatrace Api
	ApiToken = "" // Token for authenticating against the Dynatrace Api

	ConfigFile    = ""    // Configuration file to load
	EnvVarsPrefix = "DC_" // Prefix for automatic environment variables recognition
	Version       = false // Print version info and exit
)

// AddFlags adds the available cli flags
func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&EnvVarsPrefix, "envVarsPrefix", EnvVarsPrefix, "Prefix for automatic environment variables recognition")

	cmd.PersistentFlags().StringVarP(&ApiUrl, "apiurl", "u", ApiUrl, "URL for accessing the Dynatrace API")
	cmd.PersistentFlags().StringVarP(&ApiToken, "apitoken", "t", ApiToken, "Token for authenticating against the Dynatrace API")

	cmd.Flags().BoolVarP(&Version, "version", "v", Version, "Print version info and exit")
}

func LoadConfigFile() error {
	viper.SetDefault("envVarsPrefix", EnvVarsPrefix)
	viper.SetDefault("apiurl", "")
	viper.SetDefault("apitoken", "")

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Search config in home directory with name ".dynatrace-config" (without extension).
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	viper.SetConfigName(".dynatrace-config")
	viper.SetConfigFile(viper.GetString("config"))

	EnvVarsPrefix = viper.GetString("envVarsPrefix")

	viper.SetEnvPrefix(EnvVarsPrefix) // searches all env variables with that prefix
	viper.AutomaticEnv()              // read in environment variables that match)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			fmt.Println(err)
			return nil
		}
	}

	ApiUrl = viper.GetString("apiurl")
	ApiToken = viper.GetString("apitoken")

	return nil
}
