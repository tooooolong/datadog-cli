package cmd

import (
	"fmt"
	"strconv"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"

	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var monitorDeleteCmd = &cobra.Command{
	Use:   "delete <monitor-id>",
	Short: "Delete a monitor",
	Long: `Delete a Datadog monitor by ID.

Examples:
  datadog monitor delete 12345`,
	Args: cobra.ExactArgs(1),
	RunE: runMonitorDelete,
}

func init() {
	monitorCmd.AddCommand(monitorDeleteCmd)
}

func runMonitorDelete(cmd *cobra.Command, args []string) error {
	monitorID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid monitor ID %q: %w", args[0], err)
	}

	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV1.NewMonitorsApi(client)
	resp, _, err := api.DeleteMonitor(ctx, monitorID)
	if err != nil {
		return fmt.Errorf("failed to delete monitor %d: %s", monitorID, ddclient.FormatAPIError(err))
	}

	fmt.Printf("Monitor %d deleted successfully (deleted ID: %d).\n", monitorID, resp.GetDeletedMonitorId())
	return nil
}
