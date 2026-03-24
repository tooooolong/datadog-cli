package cmd

import (
	"fmt"
	"math"
	"os"
	"text/tabwriter"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"

	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var metricQueryCmd = &cobra.Command{
	Use:   "query <metric-query>",
	Short: "Query metric timeseries data",
	Long: `Query Datadog metrics and display recent values.

The query uses the same syntax as Datadog metric queries.

Examples:
  datadog metric query "avg:system.cpu.user{env:prod}"
  datadog metric query "avg:holly.balance_state.delay{env:prod}" --from 1h
  datadog metric query "sum:trace.rack.request.hits{service:panamera-rest}.as_rate()" --from 30m`,
	Args: cobra.ExactArgs(1),
	RunE: runMetricQuery,
}

var (
	metricFrom string
)

func init() {
	metricQueryCmd.Flags().StringVar(&metricFrom, "from", "1h", "Time range to query (e.g. 15m, 1h, 4h, 1d)")
	metricCmd.AddCommand(metricQueryCmd)
}

func parseDuration(s string) (time.Duration, error) {
	if len(s) > 0 && s[len(s)-1] == 'd' {
		s = s[:len(s)-1] + "h"
		d, err := time.ParseDuration(s)
		if err != nil {
			return 0, err
		}
		return d * 24, nil
	}
	return time.ParseDuration(s)
}

func runMetricQuery(cmd *cobra.Command, args []string) error {
	query := args[0]

	dur, err := parseDuration(metricFrom)
	if err != nil {
		return fmt.Errorf("invalid --from duration %q: %w", metricFrom, err)
	}

	now := time.Now()
	from := now.Add(-dur).UnixMilli()
	to := now.UnixMilli()

	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV2.NewMetricsApi(client)

	body := *datadogV2.NewScalarFormulaQueryRequest(
		*datadogV2.NewScalarFormulaRequest(
			*datadogV2.NewScalarFormulaRequestAttributes(from, []datadogV2.ScalarQuery{
				datadogV2.MetricsScalarQueryAsScalarQuery(
					&datadogV2.MetricsScalarQuery{
						Aggregator: datadogV2.METRICSAGGREGATOR_AVG,
						DataSource: datadogV2.METRICSDATASOURCE_METRICS,
						Query:      query,
						Name:       datadog.PtrString("a"),
					},
				),
			}, to),
			datadogV2.SCALARFORMULAREQUESTTYPE_SCALAR_REQUEST,
		),
	)

	resp, _, err := api.QueryScalarData(ctx, body)
	if err != nil {
		return fmt.Errorf("failed to query metrics: %s", ddclient.FormatAPIError(err))
	}

	if jsonOutput {
		return printJSON(resp)
	}

	data := resp.GetData()
	attrs := data.GetAttributes()
	columns := attrs.GetColumns()

	if len(columns) == 0 {
		fmt.Println("No data returned.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Query:\t%s\n", query)
	fmt.Fprintf(w, "Range:\t%s ago → now\n", metricFrom)

	for _, col := range columns {
		if col.DataScalarColumn != nil {
			dc := col.DataScalarColumn
			values := dc.GetValues()
			name := dc.GetName()
			fmt.Fprintf(w, "Name:\t%s\n", name)
			for _, v := range values {
				if v != nil {
					fmt.Fprintf(w, "Value:\t%s\n", formatFloat(*v))
				}
			}
		}
	}
	return w.Flush()
}

func formatFloat(f float64) string {
	if f == float64(int64(f)) && math.Abs(f) < 1e15 {
		return fmt.Sprintf("%.0f", f)
	}
	if math.Abs(f) >= 1000 {
		return fmt.Sprintf("%.2f", f)
	}
	return fmt.Sprintf("%.4f", f)
}
