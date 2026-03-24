package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"

	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var pipelineGetCmd = &cobra.Command{
	Use:   "get <pipeline-id>",
	Short: "Get pipeline details",
	Args:  cobra.ExactArgs(1),
	RunE:  runPipelineGet,
}

func init() {
	pipelineCmd.AddCommand(pipelineGetCmd)
}

func runPipelineGet(cmd *cobra.Command, args []string) error {
	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV1.NewLogsPipelinesApi(client)
	pipeline, _, err := api.GetLogsPipeline(ctx, args[0])
	if err != nil {
		return fmt.Errorf("failed to get pipeline: %s", ddclient.FormatAPIError(err))
	}

	if jsonOutput {
		return printJSON(pipeline)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID:\t%s\n", derefStr(pipeline.Id))
	fmt.Fprintf(w, "Name:\t%s\n", pipeline.Name)
	enabled := "no"
	if pipeline.IsEnabled != nil && *pipeline.IsEnabled {
		enabled = "yes"
	}
	fmt.Fprintf(w, "Enabled:\t%s\n", enabled)
	if pipeline.Filter != nil && pipeline.Filter.Query != nil {
		fmt.Fprintf(w, "Filter:\t%s\n", *pipeline.Filter.Query)
	}
	if err := w.Flush(); err != nil {
		return err
	}

	if len(pipeline.Processors) == 0 {
		fmt.Println("\nNo processors.")
		return nil
	}

	fmt.Printf("\nProcessors (%d):\n", len(pipeline.Processors))
	for i, p := range pipeline.Processors {
		printProcessor(i, p)
	}
	return nil
}

func printProcessor(index int, p datadogV1.LogsProcessor) {
	switch {
	case p.LogsCategoryProcessor != nil:
		cp := p.LogsCategoryProcessor
		enabled := cp.IsEnabled != nil && *cp.IsEnabled
		fmt.Printf("  [%d] CategoryProcessor (target=%s, enabled=%v)\n", index, cp.Target, enabled)
		for _, cat := range cp.Categories {
			query := ""
			if cat.Filter != nil && cat.Filter.Query != nil {
				query = *cat.Filter.Query
			}
			fmt.Printf("       %s → %s\n", query, derefStr(cat.Name))
		}

	case p.LogsAttributeRemapper != nil:
		ar := p.LogsAttributeRemapper
		enabled := ar.IsEnabled != nil && *ar.IsEnabled
		srcType := derefStr(ar.SourceType)
		tgtType := derefStr(ar.TargetType)
		fmt.Printf("  [%d] AttributeRemapper (sources=%v → target=%s, %s→%s, enabled=%v)\n",
			index, ar.Sources, ar.Target, srcType, tgtType, enabled)

	case p.LogsGrokParser != nil:
		gp := p.LogsGrokParser
		fmt.Printf("  [%d] GrokParser (source=%s, name=%s)\n", index, gp.Source, derefStr(gp.Name))

	case p.LogsDateRemapper != nil:
		fmt.Printf("  [%d] DateRemapper (sources=%v)\n", index, p.LogsDateRemapper.Sources)

	case p.LogsStatusRemapper != nil:
		fmt.Printf("  [%d] StatusRemapper (sources=%v)\n", index, p.LogsStatusRemapper.Sources)

	case p.LogsServiceRemapper != nil:
		fmt.Printf("  [%d] ServiceRemapper (sources=%v)\n", index, p.LogsServiceRemapper.Sources)

	case p.LogsMessageRemapper != nil:
		fmt.Printf("  [%d] MessageRemapper (sources=%v)\n", index, p.LogsMessageRemapper.Sources)

	case p.LogsURLParser != nil:
		fmt.Printf("  [%d] URLParser (sources=%v)\n", index, p.LogsURLParser.Sources)

	case p.LogsUserAgentParser != nil:
		fmt.Printf("  [%d] UserAgentParser (sources=%v)\n", index, p.LogsUserAgentParser.Sources)

	case p.LogsPipelineProcessor != nil:
		pp := p.LogsPipelineProcessor
		fmt.Printf("  [%d] NestedPipeline (name=%s, processors=%d)\n", index, derefStr(pp.Name), len(pp.Processors))

	case p.LogsTraceRemapper != nil:
		fmt.Printf("  [%d] TraceRemapper (sources=%v)\n", index, p.LogsTraceRemapper.Sources)

	case p.LogsLookupProcessor != nil:
		lp := p.LogsLookupProcessor
		fmt.Printf("  [%d] LookupProcessor (source=%s → target=%s)\n", index, lp.Source, lp.Target)

	case p.LogsArithmeticProcessor != nil:
		ap := p.LogsArithmeticProcessor
		fmt.Printf("  [%d] ArithmeticProcessor (expr=%s → target=%s)\n", index, ap.Expression, ap.Target)

	case p.LogsStringBuilderProcessor != nil:
		sb := p.LogsStringBuilderProcessor
		fmt.Printf("  [%d] StringBuilderProcessor (template=%s → target=%s)\n", index, sb.Template, sb.Target)

	case p.LogsGeoIPParser != nil:
		fmt.Printf("  [%d] GeoIPParser (sources=%v)\n", index, p.LogsGeoIPParser.Sources)

	default:
		fmt.Printf("  [%d] UnknownProcessor\n", index)
	}
}
