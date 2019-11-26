package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/rockaut/dynatrace-config/cmd/config"
	"github.com/rockaut/dynatrace-config/cmd/cli"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dynatrace-config",
		Short: "A small utility to pass Dynatrace API calls with files.",
		Long: `A small utility to pass Dynatrace API calls with files.
This might be especially usefull in gitops scenarios.`,
		PersistentPreRunE:	readConfig,
		PreRunE:			preFlight,
		RunE:				start,
		SilenceErrors:		true,
		SilenceUsage:		true,
	}

	// shaman version information (populated by go linker)
	// -ldflags="-X main.version=${tag} -X main.commit=${commit}"
	version			= "unset"
	commit			= "unset"
)

// init add supported cli commands/flags
func init() {
	rootCmd.AddCommand(cli.ApplyCommand)
	rootCmd.AddCommand(cli.DiffCommand)
	rootCmd.AddCommand(cli.GetCommand)
	rootCmd.AddCommand(cli.ValidateCommand)
	config.AddFlags(rootCmd)
}

// Execute is the main entrypoint and called in main.main
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readConfig(ccmd *cobra.Command, args []string) error {
	if err := config.LoadConfigFile(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	return nil
}

func preFlight(ccmd *cobra.Command, args []string) error {
	if config.Version {
		fmt.Printf("dynatrace-config %s (%s)\n", version, commit)
		os.Exit(0)
	}

	return nil
}

func start(ccmd *cobra.Command, args []string) error {
	return fmt.Errorf("no command given")
}
