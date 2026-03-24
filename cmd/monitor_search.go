package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var monitorSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search monitors",
	Long: `Search Datadog monitors using a query string.

Examples:
  datadog monitor search "type:metric status:alert"
  datadog monitor search "tag:env:prod"`,
	Args: cobra.ExactArgs(1),
	RunE: runMonitorSearch,
}

var (
	searchPage    int64
	searchPerPage int64
	searchSort    string
)

func init() {
	monitorSearchCmd.Flags().Int64Var(&searchPage, "page", 0, "Page number (0-indexed)")
	monitorSearchCmd.Flags().Int64Var(&searchPerPage, "per-page", 50, "Number of results per page")
	monitorSearchCmd.Flags().StringVar(&searchSort, "sort", "", "Sort field (e.g. name, status, tags)")
	monitorCmd.AddCommand(monitorSearchCmd)
}

func runMonitorSearch(cmd *cobra.Command, args []string) error {
	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV1.NewMonitorsApi(client)
	opts := datadogV1.NewSearchMonitorsOptionalParameters().
		WithQuery(args[0]).
		WithPage(searchPage).
		WithPerPage(searchPerPage)

	if searchSort != "" {
		opts = opts.WithSort(searchSort)
	}

	resp, _, err := api.SearchMonitors(ctx, *opts)
	if err != nil {
		return fmt.Errorf("failed to search monitors: %w", err)
	}

	if jsonOutput {
		return printJSON(resp)
	}

	monitors := resp.GetMonitors()
	meta := resp.GetMetadata()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tSTATUS\tTAGS")
	for _, m := range monitors {
		id := derefInt64(m.Id)
		name := derefStr(m.Name)
		status := ""
		if m.Status != nil {
			status = string(*m.Status)
		}
		tags := strings.Join(m.Tags, ",")
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", id, name, status, tags)
	}
	if err := w.Flush(); err != nil {
		return err
	}

	totalCount := meta.GetTotalCount()
	page := meta.GetPage()
	perPage := meta.GetPerPage()
	fmt.Fprintf(os.Stderr, "\nShowing page %d (%d per page), total: %d\n", page, perPage, totalCount)
	return nil
}
