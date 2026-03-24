package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var monitorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List monitors",
	Long:  "List Datadog monitors with optional filters.",
	RunE:  runMonitorList,
}

var (
	listTags     string
	listName     string
	listPage     int64
	listPageSize int32
)

func init() {
	monitorListCmd.Flags().StringVar(&listTags, "tags", "", "Comma-separated list of tags to filter by")
	monitorListCmd.Flags().StringVar(&listName, "name", "", "Filter monitors by name")
	monitorListCmd.Flags().Int64Var(&listPage, "page", 0, "Page number (0-indexed)")
	monitorListCmd.Flags().Int32Var(&listPageSize, "page-size", 50, "Number of monitors per page")
	monitorCmd.AddCommand(monitorListCmd)
}

func runMonitorList(cmd *cobra.Command, args []string) error {
	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV1.NewMonitorsApi(client)
	opts := datadogV1.NewListMonitorsOptionalParameters().
		WithPage(listPage).
		WithPageSize(listPageSize)

	if listTags != "" {
		opts = opts.WithTags(listTags)
	}
	if listName != "" {
		opts = opts.WithName(listName)
	}

	monitors, _, err := api.ListMonitors(ctx, *opts)
	if err != nil {
		return fmt.Errorf("failed to list monitors: %w", err)
	}

	if jsonOutput {
		return printJSON(monitors)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tTYPE\tSTATUS\tTAGS")
	for _, m := range monitors {
		id := derefInt64(m.Id)
		name := derefStr(m.Name)
		mType := string(m.Type)
		status := ""
		if m.OverallState != nil {
			status = string(*m.OverallState)
		}
		tags := strings.Join(m.Tags, ",")
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", id, name, mType, status, tags)
	}
	return w.Flush()
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func derefStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}
