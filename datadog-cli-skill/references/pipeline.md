# datadog pipeline

Manage Datadog log processing pipelines.

## list

List all log pipelines.

```bash
datadog pipeline list
```

**Output columns:** ID, NAME, ENABLED, FILTER, PROCESSORS (count)

## get

Get pipeline details including all processors.

```bash
datadog pipeline get <pipeline-id>
```

Shows pipeline metadata and a list of processors with type-specific details (Category Processor rules, Remapper mappings, Grok patterns, etc.).

```bash
datadog pipeline get gX3YuAs5RkmKJQOpVKhBsA
datadog pipeline get gX3YuAs5RkmKJQOpVKhBsA --json
```

## create

Create a new log pipeline.

```bash
datadog pipeline create [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--name` | (required) | Pipeline name |
| `--filter` | | Log filter query (e.g. `source:nginx`, `service:myapp`, `*`) |
| `--enabled` | `false` | Enable pipeline immediately |

```bash
datadog pipeline create --name "Team Tagging" --filter "*" --enabled
datadog pipeline create --name "UC Logs" --filter "service:(counter-rest OR panamera-rest)"
```

## delete

Delete a log pipeline.

```bash
datadog pipeline delete <pipeline-id>
```

## add-category

Add a Category Processor to an existing pipeline. Categorizes logs into groups based on query rules.

```bash
datadog pipeline add-category <pipeline-id> [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--target` | (required) | Target attribute name |
| `--rule` | (required, repeatable) | Rule in format `query=value` |
| `--processor-name` | `Category Processor` | Processor display name |

Rules use Datadog log query syntax for the query part, `=` separates query from category value.

```bash
datadog pipeline add-category abc123 --target team \
  --rule "service:(counter-rest OR nagoya OR panamera-rest)=plat" \
  --rule "service:(smaug OR engine OR strategy-rpc)=spot" \
  --rule "service:(contract-api OR contract-ws)=contract"
```

## add-remapper

Add an Attribute Remapper to an existing pipeline. Maps attributes between types (attribute <-> tag).

```bash
datadog pipeline add-remapper <pipeline-id> [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `--source` | (required) | Source attribute name |
| `--target` | (required) | Target attribute name |
| `--source-type` | `attribute` | Source type: `attribute` or `tag` |
| `--target-type` | `tag` | Target type: `attribute` or `tag` |
| `--preserve-source` | `false` | Keep source after remapping |
| `--override` | `false` | Override target if exists |
| `--processor-name` | `Attribute Remapper` | Processor display name |

```bash
# Convert a log attribute to a tag
datadog pipeline add-remapper abc123 \
  --source team --target team \
  --source-type attribute --target-type tag
```

## Full Workflow: Team Tag Pipeline

```bash
# 1. Create pipeline
datadog pipeline create --name "Team Tagging" --filter "*" --enabled
# Note the returned pipeline ID

# 2. Add category processor to map services to teams
datadog pipeline add-category <id> --target team \
  --rule "service:(counter-rest OR nagoya)=plat" \
  --rule "service:(smaug OR engine)=spot" \
  --rule "service:(contract-api)=contract"

# 3. Add remapper to convert attribute to tag
datadog pipeline add-remapper <id> \
  --source team --target team \
  --source-type attribute --target-type tag

# 4. Verify
datadog pipeline get <id>
```
