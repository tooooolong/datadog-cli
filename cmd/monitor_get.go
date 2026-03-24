package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var monitorGetCmd = &cobra.Command{
	Use:   "get <monitor-id>",
	Short: "Get a monitor by ID",
	Long:  "Retrieve detailed information about a specific Datadog monitor.",
	Args:  cobra.ExactArgs(1),
	RunE:  runMonitorGet,
}

func init() {
	monitorCmd.AddCommand(monitorGetCmd)
}

func runMonitorGet(cmd *cobra.Command, args []string) error {
	monitorID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid monitor ID %q: %w", args[0], err)
	}

	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV1.NewMonitorsApi(client)
	monitor, _, err := api.GetMonitor(ctx, monitorID)
	if err != nil {
		return fmt.Errorf("failed to get monitor %d: %w", monitorID, err)
	}

	if jsonOutput {
		return printJSON(monitor)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID:\t%d\n", derefInt64(monitor.Id))
	fmt.Fprintf(w, "Name:\t%s\n", derefStr(monitor.Name))
	fmt.Fprintf(w, "Type:\t%s\n", monitor.Type)
	fmt.Fprintf(w, "Query:\t%s\n", monitor.Query)
	fmt.Fprintf(w, "Message:\t%s\n", derefStr(monitor.Message))
	if monitor.OverallState != nil {
		fmt.Fprintf(w, "Status:\t%s\n", *monitor.OverallState)
	}
	if monitor.Priority.IsSet() && monitor.Priority.Get() != nil {
		fmt.Fprintf(w, "Priority:\t%d\n", *monitor.Priority.Get())
	}
	fmt.Fprintf(w, "Tags:\t%s\n", strings.Join(monitor.Tags, ", "))
	if monitor.Created != nil {
		fmt.Fprintf(w, "Created:\t%s\n", monitor.Created.Format("2006-01-02 15:04:05"))
	}
	if monitor.Modified != nil {
		fmt.Fprintf(w, "Modified:\t%s\n", monitor.Modified.Format("2006-01-02 15:04:05"))
	}
	return w.Flush()
}
