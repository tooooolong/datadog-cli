package cmd

import (
	"fmt"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"

	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var pipelineAddCategoryCmd = &cobra.Command{
	Use:   "add-category <pipeline-id>",
	Short: "Add a category processor to a pipeline",
	Long: `Add a Category Processor to an existing pipeline.

Rules use the format "query=category_value".

Examples:
  datadog pipeline add-category <id> --target team \
    --rule "service:(counter-rest OR nagoya OR panamera-rest)=plat" \
    --rule "service:(smaug OR engine OR strategy-rpc)=spot" \
    --rule "service:(contract-api OR contract-ws)=contract"`,
	Args: cobra.ExactArgs(1),
	RunE: runPipelineAddCategory,
}

var (
	categoryTarget string
	categoryRules  []string
	categoryName   string
)

func init() {
	f := pipelineAddCategoryCmd.Flags()
	f.StringVar(&categoryTarget, "target", "", "Target attribute name (required)")
	f.StringArrayVar(&categoryRules, "rule", nil, "Category rule in format 'query=value' (repeatable)")
	f.StringVar(&categoryName, "processor-name", "Category Processor", "Processor display name")
	pipelineAddCategoryCmd.MarkFlagRequired("target")
	pipelineAddCategoryCmd.MarkFlagRequired("rule")
	pipelineCmd.AddCommand(pipelineAddCategoryCmd)
}

func runPipelineAddCategory(cmd *cobra.Command, args []string) error {
	pipelineID := args[0]

	categories, err := parseRules(categoryRules)
	if err != nil {
		return err
	}

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
		LogsCategoryProcessor: &datadogV1.LogsCategoryProcessor{
			Target:     categoryTarget,
			Categories: categories,
			IsEnabled:  datadog.PtrBool(true),
			Name:       &categoryName,
			Type:       datadogV1.LOGSCATEGORYPROCESSORTYPE_CATEGORY_PROCESSOR,
		},
	}

	pipeline.Processors = append(pipeline.Processors, processor)

	_, _, err = api.UpdateLogsPipeline(ctx, pipelineID, pipeline)
	if err != nil {
		return fmt.Errorf("failed to update pipeline: %s", ddclient.FormatAPIError(err))
	}

	fmt.Printf("Category processor added to pipeline %s (%d rules, target=%s).\n", pipelineID, len(categories), categoryTarget)
	return nil
}

func parseRules(rules []string) ([]datadogV1.LogsCategoryProcessorCategory, error) {
	var categories []datadogV1.LogsCategoryProcessorCategory
	for _, rule := range rules {
		idx := strings.LastIndex(rule, "=")
		if idx <= 0 {
			return nil, fmt.Errorf("invalid rule format %q, expected 'query=value'", rule)
		}
		query := rule[:idx]
		value := rule[idx+1:]
		categories = append(categories, datadogV1.LogsCategoryProcessorCategory{
			Filter: &datadogV1.LogsFilter{
				Query: &query,
			},
			Name: &value,
		})
	}
	return categories, nil
}
