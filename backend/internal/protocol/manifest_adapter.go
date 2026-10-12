// Streamed JSON audio parsing added by upstream v1.6.0.
package protocol

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func (a manifestAdapter) ParseCreateWithRequest(_ context.Context, request GenerationRequest, body []byte) (CreateResult, error) {
	if a.manifest.Response.StreamedJSONAudio {
		format := ""
		if request.Extra != nil {
			format = manifestString(request.Extra["audioFormat"])
		}
		return streamedJSONAudioCreateResult(body, a.manifest.Response.ResultKind, format)
	}
	if a.manifest.Response.BinaryPayload {
		return binaryPayloadCreateResult(a.manifest.Response.ResultKind, body)
	}
	payload, err := decodeObject(body)
	if err != nil {
		return CreateResult{}, err
	}
	if err := a.validateResponse(payload, PollContext{}); err != nil {
		return CreateResult{}, err
	}
	return a.parse(payload, PollContext{}), nil
}

// streamedJSONAudioCreateResult decodes the chunked JSON response used by
// Doubao's unidirectional TTS endpoint. Go exposes HTTP chunked transfer as a
// single byte stream, so the decoder accepts both concatenated JSON objects and
// newline/SSE-style data frames.
const maxStreamedJSONAudioBytes = 64 << 20

func streamedJSONAudioCreateResult(body []byte, resultKind string, format string) (CreateResult, error) {
	if len(body) == 0 {
		return CreateResult{}, fmt.Errorf("streamed JSON audio response is empty")
	}
	if len(body) > maxStreamedJSONAudioBytes {
		return CreateResult{}, fmt.Errorf("streamed JSON audio response exceeds %d bytes", maxStreamedJSONAudioBytes)
	}
	frames, err := decodeJSONFrames(body)
	if err != nil {
		return CreateResult{}, fmt.Errorf("decode streamed JSON audio response: %w", err)
	}
	var audio bytes.Buffer
	var usage map[string]any
	for _, frame := range frames {
		code := strings.ToLower(strings.TrimSpace(manifestString(firstNonNilPathValue(frame, "code", "status_code", "statusCode"))))
		if code != "" && code != "0" && code != "ok" && code != "success" && code != "succeeded" {
			message := strings.TrimSpace(manifestString(firstNonNilPathValue(frame, "message", "msg", "error.message")))
			if message == "" {
				message = "上游返回失败状态"
			}
			return CreateResult{}, fmt.Errorf("流式音频合成失败（code %s）：%s", code, message)
		}
		if value := firstPathValue(frame, "data", "audio", "audio.data"); value != "" {
			decoded, decodeErr := decodeStreamedAudioBase64(value)
			if decodeErr != nil {
				return CreateResult{}, fmt.Errorf("decode streamed audio data: %w", decodeErr)
			}
			_, _ = audio.Write(decoded)
		}
		if rawUsage := pathValue(frame, "usage"); rawUsage != nil {
			usage = manifestObject(rawUsage)
		}
	}
	if audio.Len() == 0 {
		return CreateResult{}, fmt.Errorf("streamed JSON audio response contains no audio data")
	}
	mimeType := streamedAudioMIMEType(format)
	reference := MediaReference{DataURL: "data:" + mimeType + ";base64," + base64.StdEncoding.EncodeToString(audio.Bytes()), MIMEType: mimeType}
	return CreateResult{Status: StatusSucceeded, Result: &Result{Audios: []MediaReference{reference}, Usage: usage}}, nil
}

func decodeJSONFrames(body []byte) ([]map[string]any, error) {
	decoder := json.NewDecoder(bytes.NewReader(body))
	frames := make([]map[string]any, 0, 4)
	for {
		var frame map[string]any
		err := decoder.Decode(&frame)
		if err == io.EOF {
			return frames, nil
		}
		if err != nil {
			break
		}
		if frame != nil {
			frames = append(frames, frame)
		}
	}
	// Some gateways wrap each JSON object in an SSE-like "data:" line.
	for _, line := range strings.Split(string(body), "\n") {
		line = strings.TrimSpace(strings.TrimSuffix(line, "\r"))
		line = strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if line == "" || line == "[DONE]" {
			continue
		}
		var frame map[string]any
		if err := json.Unmarshal([]byte(line), &frame); err != nil {
			return nil, err
		}
		frames = append(frames, frame)
	}
	if len(frames) == 0 {
		return nil, fmt.Errorf("no JSON frames found")
	}
	return frames, nil
}

func streamedAudioMIMEType(format string) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "wav":
		return "audio/wav"
	case "ogg_opus", "ogg":
		return "audio/ogg"
	case "pcm":
		return "audio/pcm"
	default:
		return "audio/mpeg"
	}
}

func firstNonNilPathValue(payload map[string]any, paths ...string) any {
	for _, path := range paths {
		if value := pathValue(payload, path); value != nil {
			return value
		}
	}
	return nil
}

func decodeStreamedAudioBase64(value string) ([]byte, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err == nil {
		return decoded, nil
	}
	decoded, rawErr := base64.RawStdEncoding.DecodeString(value)
	if rawErr == nil {
		return decoded, nil
	}
	return nil, err
}
