package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	GetCommand = &cobra.Command{
		Use:   "get",
		Short: "Fetches the actual config from the given API call.",
		Run:   runGet,
	}
)

func init() {
	GetCommand.Flags().StringP("outfile", "o", "stdout", "The file to write to.")
}

func runGet(cmd *cobra.Command, args []string) {
	fmt.Println("get called")
}
