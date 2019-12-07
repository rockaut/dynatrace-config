package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	DiffCommand = &cobra.Command{
		Use:   "diff",
		Short: "Shows the differences for the given configuration and the acutal configuration.",
		Run:   runDiff,
	}
)

func runDiff(cmd *cobra.Command, args []string) {
	fmt.Println("diff called")
}
