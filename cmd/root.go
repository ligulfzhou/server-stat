package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"term-server-stat/cmd/commands"
)

var rootCmd = &cobra.Command{
	Use:   "server-stat",
	Short: "get server-stat in your ssh/config file",
	Long:  "get server-stat in your ssh/config file",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(commands.StartCmd)
}
