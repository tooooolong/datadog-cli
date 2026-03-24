package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"

	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var logSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search logs",
	Long: `Search Datadog logs using a query string.

Query syntax follows Datadog log search:
  service:lego-v2-rest                   — logs from a specific service
  service:lego-v2-rest status:error      — error logs
  @http.status_code:5*                   — 5xx responses

Examples:
  datadog log search "service:lego-v2-rest" --from 1h
  datadog log search "service:panamera-rest status:error" --from 15m --limit 5`,
	Args: cobra.ExactArgs(1),
	RunE: runLogSearch,
}

var (
	logFrom  string
	logLimit int32
)

func init() {
	logSearchCmd.Flags().StringVar(&logFrom, "from", "15m", "Time range (e.g. 5m, 1h, 24h)")
	logSearchCmd.Flags().Int32Var(&logLimit, "limit", 10, "Maximum number of logs to return")
	logCmd.AddCommand(logSearchCmd)
}

func runLogSearch(cmd *cobra.Command, args []string) error {
	dur, err := parseDuration(logFrom)
	if err != nil {
		return fmt.Errorf("invalid --from duration %q: %w", logFrom, err)
	}

	now := time.Now()
	from := now.Add(-dur)

	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV2.NewLogsApi(client)
	sort := datadogV2.LOGSSORT_TIMESTAMP_DESCENDING
	opts := datadogV2.NewListLogsGetOptionalParameters().
		WithFilterQuery(args[0]).
		WithFilterFrom(from).
		WithFilterTo(now).
		WithSort(sort).
		WithPageLimit(logLimit)

	resp, _, err := api.ListLogsGet(ctx, *opts)
	if err != nil {
		return fmt.Errorf("failed to search logs: %s", ddclient.FormatAPIError(err))
	}

	if jsonOutput {
		return printJSON(resp)
	}

	logs := resp.GetData()
	if len(logs) == 0 {
		fmt.Println("No logs found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIME\tSTATUS\tSERVICE\tMESSAGE")
	for _, l := range logs {
		attrs := l.GetAttributes()
		ts := ""
		if attrs.Timestamp != nil {
			ts = attrs.Timestamp.Format("15:04:05")
		}
		status := derefStr(attrs.Status)
		svc := derefStr(attrs.Service)
		msg := derefStr(attrs.Message)
		msg = strings.ReplaceAll(msg, "\n", " ")
		if len(msg) > 100 {
			msg = msg[:97] + "..."
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", ts, status, svc, msg)
	}
	return w.Flush()
}
