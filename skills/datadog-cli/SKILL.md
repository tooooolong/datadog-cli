---
name: datadog-cli
description: Use the `datadog` CLI to manage Datadog resources — monitors, metrics, events, logs, services, errors, and pipelines. Invoke this skill whenever the user asks to query, create, update, or delete Datadog monitors, search logs or errors, check metric values, list APM services, or manage log pipelines. Also trigger when the user mentions Datadog observability tasks like "check the error rate", "look at monitors", "search logs for errors", "list services", or "set up a log pipeline".
---

# datadog CLI

A command-line tool for interacting with the Datadog API. Built with Go, cobra, and `datadog-api-client-go/v2`.

## Prerequisites

Set environment variables before use:

```bash
export DD_API_KEY="your-api-key"
export DD_APP_KEY="your-app-key"
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output in JSON format (default: table) |
| `--site <site>` | Datadog site (default: `datadoghq.com`, e.g. `datadoghq.eu`, `us5.datadoghq.com`) |

## Command Map

```
datadog
├── monitor        Manage monitors (alerts)
│   ├── list       List monitors with filters
│   ├── get        Get monitor details by ID
│   ├── search     Search monitors by query
│   ├── create     Create a monitor (standard or log-based formula)
│   ├── update     Update monitor properties and thresholds
│   └── delete     Delete a monitor
│
├── metric         Query metric data
│   └── query      Query scalar metric values
│
├── event          Query events
│   └── list       List events with time range and filters
│
├── log            Query logs
│   └── search     Search logs by query
│
├── service        Query APM services
│   └── list       List active services by environment
│
├── error          Error Tracking (aliases: error, errors)
│   ├── search     Search error issues
│   └── get        Get error issue details
│
└── pipeline       Manage log processing pipelines
    ├── list       List all pipelines
    ├── get        Get pipeline details with processors
    ├── create     Create a new pipeline
    ├── delete     Delete a pipeline
    ├── add-category   Add a Category Processor
    └── add-remapper   Add an Attribute Remapper
```

## Quick Reference

For detailed usage of each command group, read the corresponding reference file:

| Task | Command | Reference |
|------|---------|-----------|
| List/search/create/update/delete monitors | `datadog monitor ...` | `references/monitor.md` |
| Query metric values to validate thresholds | `datadog metric query ...` | `references/metric.md` |
| Search events and alert history | `datadog event list ...` | `references/event.md` |
| Search logs for errors or patterns | `datadog log search ...` | `references/log.md` |
| List active APM services | `datadog service list ...` | `references/service.md` |
| Search/view Error Tracking issues | `datadog error ...` | `references/error.md` |
| Manage log pipelines and processors | `datadog pipeline ...` | `references/pipeline.md` |

## Common Workflows

### Investigate a service's health

```bash
datadog service list --env prod                    # confirm service is active
datadog log search "service:myapp status:error" --from 1h
datadog error search "service:myapp" --track trace --from 24h
datadog metric query "avg:trace.http.request.errors{service:myapp}.as_count()" --from 1h
```

### Audit and fix monitors

```bash
datadog monitor list --json                        # export all monitors
datadog monitor get <id>                           # inspect a specific one
datadog metric query "<monitor-query>" --from 1h   # validate threshold against actual values
datadog monitor update <id> --threshold-critical 100 --threshold-warning 80 --renotify 30
```

### Set up team-based log tagging

```bash
datadog pipeline create --name "Team Tagging" --filter "*" --enabled
datadog pipeline add-category <id> --target team \
  --rule "service:(svc-a OR svc-b)=team-alpha" \
  --rule "service:(svc-c OR svc-d)=team-beta"
datadog pipeline add-remapper <id> --source team --target team \
  --source-type attribute --target-type tag
```

### Create a log-based error rate monitor

```bash
datadog monitor create --type "log alert" \
  --name "My service error rate" \
  --formula "errors / total" \
  --log-query-error "service:myapp status:error" \
  --log-query-total "service:myapp" \
  --threshold-critical 0.1 --threshold-warning 0.05 \
  --window "5m" --message "@pagerduty-team" --tags "env:prod,team:myteam"
```
