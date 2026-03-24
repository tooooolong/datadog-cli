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

var eventListCmd = &cobra.Command{
	Use:   "list",
	Short: "List events",
	Long: `List Datadog events with optional filters.

Query syntax follows Datadog event search:
  source:monitor                      — all monitor events
  source:monitor monitor_id:12345     — events for a specific monitor
  source:monitor status:alert         — alert events
  priority:all                        — all priorities

Examples:
  datadog event list --query "source:monitor" --from 24h
  datadog event list --query "source:monitor monitor_id:28140542" --from 7d
  datadog event list --query "source:monitor status:alert" --from 1h`,
	RunE: runEventList,
}

var (
	eventQuery string
	eventFrom  string
	eventLimit int32
)

func init() {
	eventListCmd.Flags().StringVar(&eventQuery, "query", "source:monitor", "Event search query")
	eventListCmd.Flags().StringVar(&eventFrom, "from", "24h", "Time range (e.g. 1h, 24h, 7d)")
	eventListCmd.Flags().Int32Var(&eventLimit, "limit", 20, "Maximum number of events to return")
	eventCmd.AddCommand(eventListCmd)
}

func runEventList(cmd *cobra.Command, args []string) error {
	dur, err := parseDuration(eventFrom)
	if err != nil {
		return fmt.Errorf("invalid --from duration %q: %w", eventFrom, err)
	}

	now := time.Now()
	from := now.Add(-dur).Format(time.RFC3339)
	to := now.Format(time.RFC3339)

	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV2.NewEventsApi(client)
	opts := datadogV2.NewListEventsOptionalParameters().
		WithFilterQuery(eventQuery).
		WithFilterFrom(from).
		WithFilterTo(to).
		WithPageLimit(eventLimit)

	resp, _, err := api.ListEvents(ctx, *opts)
	if err != nil {
		return fmt.Errorf("failed to list events: %s", ddclient.FormatAPIError(err))
	}

	if jsonOutput {
		return printJSON(resp)
	}

	events := resp.GetData()
	if len(events) == 0 {
		fmt.Println("No events found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIME\tMESSAGE\tTAGS")
	for _, e := range events {
		attrs := e.GetAttributes()
		ts := ""
		if attrs.Timestamp != nil {
			ts = attrs.Timestamp.Format("2006-01-02 15:04:05")
		}
		msg := derefStr(attrs.Message)
		// Truncate long messages
		if len(msg) > 80 {
			msg = msg[:77] + "..."
		}
		// Replace newlines
		msg = strings.ReplaceAll(msg, "\n", " ")
		tags := strings.Join(attrs.Tags, ",")
		fmt.Fprintf(w, "%s\t%s\t%s\n", ts, msg, tags)
	}
	return w.Flush()
}
