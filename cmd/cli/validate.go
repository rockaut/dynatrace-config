package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	ValidateCommand = &cobra.Command{
		Use:   "validate",
		Short: "Validates the given configuration against the API.",
		Run:   runValidate,
	}
)

func runValidate(cmd *cobra.Command, args []string) {
	fmt.Println("validate called")
}
