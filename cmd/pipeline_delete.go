package cmd

import (
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"

	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var pipelineDeleteCmd = &cobra.Command{
	Use:   "delete <pipeline-id>",
	Short: "Delete a log pipeline",
	Args:  cobra.ExactArgs(1),
	RunE:  runPipelineDelete,
}

func init() {
	pipelineCmd.AddCommand(pipelineDeleteCmd)
}

func runPipelineDelete(cmd *cobra.Command, args []string) error {
	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV1.NewLogsPipelinesApi(client)
	_, err = api.DeleteLogsPipeline(ctx, args[0])
	if err != nil {
		return fmt.Errorf("failed to delete pipeline %s: %s", args[0], ddclient.FormatAPIError(err))
	}

	fmt.Printf("Pipeline %s deleted successfully.\n", args[0])
	return nil
}
