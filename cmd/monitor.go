package cmd

import "github.com/spf13/cobra"

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Manage Datadog monitors",
	Long:  "Query and update Datadog monitors.",
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
