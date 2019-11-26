package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	ApplyCommand = &cobra.Command{
		Use:   "apply",
		Short: "Applys the given configuration to the Dynatrace API",
		Run:   runApply,
	}
)

func runApply(cmd *cobra.Command, args []string) {
	fmt.Println("apply called")
}
