package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"

	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var monitorCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a monitor",
	Long: `Create a Datadog monitor. Supports standard metric monitors and log-based formula monitors.

For a standard metric monitor:
  datadog monitor create --type "query alert" --name "CPU high" \
    --query 'avg(last_5m):avg:system.cpu.user{env:prod} > 90' \
    --message "@pagerduty" --tags "env:prod,team:sre"

For a log-based error rate monitor (formula):
  datadog monitor create --type "query alert" --name "Lego error rate" \
    --formula "error_count / total_count" \
    --log-query-error "service:lego-v2-rest status:error" \
    --log-query-total "service:lego-v2-rest" \
    --threshold-critical 0.1 --threshold-warning 0.05 \
    --window "5m" \
    --message "@pagerduty-UC" --tags "env:prod,team:plat"`,
	RunE: runMonitorCreate,
}

var (
	createName              string
	createType              string
	createQuery             string
	createMessage           string
	createTags              string
	createThresholdCritical float64
	createThresholdWarning  float64
	// Formula log-based monitor fields
	createFormula       string
	createLogQueryError string
	createLogQueryTotal string
	createWindow        string
)

func init() {
	f := monitorCreateCmd.Flags()
	f.StringVar(&createName, "name", "", "Monitor name (required)")
	f.StringVar(&createType, "type", "query alert", "Monitor type")
	f.StringVar(&createQuery, "query", "", "Monitor query (for standard monitors)")
	f.StringVar(&createMessage, "message", "", "Notification message")
	f.StringVar(&createTags, "tags", "", "Comma-separated tags")
	f.Float64Var(&createThresholdCritical, "threshold-critical", 0, "Critical threshold")
	f.Float64Var(&createThresholdWarning, "threshold-warning", 0, "Warning threshold")
	// Formula fields
	f.StringVar(&createFormula, "formula", "", "Formula expression (e.g. 'error_count / total_count')")
	f.StringVar(&createLogQueryError, "log-query-error", "", "Log search query for numerator (error count)")
	f.StringVar(&createLogQueryTotal, "log-query-total", "", "Log search query for denominator (total count)")
	f.StringVar(&createWindow, "window", "5m", "Evaluation window (e.g. 5m, 15m, 1h)")

	monitorCreateCmd.MarkFlagRequired("name")
	monitorCmd.AddCommand(monitorCreateCmd)
}

func runMonitorCreate(cmd *cobra.Command, args []string) error {
	isFormula := createFormula != ""

	if !isFormula && createQuery == "" {
		return fmt.Errorf("either --query (standard) or --formula + --log-query-error + --log-query-total (formula) is required")
	}
	if isFormula && (createLogQueryError == "" || createLogQueryTotal == "") {
		return fmt.Errorf("--formula requires both --log-query-error and --log-query-total")
	}

	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	var monitor datadogV1.Monitor

	if isFormula {
		monitor, err = buildFormulaMonitor()
	} else {
		monitor, err = buildStandardMonitor()
	}
	if err != nil {
		return err
	}

	api := datadogV1.NewMonitorsApi(client)
	resp, _, err := api.CreateMonitor(ctx, monitor)
	if err != nil {
		return fmt.Errorf("failed to create monitor: %s", ddclient.FormatAPIError(err))
	}

	fmt.Printf("Monitor %d created successfully.\n", derefInt64(resp.Id))

	if jsonOutput {
		return printJSON(resp)
	}
	return nil
}

func buildStandardMonitor() (datadogV1.Monitor, error) {
	monitorType, err := parseMonitorType(createType)
	if err != nil {
		return datadogV1.Monitor{}, err
	}

	m := datadogV1.Monitor{
		Name:    &createName,
		Type:    monitorType,
		Query:   createQuery,
		Message: &createMessage,
		Tags:    parseTags(createTags),
		Options: &datadogV1.MonitorOptions{
			Thresholds:  buildThresholds(),
			IncludeTags: datadog.PtrBool(true),
		},
	}
	return m, nil
}

func buildFormulaMonitor() (datadogV1.Monitor, error) {
	monitorType, err := parseMonitorType(createType)
	if err != nil {
		return datadogV1.Monitor{}, err
	}

	query := fmt.Sprintf(`formula("%s").last("%s") > %v`, createFormula, createWindow, createThresholdCritical)

	// Parse formula to extract variable names
	// e.g. "error_count / total_count" → names are "error_count" and "total_count"
	parts := strings.Fields(createFormula)
	if len(parts) != 3 || parts[1] != "/" {
		return datadogV1.Monitor{}, fmt.Errorf("formula must be in format 'var1 / var2', got: %s", createFormula)
	}
	numeratorName := parts[0]
	denominatorName := parts[2]

	variables := []datadogV1.MonitorFormulaAndFunctionQueryDefinition{
		{
			MonitorFormulaAndFunctionEventQueryDefinition: &datadogV1.MonitorFormulaAndFunctionEventQueryDefinition{
				DataSource: datadogV1.MONITORFORMULAANDFUNCTIONEVENTSDATASOURCE_LOGS,
				Name:       numeratorName,
				Search: &datadogV1.MonitorFormulaAndFunctionEventQueryDefinitionSearch{
					Query: createLogQueryError,
				},
				Indexes: []string{"*"},
				Compute: datadogV1.MonitorFormulaAndFunctionEventQueryDefinitionCompute{
					Aggregation: datadogV1.MONITORFORMULAANDFUNCTIONEVENTAGGREGATION_COUNT,
				},
				GroupBy: []datadogV1.MonitorFormulaAndFunctionEventQueryGroupBy{},
			},
		},
		{
			MonitorFormulaAndFunctionEventQueryDefinition: &datadogV1.MonitorFormulaAndFunctionEventQueryDefinition{
				DataSource: datadogV1.MONITORFORMULAANDFUNCTIONEVENTSDATASOURCE_LOGS,
				Name:       denominatorName,
				Search: &datadogV1.MonitorFormulaAndFunctionEventQueryDefinitionSearch{
					Query: createLogQueryTotal,
				},
				Indexes: []string{"*"},
				Compute: datadogV1.MonitorFormulaAndFunctionEventQueryDefinitionCompute{
					Aggregation: datadogV1.MONITORFORMULAANDFUNCTIONEVENTAGGREGATION_COUNT,
				},
				GroupBy: []datadogV1.MonitorFormulaAndFunctionEventQueryGroupBy{},
			},
		},
	}

	m := datadogV1.Monitor{
		Name:    &createName,
		Type:    monitorType,
		Query:   query,
		Message: &createMessage,
		Tags:    parseTags(createTags),
		Options: &datadogV1.MonitorOptions{
			Thresholds:      buildThresholds(),
			IncludeTags:     datadog.PtrBool(true),
			Variables:       variables,
			RequireFullWindow: datadog.PtrBool(true),
		},
	}
	return m, nil
}

func buildThresholds() *datadogV1.MonitorThresholds {
	t := &datadogV1.MonitorThresholds{}
	if createThresholdCritical != 0 {
		t.Critical = &createThresholdCritical
	}
	if createThresholdWarning != 0 {
		t.Warning = *datadog.NewNullableFloat64(&createThresholdWarning)
	}
	return t
}

func parseTags(s string) []string {
	if s == "" {
		return nil
	}
	tags := strings.Split(s, ",")
	for i := range tags {
		tags[i] = strings.TrimSpace(tags[i])
	}
	return tags
}

func parseMonitorType(s string) (datadogV1.MonitorType, error) {
	var t datadogV1.MonitorType
	err := json.Unmarshal([]byte(`"`+s+`"`), &t)
	if err != nil {
		return t, fmt.Errorf("invalid monitor type %q", s)
	}
	return t, nil
}
