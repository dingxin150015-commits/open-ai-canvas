# Speech and audio

Primary source groups:

- `api-reference/speech-recognition/`
- `api-reference/speech-synthesis/`
- `api-reference/speech-translation/`
- `api-reference/qwen-audio-realtime/`
- `api-reference/real-time-multimodal/`
- `developer-guides/speech/`
- `developer-guides/realtime-api/`

## Routing

- ASR file, realtime ASR, streaming ASR and file translation are different lifecycles.
- TTS includes Qwen Audio TTS, Qwen TTS, CosyVoice, Sambert, realtime TTS, voice design and voice cloning. Do not share payloads by modality alone.
- Realtime APIs may use WebSocket event protocols with mandatory lifecycle events and fixed empty objects. Preserve exact event ordering and required fields.

## Safety and consent

- Voice cloning and persistent voice assets require user authorization and lawful source material.
- Do not upload or retain voice samples beyond the authorized task.
- Validate sample format, duration, rate, channel count, language and size before a paid call.
- Avoid time-stretching generated speech unless the requested workflow explicitly allows it; prefer model pacing and regeneration.

## Results

- Distinguish raw audio bytes, streamed chunks, URLs, task results and synthesized metadata.
- Persist authorized outputs before temporary URLs expire and record safe request/task IDs.
