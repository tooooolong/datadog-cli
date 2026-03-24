package cmd

import "github.com/spf13/cobra"

var errorCmd = &cobra.Command{
	Use:     "error",
	Aliases: []string{"errors"},
	Short:   "Query Datadog Error Tracking issues",
	Long:    "Search and view Error Tracking issues across traces, logs, and RUM.",
}

func init() {
	rootCmd.AddCommand(errorCmd)
}
