# datadog metric

Query Datadog metric data via the V2 Scalar API.

## query

Query a metric and return the aggregated scalar value over a time range.

```bash
datadog metric query <metric-query> [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--from` | `1h` | Time range (e.g. `15m`, `1h`, `4h`, `1d`) |

The query uses standard Datadog metric query syntax: `aggregation:metric.name{tag:value}`.

**Output:** Query, Range, Name, Value.

```bash
datadog metric query "avg:system.cpu.user{env:prod}" --from 1h
datadog metric query "avg:holly.balance_state.delay{env:prod}" --from 4h
datadog metric query "sum:trace.rack.request.hits{service:panamera-rest}.as_rate()" --from 30m
datadog metric query "avg:kubernetes.containers.restarts{kube_cluster_name:prod}" --from 1h --json
```

**Required App Key scope:** `timeseries_query`

**Use case:** Validate monitor thresholds against actual metric values. If a monitor has `critical > 50` but the metric's average is `3.9`, the threshold is reasonable.
