package cmd

import "github.com/spf13/cobra"

var metricCmd = &cobra.Command{
	Use:   "metric",
	Short: "Query Datadog metrics",
	Long:  "Query and explore Datadog metric data.",
}

func init() {
	rootCmd.AddCommand(metricCmd)
}
