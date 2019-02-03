package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "hashi",
	Short: "Download and install HashiCorp tools.",
	Long:  "Hashi is a tool for downloading and installing HashiCorp tools.",
}

// Execute runs a hashi command.
func Execute(args []string) int {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(installCmd)

	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(rootCmd.OutOrStderr(), "Err: %s\n", err)
		return 1
	}

	return 0
}
