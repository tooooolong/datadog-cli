package cmd

import (
	"fmt"
	"sort"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"

	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var serviceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List active APM services",
	Long: `List all services currently reporting to Datadog APM for a given environment.

Examples:
  datadog service list --env prod
  datadog service list --env staging`,
	RunE: runServiceList,
}

var serviceEnv string

func init() {
	serviceListCmd.Flags().StringVar(&serviceEnv, "env", "prod", "Environment to filter services")
	serviceCmd.AddCommand(serviceListCmd)
}

func runServiceList(cmd *cobra.Command, args []string) error {
	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV2.NewAPMApi(client)
	resp, _, err := api.GetServiceList(ctx, serviceEnv)
	if err != nil {
		return fmt.Errorf("failed to list services: %s", ddclient.FormatAPIError(err))
	}

	if jsonOutput {
		return printJSON(resp)
	}

	data := resp.GetData()
	attrs := data.GetAttributes()
	services := attrs.GetServices()

	if len(services) == 0 {
		fmt.Println("No services found.")
		return nil
	}

	sort.Strings(services)
	fmt.Printf("Active services in env:%s (%d total):\n\n", serviceEnv, len(services))
	for _, s := range services {
		fmt.Printf("  %s\n", s)
	}
	return nil
}
