# datadog error

Query Datadog Error Tracking issues. Alias: `datadog errors`.

## search

Search Error Tracking issues across traces, logs, or RUM.

```bash
datadog error search [query] [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--track` | `trace` | Error source: `trace`, `logs`, `rum` |
| `--from` | `24h` | Time range (e.g. `1h`, `24h`, `7d`) |

If no query is provided, defaults to `*` (all issues).

**Query syntax:**

| Query | Description |
|-------|-------------|
| `service:myapp` | Errors from a specific service |
| `error.type:RuntimeError` | Errors by type |
| `*` | All error issues |

**Output columns:** ID, ERROR_TYPE, MESSAGE, COUNT, FIRST_SEEN

```bash
datadog error search --track trace --from 24h
datadog error search "service:panamera-rest" --track trace --from 7d
datadog error search --track logs --from 1h
datadog error search --track trace --from 24h --json
```

## get

Get detailed information about a specific error issue.

```bash
datadog error get <issue-id>
```

Shows: ID, Error Type, Message, State, File, Function, Languages, First/Last Seen, Versions.

```bash
datadog error get 3a4fc6ca-dcb2-11ef-8282-da7ad0900002
datadog error get 3a4fc6ca-dcb2-11ef-8282-da7ad0900002 --json
```

**Required App Key scope:** `error_tracking_read`
