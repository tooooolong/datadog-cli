package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"

	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var pipelineListCmd = &cobra.Command{
	Use:   "list",
	Short: "List log pipelines",
	RunE:  runPipelineList,
}

func init() {
	pipelineCmd.AddCommand(pipelineListCmd)
}

func runPipelineList(cmd *cobra.Command, args []string) error {
	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV1.NewLogsPipelinesApi(client)
	pipelines, _, err := api.ListLogsPipelines(ctx)
	if err != nil {
		return fmt.Errorf("failed to list pipelines: %s", ddclient.FormatAPIError(err))
	}

	if jsonOutput {
		return printJSON(pipelines)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tENABLED\tFILTER\tPROCESSORS")
	for _, p := range pipelines {
		id := derefStr(p.Id)
		enabled := "no"
		if p.IsEnabled != nil && *p.IsEnabled {
			enabled = "yes"
		}
		filter := ""
		if p.Filter != nil && p.Filter.Query != nil {
			filter = *p.Filter.Query
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n", id, p.Name, enabled, filter, len(p.Processors))
	}
	return w.Flush()
}
