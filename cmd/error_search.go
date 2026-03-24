package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"

	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var errorSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search Error Tracking issues",
	Long: `Search Datadog Error Tracking issues.

Query syntax follows Datadog error tracking search:
  service:panamera-rest             — errors from a specific service
  error.type:RuntimeError           — errors by type
  status:open                       — only open issues

Track types: trace, logs, rum

Examples:
  datadog error search "service:panamera-rest" --track trace --from 24h
  datadog error search "status:open" --track logs --from 7d
  datadog error search "" --track trace --from 1h`,
	Args: cobra.MaximumNArgs(1),
	RunE: runErrorSearch,
}

var (
	errorSearchTrack string
	errorSearchFrom  string
)

func init() {
	f := errorSearchCmd.Flags()
	f.StringVar(&errorSearchTrack, "track", "trace", "Error source: trace, logs, rum")
	f.StringVar(&errorSearchFrom, "from", "24h", "Time range (e.g. 1h, 24h, 7d)")
	errorCmd.AddCommand(errorSearchCmd)
}

func runErrorSearch(cmd *cobra.Command, args []string) error {
	query := "*"
	if len(args) > 0 && args[0] != "" {
		query = args[0]
	}

	dur, err := parseDuration(errorSearchFrom)
	if err != nil {
		return fmt.Errorf("invalid --from duration %q: %w", errorSearchFrom, err)
	}

	now := time.Now()
	from := now.Add(-dur).UnixMilli()
	to := now.UnixMilli()

	track := datadogV2.IssuesSearchRequestDataAttributesTrack(errorSearchTrack)

	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV2.NewErrorTrackingApi(client)

	body := *datadogV2.NewIssuesSearchRequest(
		*datadogV2.NewIssuesSearchRequestData(
			datadogV2.IssuesSearchRequestDataAttributes{
				From:  from,
				To:    to,
				Query: query,
				Track: &track,
			},
			datadogV2.ISSUESSEARCHREQUESTDATATYPE_SEARCH_REQUEST,
		),
	)

	opts := datadogV2.NewSearchIssuesOptionalParameters().WithInclude([]datadogV2.SearchIssuesIncludeQueryParameterItem{
		datadogV2.SEARCHISSUESINCLUDEQUERYPARAMETERITEM_ISSUE,
	})
	resp, _, err := api.SearchIssues(ctx, body, *opts)
	if err != nil {
		return fmt.Errorf("failed to search errors: %s", ddclient.FormatAPIError(err))
	}

	if jsonOutput {
		return printJSON(resp)
	}

	results := resp.GetData()
	if len(results) == 0 {
		fmt.Println("No error issues found.")
		return nil
	}

	// Collect included issues for detail display
	issueMap := make(map[string]datadogV2.IssueAttributes)
	for _, inc := range resp.GetIncluded() {
		if inc.Issue != nil {
			issueMap[inc.Issue.Id] = inc.Issue.Attributes
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tERROR_TYPE\tMESSAGE\tCOUNT\tFIRST_SEEN")
	for _, r := range results {
		id := r.Id
		attrs := r.Attributes
		count := int64(0)
		if attrs.TotalCount != nil {
			count = *attrs.TotalCount
		}

		// Get issue details from included
		ia, ok := issueMap[id]
		errorType := ""
		errorMsg := ""
		firstSeen := ""
		if ok {
			errorType = derefStr(ia.ErrorType)
			errorMsg = derefStr(ia.ErrorMessage)
			if len(errorMsg) > 60 {
				errorMsg = errorMsg[:57] + "..."
			}
			if ia.FirstSeen != nil {
				firstSeen = time.UnixMilli(*ia.FirstSeen).Format("2006-01-02 15:04")
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\n", id, errorType, errorMsg, count, firstSeen)
	}
	return w.Flush()
}
