# datadog log

Search Datadog log data.

## search

Search logs using Datadog log query syntax.

```bash
datadog log search <query> [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--from` | `15m` | Time range (e.g. `5m`, `1h`, `24h`) |
| `--limit` | `10` | Maximum logs to return |

**Query syntax** follows Datadog log search:

| Query | Description |
|-------|-------------|
| `service:myapp` | Logs from a specific service |
| `service:myapp status:error` | Error logs from a service |
| `@http.status_code:5*` | 5xx responses |
| `host:my-host` | Logs from a specific host |

**Output columns:** TIME, STATUS, SERVICE, MESSAGE

```bash
datadog log search "service:panamera-rest" --from 1h
datadog log search "service:lego-v2-rest status:error" --from 24h --limit 20
datadog log search "status:error" --from 1h --limit 5
datadog log search "service:counter-rest" --from 5m --json
```

**Required App Key scope:** `logs_read_data`
