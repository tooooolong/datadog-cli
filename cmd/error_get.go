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

var errorGetCmd = &cobra.Command{
	Use:   "get <issue-id>",
	Short: "Get error issue details",
	Args:  cobra.ExactArgs(1),
	RunE:  runErrorGet,
}

func init() {
	errorCmd.AddCommand(errorGetCmd)
}

func runErrorGet(cmd *cobra.Command, args []string) error {
	issueID := args[0]

	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV2.NewErrorTrackingApi(client)
	include := []datadogV2.GetIssueIncludeQueryParameterItem{
		datadogV2.GETISSUEINCLUDEQUERYPARAMETERITEM_ASSIGNEE,
	}
	opts := datadogV2.NewGetIssueOptionalParameters().WithInclude(include)

	resp, _, err := api.GetIssue(ctx, issueID, *opts)
	if err != nil {
		return fmt.Errorf("failed to get issue %s: %s", issueID, ddclient.FormatAPIError(err))
	}

	if jsonOutput {
		return printJSON(resp)
	}

	issue := resp.GetData()
	attrs := issue.GetAttributes()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID:\t%s\n", issue.Id)
	fmt.Fprintf(w, "Error Type:\t%s\n", derefStr(attrs.ErrorType))
	fmt.Fprintf(w, "Message:\t%s\n", derefStr(attrs.ErrorMessage))
	if attrs.State != nil {
		fmt.Fprintf(w, "State:\t%s\n", *attrs.State)
	}
	if attrs.FilePath != nil {
		fmt.Fprintf(w, "File:\t%s\n", *attrs.FilePath)
	}
	if attrs.FunctionName != nil {
		fmt.Fprintf(w, "Function:\t%s\n", *attrs.FunctionName)
	}
	if attrs.Languages != nil {
		langs := make([]string, len(attrs.Languages))
		for i, l := range attrs.Languages {
			langs[i] = string(l)
		}
		fmt.Fprintf(w, "Languages:\t%s\n", strings.Join(langs, ", "))
	}
	if attrs.FirstSeen != nil {
		fmt.Fprintf(w, "First Seen:\t%s\n", time.UnixMilli(*attrs.FirstSeen).Format("2006-01-02 15:04:05"))
	}
	if attrs.LastSeen != nil {
		fmt.Fprintf(w, "Last Seen:\t%s\n", time.UnixMilli(*attrs.LastSeen).Format("2006-01-02 15:04:05"))
	}
	if attrs.FirstSeenVersion != nil {
		fmt.Fprintf(w, "First Version:\t%s\n", *attrs.FirstSeenVersion)
	}
	if attrs.LastSeenVersion != nil {
		fmt.Fprintf(w, "Last Version:\t%s\n", *attrs.LastSeenVersion)
	}
	return w.Flush()
}
