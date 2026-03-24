# datadog monitor

Manage Datadog monitors (alerts).

## list

List monitors with optional filters.

```bash
datadog monitor list [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--tags` | | Comma-separated tags to filter (e.g. `env:prod,team:spot`) |
| `--name` | | Filter by monitor name pattern |
| `--page` | `0` | Page number (0-indexed) |
| `--page-size` | `50` | Results per page |

**Output columns:** ID, NAME, TYPE, STATUS, TAGS

```bash
datadog monitor list --tags "env:prod" --page-size 10
datadog monitor list --name "CPU" --json
```

## get

Get detailed information about a specific monitor.

```bash
datadog monitor get <monitor-id>
```

Shows: ID, Name, Type, Query, Message, Status, Priority, Tags, Created, Modified.

```bash
datadog monitor get 12345678
datadog monitor get 12345678 --json    # full JSON with options/thresholds
```

## search

Search monitors using Datadog search query syntax.

```bash
datadog monitor search <query> [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--page` | `0` | Page number |
| `--per-page` | `50` | Results per page |
| `--sort` | | Sort field (e.g. `name`, `status`) |

```bash
datadog monitor search "type:metric status:alert"
datadog monitor search "tag:env:prod"
```

## create

Create a new monitor. Supports two modes:

### Standard monitor (--query)

```bash
datadog monitor create --name "CPU high" --type "query alert" \
  --query 'avg(last_5m):avg:system.cpu.user{env:prod} > 90' \
  --threshold-critical 90 --threshold-warning 80 \
  --message "@pagerduty" --tags "env:prod,team:sre"
```

### Log-based formula monitor (--formula)

For error rate monitoring when APM is unavailable:

```bash
datadog monitor create --name "Error rate" --type "log alert" \
  --formula "errors / total" \
  --log-query-error "service:myapp status:error" \
  --log-query-total "service:myapp" \
  --threshold-critical 0.1 --threshold-warning 0.05 \
  --window "5m" --message "@pagerduty" --tags "env:prod"
```

| Flag | Description |
|------|-------------|
| `--name` | Monitor name (required) |
| `--type` | Monitor type (default: `query alert`) |
| `--query` | Standard monitor query |
| `--formula` | Formula expression (e.g. `errors / total`) |
| `--log-query-error` | Log query for numerator |
| `--log-query-total` | Log query for denominator |
| `--threshold-critical` | Critical threshold |
| `--threshold-warning` | Warning threshold |
| `--window` | Evaluation window (default: `5m`) |
| `--message` | Notification message |
| `--tags` | Comma-separated tags |

## update

Update a monitor's properties. Only specified flags are changed; others preserved.

```bash
datadog monitor update <monitor-id> [flags]
```

| Flag | Description |
|------|-------------|
| `--name` | New name |
| `--query` | New query |
| `--message` | New notification message |
| `--tags` | Comma-separated tags |
| `--priority` | Priority 1-5 (0 to clear) |
| `--threshold-critical` | Critical threshold |
| `--threshold-warning` | Warning threshold (0 to clear) |
| `--threshold-critical-recovery` | Critical recovery threshold (0 to clear) |
| `--threshold-warning-recovery` | Warning recovery threshold (0 to clear) |
| `--renotify` | Re-notification interval in minutes (0 to disable) |
| `--on-missing-data` | `default`, `show_no_data`, `show_and_notify_no_data`, `resolve` |

```bash
datadog monitor update 12345 --threshold-critical 100 --threshold-warning 80
datadog monitor update 12345 --renotify 30 --tags "env:prod,team:spot"
datadog monitor update 12345 --on-missing-data resolve
```

## delete

Delete a monitor by ID.

```bash
datadog monitor delete <monitor-id>
```
