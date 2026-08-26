# Image generation

## Qwen Image 3.0

Use the exact Qwen pages, not legacy Wanx assumptions.

Primary sources:

- `api-reference/image-generation/qwen-text-to-image.md`
- `api-reference/image-generation/qwen-text-to-image-30-async.md`
- `api-reference/image-generation/qwen-text-to-image-task-query.md`
- `api-reference/image-generation/qwen-image-editing.md`
- `developer-guides/getting-started/image-models.md`

Confirmed contract from the indexed mirror:

- Sync endpoint: `POST /services/aigc/multimodal-generation/generation` under the native `/api/v1` base.
- Request uses `model`, `input.messages`, and `parameters`.
- Text-to-image input is single-turn: one `user` message whose `content` contains exactly one text object.
- `qwen-image-3.0-pro` and `qwen-image-3.0` support text-to-image and image editing.
- Qwen Image 3.0 pixel dimensions use upstream `宽*高`; total pixel range is from `512*512` through `2048*2048`, aspect ratio 1:8 through 8:1.
- Omitting `size` lets the 3.0 model recommend a resolution. Treat UI `auto` as omission at serialization time.
- Qwen Image 3.0/2.0 output count `n` is 1-6; max/plus/image families have different limits.
- `negative_prompt` is at most 500 characters.
- 3.0 supports `prompt_extend_mode` values `direct` and `agent`; `agent` applies only where the page states it is supported.
- `enable_thinking` and `prompt_extend` interact; do not expose the combination without preserving the documented dependency.
- Sync 3.0 results appear in `output.choices[].message.content[].image`; output URLs are PNG and documented as valid for 24 hours.

## Async Qwen Image

- Qwen Image 3.0 has its own async page and still uses the message-shaped input.
- Legacy `qwen-image-plus`/`qwen-image` async requests use a different prompt-shaped input and fixed dimensions/count rules.
- Never merge these schemas into one generic payload without a model-aware adapter.

## Editing and references

- Qwen Image 3.0 editing supports model-specific input image counts, dimension limits and output counts. Read `qwen-image-editing.md` before accepting references.
- Validate reference accessibility and media constraints before creating a paid task.
- Keep mask, multi-image fusion and output count separate; support in one family does not imply support in another.

## Wan image and utilities

The mirror also includes Wan 2.1/2.5/2.6/2.7 image generation/editing, background generation, outpainting, erasure, virtual try-on, poster, segmentation, word art and third-party image APIs. Each create/query pair has separate payloads and limits.

Search exact folders under `api-reference/image-generation/`; do not reuse the Qwen message schema for Wan or utility services.
