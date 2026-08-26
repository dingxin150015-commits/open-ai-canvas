# Text, chat and tools

## API surfaces

Primary sources:

- `api-reference/chat/openai-chat.md`
- `api-reference/chat/openai-responses.md`
- `api-reference/chat/dashscope.md`
- `api-reference/chat/anthropic.md`
- `developer-guides/text-generation/`
- `developer-guides/tool-calling/`

Choose one API surface explicitly:

- OpenAI-compatible Chat Completions;
- OpenAI-compatible Responses;
- native DashScope Generation/Multimodal APIs;
- Anthropic-compatible Messages where documented.

Do not mix headers, base URLs, streaming frames or tool-call response shapes across surfaces.

## Model-specific behavior

- Thinking/reasoning defaults, `enable_thinking`, reasoning content preservation and tool-choice support differ across Qwen and third-party models.
- Preserve reasoning fields across tool-call turns when the exact guide requires it.
- Structured output, JSON schema, files, document understanding, long context, translation, deep research and batch each have dedicated pages.
- Model aliases and dated snapshots change. Use the exact model requested and consult `changelog/models.md` plus the exact API page.

## Streaming and context

- Read `developer-guides/run-and-scale/streaming.md`, `multi-turn.md`, `context-cache.md`, `token-counting.md` and the selected API reference.
- Pass cancellation through the whole stream and distinguish transport completion from model finish reason.
- Do not infer context length, cache eligibility or token price from another model.

## Tools

- Function calling, MCP, web search/scraping, image search, code interpreter and PDF understanding have separate guides.
- Validate tool schemas, loop limits and tool-result correlation. Never expose project secrets or grant external tools authority beyond the user request.
