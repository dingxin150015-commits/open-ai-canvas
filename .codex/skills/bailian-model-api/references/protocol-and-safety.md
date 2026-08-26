# Protocol and safety

## Source hierarchy

1. Exact `api-reference/` page for the requested model/operation.
2. Matching create-task/query-result pair for asynchronous APIs.
3. `developer-guides/` for workflow and best practices.
4. `resources/` and `token-plan/` for billing/account product context.
5. Historical project reports only as non-authoritative context.

## Authentication and regions

- Most DashScope HTTP APIs use `Authorization: Bearer $DASHSCOPE_API_KEY`.
- Anthropic-compatible pages may use the Anthropic header contract; do not force Bearer authentication across every compatibility surface.
- Common Beijing endpoints in this mirror include:
  - Native API: `https://dashscope.aliyuncs.com/api/v1`
  - OpenAI-compatible: `https://dashscope.aliyuncs.com/compatible-mode/v1`
  - WebSocket surfaces under `/api-ws/v1/...`
- Token Plan uses separate endpoints and entitlements. Never silently substitute a Token Plan endpoint for ordinary Bailian API access.
- Region, account entitlement, model availability and quota are runtime facts. Verify them before a real call.

Primary sources:

- `api-reference/preparation/api-key.md`
- `api-reference/preparation/export-api-key-env.md`
- `api-reference/chat/openai-chat.md`
- `api-reference/chat/openai-responses.md`
- `api-reference/chat/dashscope.md`
- `api-reference/chat/anthropic.md`

## Sync and async

- Do not classify an entire modality as sync or async. The exact operation/model page decides.
- Async create requests typically require `X-DashScope-Async: enable`, return `output.task_id`, and are queried with `GET /tasks/{task_id}`.
- Model/task pages define valid states. Common states include `PENDING`, `RUNNING`, `SUCCEEDED`, `FAILED`; some services add other states.
- Poll with a bounded interval, deadline, cancellation propagation, and terminal error mapping. Do not retry create blindly after an ambiguous response; preserve idempotency or reconcile by task/request ID where supported.
- Task/result retention windows differ by service. Read the matching query page.

Primary sources:

- `developer-guides/run-and-scale/async-task-management.md`
- `api-reference/image-generation/qwen-text-to-image-30-async.md`
- `api-reference/image-generation/qwen-text-to-image-task-query.md`
- create/query pairs under `api-reference/video-generation/`

## Temporary results

- Many image/video result URLs are signed OSS URLs and commonly expire after 24 hours, but the exact page controls.
- Download only after an authorized generation. Validate HTTP status, media type, non-empty content and size before persisting.
- Keep `request_id` and `task_id` without logging secrets or full sensitive inputs.

## Errors and retries

- Preserve provider `request_id`, error code and safe message.
- Authentication, invalid parameters, content safety, quota, throttling, unsupported model and inaccessible media are different failure classes.
- Retry only transient, documented failures with bounded attempts and backoff. Never retry validation, authentication, billing or safety errors as if transient.
- See `api-reference/preparation/error-messages.md` and the exact operation response schema.
