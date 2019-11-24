package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var envPrefixDefault string = "DC_"

var rootCmd = &cobra.Command{
	Use:   "dynatrace-config",
	Short: "A small utility to pass Dynatrace API calls with files.",
	Long: `A small utility to pass Dynatrace API calls with files.
This might be especially usefull in gitops scenarios.`,
}

// Execute is the main entrypoint and called in main.main
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dynatrace-config.yaml)")

	// local flags
	// none for now
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("envVarsPrefix", envPrefixDefault)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".dynatrace-config" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".dynatrace-config")
	}

	viper.SetEnvPrefix(viper.GetString("envVarsPrefix")) // searches all env variables with that prefix
	viper.AutomaticEnv()                                 // read in environment variables that match

	//fmt.Println("Using envVarsPrefix:", viper.GetString("envVarsPrefix"))
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
