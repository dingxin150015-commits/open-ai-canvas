# Video generation

Video APIs in the mirror are predominantly asynchronous but not uniform. Always read both the create-task and query-result pages for the exact model family and operation.

## Wan 2.7

Primary source folders:

- `api-reference/video-generation/wan27-text-to-video/`
- `api-reference/video-generation/wan27-image-to-video/`
- `api-reference/video-generation/wan27-reference-to-video/`
- `api-reference/video-generation/wan27-video-editing/`

Confirmed distinctions:

- Text-to-video and reference-to-video use `resolution` (`720P` or `1080P`) plus `ratio`, not arbitrary pixel `size`.
- Reference-to-video supports ratios documented per page and may infer ratio from `first_frame`.
- Snapshot model IDs and rolling aliases coexist. Preserve the exact requested ID; do not silently upgrade.
- Negative prompt, duration, audio, watermark and media constraints are operation/model specific.

## Wan 3.0 and other Wan operations

- Read `api-reference/video-generation/wan30-video/` for Wan 3.0.
- First-frame, first/last-frame, reference-to-video, speech-to-video, animation, editing and character replacement have separate folders and constraints.

## HappyHorse

Primary source folders:

- `happyhorse-text-to-video/`
- `happyhorse-image-to-video/`
- `happyhorse-reference-to-video/`
- `happyhorse-video-editing/`

Do not inherit Wan media counts or operation semantics. HappyHorse image/video reference support differs by operation.

## Other providers

Vidu, PixVerse, Kling and specialized video services use different model IDs, endpoints and constraints. Route through their exact source folders.

## Adapter rules

- Capability profiles must be per concrete model and operation.
- Reject unsupported duration, resolution, ratio, media count or audio options before task creation.
- Never silently convert 480P to 720P, 16 seconds to 15 seconds, or discard extra media.
- Preserve `task_id`, poll terminal status, surface provider code/message, and persist authorized output before signed URL expiry.
