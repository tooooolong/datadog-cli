package cmd

import "github.com/spf13/cobra"

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Query Datadog logs",
	Long:  "Search and view Datadog log data.",
}

func init() {
	rootCmd.AddCommand(logCmd)
}
