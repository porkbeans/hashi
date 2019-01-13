package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = cobra.Command{
	Use:   "hashi",
	Short: "Download and install Hashicorp tools.",
	Long:  "Hashi is a tool for downloading and installing Hashicorp tools.",
}

func Execute() {
	rootCmd.SetOutput(os.Stdout)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.Execute()
}
