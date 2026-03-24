package cmd

import "github.com/spf13/cobra"

var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Manage Datadog log pipelines",
	Long:  "List, create, and configure Datadog log processing pipelines.",
}

func init() {
	rootCmd.AddCommand(pipelineCmd)
}
