package cmd

import "github.com/spf13/cobra"

var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "Query Datadog events",
	Long:  "Search and view Datadog events and alert history.",
}

func init() {
	rootCmd.AddCommand(eventCmd)
}
