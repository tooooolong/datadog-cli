package cmd

import (
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"

	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var pipelineCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a log pipeline",
	Long: `Create a new Datadog log processing pipeline.

Examples:
  datadog pipeline create --name "Team Tagging" --filter "*" --enabled
  datadog pipeline create --name "UC Logs" --filter "service:(counter-rest OR panamera-rest)"`,
	RunE: runPipelineCreate,
}

var (
	pipelineCreateName    string
	pipelineCreateFilter  string
	pipelineCreateEnabled bool
)

func init() {
	f := pipelineCreateCmd.Flags()
	f.StringVar(&pipelineCreateName, "name", "", "Pipeline name (required)")
	f.StringVar(&pipelineCreateFilter, "filter", "", "Log filter query")
	f.BoolVar(&pipelineCreateEnabled, "enabled", false, "Enable pipeline immediately")
	pipelineCreateCmd.MarkFlagRequired("name")
	pipelineCmd.AddCommand(pipelineCreateCmd)
}

func runPipelineCreate(cmd *cobra.Command, args []string) error {
	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	body := datadogV1.LogsPipeline{
		Name:      pipelineCreateName,
		IsEnabled: datadog.PtrBool(pipelineCreateEnabled),
	}
	if pipelineCreateFilter != "" {
		body.Filter = &datadogV1.LogsFilter{
			Query: &pipelineCreateFilter,
		}
	}

	api := datadogV1.NewLogsPipelinesApi(client)
	resp, _, err := api.CreateLogsPipeline(ctx, body)
	if err != nil {
		return fmt.Errorf("failed to create pipeline: %s", ddclient.FormatAPIError(err))
	}

	fmt.Printf("Pipeline %s created successfully (id: %s).\n", resp.Name, derefStr(resp.Id))

	if jsonOutput {
		return printJSON(resp)
	}
	return nil
}
