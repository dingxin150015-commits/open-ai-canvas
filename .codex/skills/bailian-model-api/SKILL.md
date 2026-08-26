---
name: bailian-model-api
description: Use for Alibaba Cloud Bailian or Qianwen Model Studio API integration, debugging, capability mapping, and request design across Qwen text/image, Wan/HappyHorse video, speech/audio, embeddings, reranking, and platform APIs. Ground answers in the indexed local official documentation and preserve exact model-specific sync/async contracts. Do not activate for generic AI prompts that do not involve Bailian APIs or model capabilities.
---

# Bailian Model API

Use the local official documentation mirror as the primary source. Do not answer protocol questions from historical project reports or memory when an exact model/API document is available.

## Source workflow

1. Identify the exact task, model ID, region, protocol, SDK/language, and whether the user asks for analysis, code, debugging, or a real call.
2. Search `references/source-index.md` or run:

   ```powershell
   python scripts/search_official_docs.py <model-or-topic>
   ```

3. Read the exact API page and, for asynchronous APIs, its matching query-result page. Read a developer guide only for workflow guidance; API reference fields win when they conflict.
4. Preserve model-specific limits. Never infer one Qwen/Wan/HappyHorse/CosyVoice model's parameters from another model or from a generic protocol name.
5. Cite the local relative source path and state whether the claim is exact, inferred, or not present in the indexed mirror.

The complete 473-document manifest, hashes, and titles are in `references/source-index.md`; machine-readable metadata is in `references/source-catalog.json`.

## Route by task

- Text, Responses, Chat Completions, thinking, tool calling, files, batch: read `references/text-and-tools.md`.
- Qwen Image, Wan image generation/editing, image utilities: read `references/image-generation.md`.
- Wan, HappyHorse, Vidu, PixVerse, Kling, video editing: read `references/video-generation.md`.
- ASR, realtime audio, TTS, voice design/cloning, translation: read `references/speech-and-audio.md`.
- Embeddings, reranking, platform files/models/fine-tuning or billing: read `references/embeddings-and-platform.md`.
- Authentication, endpoints, async polling, secrets, output persistence, error handling: read `references/protocol-and-safety.md`.
- Updating the mirror or skill: read `references/maintenance.md`.

## Non-negotiable boundaries

- Internal product capability, UI values, and upstream payload format are separate layers.
- Use `x` for normalized internal pixel sizes when the application contract requires it; convert to the exact provider format only at the adapter boundary.
- `auto`, aspect ratio, and pixel size are distinct. For Qwen Image 3.0, omitting `size` lets the model recommend a resolution; do not serialize the literal string `auto` unless the exact API page says so.
- Do not silently clamp, coerce, drop, or replace unsupported durations, resolutions, media, output counts, tools, or model IDs. Reject or surface a precise validation error.
- Store API keys in approved secret storage or environment variables. Never place them in URLs, committed files, logs, task metadata, or examples with real values.
- Real API calls may consume quota or money. Before any call, state provider, account/region, model, endpoint, parameters, expected charge/quota, retries, and authorization; wait for explicit approval.
- Generated media URLs are often temporary. Materialize authorized results promptly into controlled storage and retain request/task IDs for audit.

## Deliverable standard

For an integration or diagnosis, report:

- exact model and API mode;
- endpoint and authentication shape;
- request/response schema and model-specific limits;
- sync/async lifecycle, polling and terminal states;
- output persistence and URL expiry;
- error and retry behavior;
- confirmed facts versus unresolved runtime/account questions;
- source document paths used.

Do not claim an integration works from documentation or compilation alone. Runtime verification requires an authorized call or a faithful mock/fixture test.
