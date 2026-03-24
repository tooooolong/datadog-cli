# datadog-cli

A command-line tool for managing Datadog resources — monitors, metrics, events, logs, services, errors, and log pipelines.

Built with Go, [cobra](https://github.com/spf13/cobra), and [datadog-api-client-go](https://github.com/DataDog/datadog-api-client-go).

> **Note:** This tool covers a subset of commonly used Datadog APIs, not the full platform. Currently supported: Monitors, Metrics Query, Events, Logs Search, APM Services, Error Tracking, and Log Pipelines. Features like Dashboards, Synthetics, SLOs, Notebooks, Downtimes, and Incidents are not yet implemented. Contributions welcome.

## Skill

```bash
npx skills add tooooolong/datadog-cli@datadog-cli-skill
```

Installs the `datadog-cli` skill for your agent, enabling AI-assisted Datadog operations.

## Install

```bash
go install github.com/tooooolong/datadog-cli@latest
```

Or build from source:

```bash
git clone https://github.com/tooooolong/datadog-cli.git
cd datadog-cli
make install   # installs to ~/.local/bin/datadog
```

## Setup

```bash
datadog login
```

Credentials are saved to `~/.local/config/datadog-cli/config.json`. Environment variables `DD_API_KEY` / `DD_APP_KEY` take priority over the saved config if set.

### Required App Key Scopes

| Command | Scope |
|---------|-------|
| `monitor` | `monitors_read`, `monitors_write` |
| `metric query` | `timeseries_query` |
| `event list` | `events_read` |
| `log search` | `logs_read_data` |
| `service list` | `apm_read` |
| `error search/get` | `error_tracking_read` |
| `pipeline` | `logs_pipelines_read`, `logs_pipelines_write` |

## Commands

```
datadog
├── monitor                     Manage monitors
│   ├── list                    List with filters (--tags, --name, --page)
│   ├── get <id>                Get details
│   ├── search <query>          Search (--page, --per-page, --sort)
│   ├── create                  Create standard or log-based formula monitor
│   ├── update <id>             Update properties, thresholds, renotify
│   └── delete <id>             Delete
│
├── metric
│   └── query <query>           Query scalar value (--from)
│
├── event
│   └── list                    List events (--query, --from, --limit)
│
├── log
│   └── search <query>          Search logs (--from, --limit)
│
├── service
│   └── list                    List APM services (--env)
│
├── error (alias: errors)
│   ├── search [query]          Search Error Tracking issues (--track, --from)
│   └── get <id>                Get issue details
│
└── pipeline
    ├── list                    List pipelines
    ├── get <id>                Get with processor details
    ├── create                  Create pipeline (--name, --filter, --enabled)
    ├── delete <id>             Delete pipeline
    ├── add-category <id>       Add Category Processor (--target, --rule)
    └── add-remapper <id>       Add Attribute Remapper (--source, --target)
```

### Global Flags

```
--json          Output in JSON format (default: table)
--site string   Datadog site (default: datadoghq.com)
```

## Examples

```bash
# List monitors filtered by tag
datadog monitor list --tags "env:prod" --page-size 10

# Check a metric value to validate a monitor threshold
datadog metric query "avg:system.cpu.user{env:prod}" --from 1h

# Search error logs from a service
datadog log search "service:myapp status:error" --from 1h

# View top Error Tracking issues
datadog error search --track trace --from 24h

# List active APM services
datadog service list --env prod

# Create a log-based error rate monitor
datadog monitor create --type "log alert" \
  --name "My service errors" \
  --formula "errors / total" \
  --log-query-error "service:myapp status:error" \
  --log-query-total "service:myapp" \
  --threshold-critical 0.1 --window "5m" \
  --message "@pagerduty" --tags "env:prod"

# Set up a pipeline to tag logs by team
datadog pipeline create --name "Team Tagging" --filter "*" --enabled
datadog pipeline add-category <id> --target team \
  --rule "service:(svc-a OR svc-b)=team-alpha" \
  --rule "service:(svc-c)=team-beta"
datadog pipeline add-remapper <id> \
  --source team --target team \
  --source-type attribute --target-type tag
```

## License

MIT
