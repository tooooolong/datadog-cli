package cmd

import "github.com/spf13/cobra"

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Query Datadog APM services",
	Long:  "List and explore Datadog APM services.",
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
