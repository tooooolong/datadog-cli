package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	ddSite     string
)

var rootCmd = &cobra.Command{
	Use:   "datadog",
	Short: "Datadog CLI tool",
	Long:  "A command-line interface for interacting with the Datadog API.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().StringVar(&ddSite, "site", "", "Datadog site (e.g. datadoghq.eu, us5.datadoghq.com)")
}
