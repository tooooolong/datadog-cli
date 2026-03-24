package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/tooooolong/datadog-cli/internal/ddclient"
	"github.com/spf13/cobra"
)

var monitorUpdateCmd = &cobra.Command{
	Use:   "update <monitor-id>",
	Short: "Update a monitor",
	Long: `Update a Datadog monitor's properties.

Examples:
  datadog monitor update 12345 --name "New Name"
  datadog monitor update 12345 --tags "env:prod,team:infra"
  datadog monitor update 12345 --query "avg(...) > 100" --threshold-critical 100
  datadog monitor update 12345 --threshold-warning 80 --renotify 30`,
	Args: cobra.ExactArgs(1),
	RunE: runMonitorUpdate,
}

var (
	updateName              string
	updateTags              string
	updateMessage           string
	updateQuery             string
	updatePriority          int64
	updateThresholdCritical         float64
	updateThresholdWarning          float64
	updateThresholdCriticalRecovery float64
	updateThresholdWarningRecovery  float64
	updateRenotify                  int64
	updateOnMissingData             string
)

func init() {
	f := monitorUpdateCmd.Flags()
	f.StringVar(&updateName, "name", "", "New monitor name")
	f.StringVar(&updateTags, "tags", "", "Comma-separated tags")
	f.StringVar(&updateMessage, "message", "", "New monitor message")
	f.StringVar(&updateQuery, "query", "", "New monitor query")
	f.Int64Var(&updatePriority, "priority", 0, "Priority (1-5, 0 to clear)")
	f.Float64Var(&updateThresholdCritical, "threshold-critical", 0, "Critical threshold")
	f.Float64Var(&updateThresholdWarning, "threshold-warning", 0, "Warning threshold (0 to clear)")
	f.Float64Var(&updateThresholdCriticalRecovery, "threshold-critical-recovery", 0, "Critical recovery threshold (0 to clear)")
	f.Float64Var(&updateThresholdWarningRecovery, "threshold-warning-recovery", 0, "Warning recovery threshold (0 to clear)")
	f.Int64Var(&updateRenotify, "renotify", -1, "Re-notification interval in minutes (0 to disable)")
	f.StringVar(&updateOnMissingData, "on-missing-data", "", "Behavior on missing data: default, show_no_data, show_and_notify_no_data, resolve")
	monitorCmd.AddCommand(monitorUpdateCmd)
}

func runMonitorUpdate(cmd *cobra.Command, args []string) error {
	monitorID, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid monitor ID %q: %w", args[0], err)
	}

	ctx, client, err := ddclient.NewClient(ddSite)
	if err != nil {
		return err
	}

	api := datadogV1.NewMonitorsApi(client)

	current, _, err := api.GetMonitor(ctx, monitorID)
	if err != nil {
		return fmt.Errorf("failed to get monitor %d: %w", monitorID, err)
	}

	body := datadogV1.NewMonitorUpdateRequest()
	body.SetType(current.Type)
	body.SetQuery(current.Query)

	changed := false

	if cmd.Flags().Changed("name") {
		body.SetName(updateName)
		changed = true
	} else {
		body.Name = current.Name
	}

	if cmd.Flags().Changed("message") {
		body.SetMessage(updateMessage)
		changed = true
	} else {
		body.Message = current.Message
	}

	if cmd.Flags().Changed("query") {
		body.SetQuery(updateQuery)
		changed = true
	}

	if cmd.Flags().Changed("tags") {
		tags := strings.Split(updateTags, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
		body.SetTags(tags)
		changed = true
	} else {
		body.SetTags(current.Tags)
	}

	if cmd.Flags().Changed("priority") {
		if updatePriority == 0 {
			body.Priority = *datadog.NewNullableInt64(nil)
		} else {
			body.SetPriority(updatePriority)
		}
		changed = true
	} else {
		body.Priority = current.Priority
	}

	// Copy current options and apply overrides
	opts := datadogV1.MonitorOptions{}
	if current.Options != nil {
		opts = datadogV1.MonitorOptions(*current.Options)
	}

	thresholdChanged := cmd.Flags().Changed("threshold-critical") || cmd.Flags().Changed("threshold-warning") ||
		cmd.Flags().Changed("threshold-critical-recovery") || cmd.Flags().Changed("threshold-warning-recovery")
	if thresholdChanged {
		if opts.Thresholds == nil {
			opts.Thresholds = &datadogV1.MonitorThresholds{}
		}
		if cmd.Flags().Changed("threshold-critical") {
			opts.Thresholds.Critical = &updateThresholdCritical
		}
		if cmd.Flags().Changed("threshold-warning") {
			if updateThresholdWarning == 0 {
				opts.Thresholds.Warning = *datadog.NewNullableFloat64(nil)
			} else {
				opts.Thresholds.Warning = *datadog.NewNullableFloat64(&updateThresholdWarning)
			}
		}
		if cmd.Flags().Changed("threshold-critical-recovery") {
			if updateThresholdCriticalRecovery == 0 {
				opts.Thresholds.CriticalRecovery = *datadog.NewNullableFloat64(nil)
			} else {
				opts.Thresholds.CriticalRecovery = *datadog.NewNullableFloat64(&updateThresholdCriticalRecovery)
			}
		}
		if cmd.Flags().Changed("threshold-warning-recovery") {
			if updateThresholdWarningRecovery == 0 {
				opts.Thresholds.WarningRecovery = *datadog.NewNullableFloat64(nil)
			} else {
				opts.Thresholds.WarningRecovery = *datadog.NewNullableFloat64(&updateThresholdWarningRecovery)
			}
		}
		changed = true
	}

	if cmd.Flags().Changed("renotify") {
		if updateRenotify < 0 {
			return fmt.Errorf("--renotify must be >= 0 (0 to disable)")
		}
		opts.RenotifyInterval = *datadog.NewNullableInt64(&updateRenotify)
		changed = true
	}

	if cmd.Flags().Changed("on-missing-data") {
		v := datadogV1.OnMissingDataOption(updateOnMissingData)
		opts.OnMissingData = &v
		changed = true
	}

	body.SetOptions(opts)

	if !changed {
		return fmt.Errorf("no update flags specified; use --name, --tags, --message, --query, --priority, --threshold-critical, --threshold-warning, or --renotify")
	}

	updated, _, err := api.UpdateMonitor(ctx, monitorID, *body)
	if err != nil {
		return fmt.Errorf("failed to update monitor %d: %s", monitorID, ddclient.FormatAPIError(err))
	}

	fmt.Printf("Monitor %d updated successfully.\n", monitorID)

	if jsonOutput {
		return printJSON(updated)
	}
	return nil
}
