package cmd

import (
	"fmt"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"

	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var pipelineAddRemapperCmd = &cobra.Command{
	Use:   "add-remapper <pipeline-id>",
	Short: "Add an attribute remapper to a pipeline",
	Long: `Add an Attribute Remapper processor to an existing pipeline.

Examples:
  datadog pipeline add-remapper <id> \
    --source team --target team \
    --source-type attribute --target-type tag

  datadog pipeline add-remapper <id> \
    --source env --target env \
    --source-type attribute --target-type tag --preserve-source`,
	Args: cobra.ExactArgs(1),
	RunE: runPipelineAddRemapper,
}

var (
	remapperSource         string
	remapperTarget         string
	remapperSourceType     string
	remapperTargetType     string
	remapperPreserveSource bool
	remapperOverride       bool
	remapperName           string
)

func init() {
	f := pipelineAddRemapperCmd.Flags()
	f.StringVar(&remapperSource, "source", "", "Source attribute (required)")
	f.StringVar(&remapperTarget, "target", "", "Target attribute (required)")
	f.StringVar(&remapperSourceType, "source-type", "attribute", "Source type (attribute or tag)")
	f.StringVar(&remapperTargetType, "target-type", "tag", "Target type (attribute or tag)")
	f.BoolVar(&remapperPreserveSource, "preserve-source", false, "Keep source attribute after remapping")
	f.BoolVar(&remapperOverride, "override", false, "Override target if it already exists")
	f.StringVar(&remapperName, "processor-name", "Attribute Remapper", "Processor display name")
	pipelineAddRemapperCmd.MarkFlagRequired("source")
	pipelineAddRemapperCmd.MarkFlagRequired("target")
	pipelineCmd.AddCommand(pipelineAddRemapperCmd)
}

func runPipelineAddRemapper(cmd *cobra.Command, args []string) error {
	pipelineID := args[0]

	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV1.NewLogsPipelinesApi(client)

	pipeline, _, err := api.GetLogsPipeline(ctx, pipelineID)
	if err != nil {
		return fmt.Errorf("failed to get pipeline %s: %s", pipelineID, ddclient.FormatAPIError(err))
	}

	processor := datadogV1.LogsProcessor{
		LogsAttributeRemapper: &datadogV1.LogsAttributeRemapper{
			Sources:            []string{remapperSource},
			Target:             remapperTarget,
			SourceType:         &remapperSourceType,
			TargetType:         &remapperTargetType,
			PreserveSource:     datadog.PtrBool(remapperPreserveSource),
			OverrideOnConflict: datadog.PtrBool(remapperOverride),
			IsEnabled:          datadog.PtrBool(true),
			Name:               &remapperName,
			Type:               datadogV1.LOGSATTRIBUTEREMAPPERTYPE_ATTRIBUTE_REMAPPER,
		},
	}

	pipeline.Processors = append(pipeline.Processors, processor)

	_, _, err = api.UpdateLogsPipeline(ctx, pipelineID, pipeline)
	if err != nil {
		return fmt.Errorf("failed to update pipeline: %s", ddclient.FormatAPIError(err))
	}

	fmt.Printf("Attribute remapper added to pipeline %s (%s:%s → %s:%s).\n",
		pipelineID, remapperSourceType, remapperSource, remapperTargetType, remapperTarget)
	return nil
}
