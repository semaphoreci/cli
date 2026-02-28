# Design: Surface after-pipeline jobs in pipeline describe

## Problem

When a pipeline has after-pipeline jobs (cleanup/publish steps that run after the main pipeline), `sem get pipeline` doesn't show them. The REST API already returns `with_after_task` and `after_task_id` fields, but the CLI's `PipelineV1Alpha` model doesn't declare them, so `json.Unmarshal` silently drops them.

This matches the issue fixed for the MCP server in [semaphoreio/semaphore#866](https://github.com/semaphoreio/semaphore/pull/866).

## Approach (Phase 1: flag only)

Add `WithAfterTask` and `AfterTaskID` fields to the `PipelineV1Alpha.Pipeline` struct with `omitempty` tags. No changes to the API client or describe command — the existing `DescribePpl` call already returns these fields.

### Changes

- `api/models/pipeline_v1_alpha.go`: Add two fields to the `Pipeline` inner struct.

### Output

Pipelines with after-tasks will show:

```yaml
pipeline:
  ppl_id: abc-123
  name: Deploy
  state: done
  result: passed
  with_after_task: true
  after_task_id: zebra-task-456
blocks:
  - name: Build
    ...
```

Pipelines without after-tasks are unchanged (omitempty).

## Future work (Phase 2: resolve job IDs)

The v1alpha `describe_topology` REST endpoint currently drops the `after_pipeline` field from its response. Once the server-side response formatter is updated to include it, the CLI can:

1. Add a `DescribeTopology` method to the pipelines API client
2. Call it when `with_after_task=true` to resolve actual job IDs
3. Display after-pipeline job IDs so users can run `sem logs <job_id>`
