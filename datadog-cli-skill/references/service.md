# datadog service

Query Datadog APM services.

## list

List all active services reporting to APM for a given environment.

```bash
datadog service list [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--env` | `prod` | Environment to filter |

**Output:** Sorted list of service names with total count.

```bash
datadog service list --env prod
datadog service list --env staging
datadog service list --env prod --json
```

**Required App Key scope:** `apm_read`

**Use case:** Cross-reference with monitors to find stale monitors pointing to services that no longer exist.
