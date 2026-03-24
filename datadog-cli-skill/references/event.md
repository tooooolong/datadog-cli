# datadog event

Query Datadog events and alert history.

## list

List events with filters and time range.

```bash
datadog event list [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--query` | `source:monitor` | Event search query |
| `--from` | `24h` | Time range (e.g. `1h`, `24h`, `7d`) |
| `--limit` | `20` | Maximum events to return |

**Query syntax** follows Datadog event search:

| Query | Description |
|-------|-------------|
| `source:monitor` | All monitor events |
| `source:monitor monitor_id:12345` | Events for a specific monitor |
| `source:monitor status:alert` | Alert trigger events |
| `*` | All events |

**Output columns:** TIME, MESSAGE, TAGS

```bash
datadog event list --query "source:monitor" --from 24h
datadog event list --query "source:monitor monitor_id:28140542" --from 7d
datadog event list --query "*" --from 1h --limit 5
datadog event list --query "source:kubernetes" --from 4h --json
```

**Required App Key scope:** `events_read`
