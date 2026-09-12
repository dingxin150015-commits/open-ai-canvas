package service

import (
	"context"
	"encoding/json"
	"infinite-canvas/backend/internal/model"
	"infinite-canvas/backend/internal/protocol"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBailianOfficialLongAudioAndContinuationDimensions(t *testing.T) {
	adapter := migratedBailianAdapter(t, "dashscope-video")
	input := canvasGenerationInput{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan2.7-t2v", VideoSeconds: "5"}, ReferenceAudios: []providerMedia{{URL: "https://example.com/a.wav", MimeType: "audio/wav", DurationMs: 6000}}}
	if _, err := adapter.BuildCreate(context.Background(), protocol.RequestContext{Request: protocolRequestFromInput(input)}); err != nil {
		t.Fatal(err)
	}
	input.Config.Model = "wan2.7-i2v"
	input.Prompt = ""
	input.ReferenceAudios = nil
	input.ReferenceVideos = []providerMedia{{URL: "https://example.com/a.mp4", MimeType: "video/mp4", DurationMs: 3000, Width: 720, Height: 1280}}
	if _, err := adapter.BuildCreate(context.Background(), protocol.RequestContext{Request: protocolRequestFromInput(input)}); err != nil {
		t.Fatal(err)
	}
	input.ReferenceVideos[0].Width = 4097
	if _, err := adapter.BuildCreate(context.Background(), protocol.RequestContext{Request: protocolRequestFromInput(input)}); err == nil {
		t.Fatal("oversize continuation accepted")
	}
}

func TestBailianQwenTransportContracts(t *testing.T) {
	adapter := migratedBailianAdapter(t, "dashscope-image")
	for _, mode := range []string{"sync", "async", "openai"} {
		t.Run(mode, func(t *testing.T) {
			input := canvasGenerationInput{Mode: "image", Prompt: "fixture", Config: providerConfig{Model: "qwen-image-3.0-pro", Count: "6", Size: "1024x1024"}, Metadata: map[string]any{"request_mode": mode, "seed": "0"}}
			spec, err := adapter.BuildCreate(context.Background(), protocol.RequestContext{Request: protocolRequestFromInput(input)})
			if err != nil {
				t.Fatal(err)
			}
			raw, _ := json.Marshal(spec.Body)
			if !strings.Contains(string(raw), `"n":6`) {
				t.Fatalf("count lost: %s", raw)
			}
			if mode == "openai" && (!strings.Contains(spec.Path, "compatible-mode") || !strings.Contains(string(raw), `"size":"1024x1024"`)) {
				t.Fatalf("bad compatible request: %#v", spec)
			}
			if mode == "async" && !strings.Contains(spec.Path, "/image-generation/") {
				t.Fatalf("bad async path: %s", spec.Path)
			}
		})
	}
	created, err := adapter.ParseCreate(context.Background(), []byte(`{"output":{"task_id":"t1","task_status":"PENDING"}}`))
	if err != nil || created.Status != protocol.StatusPending || created.TaskID != "t1" {
		t.Fatalf("async create: %#v %v", created, err)
	}
	result, err := adapter.ParsePoll(context.Background(), protocol.PollContext{TaskID: "t1"}, []byte(`{"output":{"task_status":"SUCCEEDED","results":[{"url":"https://example.com/a.png"}]}}`))
	if err != nil || result.Result == nil || len(result.Result.Images) != 1 {
		t.Fatalf("async result: %#v %v", result, err)
	}
}

func TestBailianDurationMetering(t *testing.T) {
	log := &model.ApiCallLog{Model: "wan2.7-r2v", Capability: "video", Status: model.ApiCallStatusSucceeded}
	(&Service{}).enrichAPICallLogPayload(log, map[string]any{"output": map[string]any{"task_status": "SUCCEEDED"}, "usage": map[string]any{"duration": 8.25, "input_video_duration": 3.25, "output_video_duration": 5}})
	if log.UsageDurationMs != 8250 || log.InputVideoDurationMs != 3250 || log.OutputVideoDurationMs != 5000 || !log.UsageAvailable {
		t.Fatalf("bad usage: %#v", log)
	}
}

func TestBailianEmptyPromptAdmissionRequiresMedia(t *testing.T) {
	input := map[string]any{"mode": "video", "config": map[string]any{"model": "wan2.7-i2v"}}
	if bailianTaskAllowsEmptyPrompt(input) {
		t.Fatal("media-less empty prompt accepted")
	}
	input["referenceImages"] = []any{map[string]any{"id": "image"}}
	if !bailianTaskAllowsEmptyPrompt(input) {
		t.Fatal("image-only prompt rejected")
	}
	input["config"] = map[string]any{"model": "happyhorse-1.1-r2v"}
	if bailianTaskAllowsEmptyPrompt(input) {
		t.Fatal("R2V empty prompt accepted")
	}
}

func TestBailianOpenAIHeaderRequestID(t *testing.T) {
	if id := providerRequestIDFromHeader(http.Header{"X-Request-Id": []string{"req-qwen-1"}}); id != "req-qwen-1" {
		t.Fatal(id)
	}
	for _, value := range []string{"Bearer secret", "https://example.com/?key=secret", strings.Repeat("a", 161)} {
		if providerRequestIDFromHeader(http.Header{"X-Request-Id": []string{value}}) != "" {
			t.Fatal("unsafe request ID accepted")
		}
	}
}

func TestBailianMeteringRequiresTerminalReceiptAndSupportsSKUKey(t *testing.T) {
	for _, status := range []string{"PENDING", "RUNNING", "FAILED", "SUCCEEDED"} {
		log := &model.ApiCallLog{Model: "private-sku", Path: "/api/v1/tasks/fixture", Capability: "video"}
		(&Service{}).enrichAPICallLogPayload(log, map[string]any{"output": map[string]any{"task_status": status}, "usage": map[string]any{"duration": 8.25}})
		if (log.UsageDurationMs > 0) != (status == "SUCCEEDED") {
			t.Fatalf("status %s got %d", status, log.UsageDurationMs)
		}
	}
}

func TestBailianQwenAllTransportsMockLifecycle(t *testing.T) {
	for _, mode := range []string{"sync", "async", "openai"} {
		t.Run(mode, func(t *testing.T) {
			creates, polls, downloads := 0, 0, 0
			var server *httptest.Server
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.URL.Path == "/result.png" {
					downloads++
					if r.Header.Get("Authorization") != "" {
						t.Error("result download leaked authentication")
					}
					w.Header().Set("Content-Type", "image/png")
					data, _, err := mediaBytes(providerMedia{DataURL: qwenImageTestPNG})
					if err != nil {
						t.Error(err)
						return
					}
					w.Write(data)
					return
				}
				if r.Method == "POST" {
					creates++
					if mode == "async" {
						if r.Header.Get("X-DashScope-Async") != "enable" {
							t.Error("async header missing")
						}
						json.NewEncoder(w).Encode(map[string]any{"output": map[string]any{"task_id": "qwen-task", "task_status": "PENDING"}})
						return
					}
				} else {
					polls++
				}
				url := server.URL + "/result.png"
				if mode == "openai" {
					json.NewEncoder(w).Encode(map[string]any{"data": []any{map[string]any{"url": url}}})
					return
				}
				if mode == "async" {
					json.NewEncoder(w).Encode(map[string]any{"output": map[string]any{"task_status": "SUCCEEDED", "results": []any{map[string]any{"url": url}}}})
					return
				}
				json.NewEncoder(w).Encode(map[string]any{"output": map[string]any{"choices": []any{map[string]any{"message": map[string]any{"content": []any{map[string]any{"image": url}}}}}}})
			}))
			defer server.Close()
			config := providerConfig{InterfaceType: "dashscope-image", Model: "qwen-image-3.0-pro", BaseURL: server.URL, APIKey: "fixture", AllowLocalChannel: true, Count: "1", Size: "auto"}
			ctx := withProviderOutboundPolicy(context.Background(), config)
			result, err := runProtocolAdapterTaskWithTiming(ctx, canvasGenerationInput{Mode: "image", Prompt: "fixture", Config: config, Metadata: map[string]any{"request_mode": mode}}, migratedBailianAdapter(t, "dashscope-image"), protocolPollTiming{})
			if err != nil || result["mode"] != "image" || creates != 1 || downloads != 1 {
				t.Fatalf("creates=%d polls=%d downloads=%d error=%v", creates, polls, downloads, err)
			}
			if mode == "async" && polls != 1 {
				t.Fatalf("polls=%d", polls)
			}
		})
	}
}
