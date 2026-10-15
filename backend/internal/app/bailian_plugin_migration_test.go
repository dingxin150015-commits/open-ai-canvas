package app

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"yingce/backend/internal/model"
	"yingce/backend/internal/protocol"
)

func migratedBailianAdapter(t *testing.T, id string) protocol.Adapter {
	t.Helper()
	adapter, ok := loadOfficialFallbackRegistry().Get(id)
	if !ok || adapter.Metadata().Execution != "declarative" {
		t.Fatalf("packaged declarative provider %s unavailable", id)
	}
	return adapter
}

func canonicalFixture(t *testing.T, value any) any {
	t.Helper()
	raw, e := json.Marshal(value)
	if e != nil {
		t.Fatal(e)
	}
	var result any
	if e = json.Unmarshal(raw, &result); e != nil {
		t.Fatal(e)
	}
	return result
}

func assertPluginMatchesOracle(t *testing.T, input canvasGenerationInput) {
	t.Helper()
	id := "dashscope-video"
	var expected map[string]any
	var err error
	if input.Mode == "image" {
		id = "dashscope-image"
		expected, err = buildQwenImage30Request(input)
	} else {
		expected, err = buildDashScopeVideoRequest(input)
	}
	if err != nil {
		t.Fatalf("invalid oracle fixture: %v", err)
	}
	spec, err := migratedBailianAdapter(t, id).BuildCreate(context.Background(), protocol.RequestContext{Request: protocolRequestFromInput(input)})
	if err != nil {
		t.Fatalf("plugin rejected accepted oracle fixture: %v", err)
	}
	if !reflect.DeepEqual(canonicalFixture(t, expected), canonicalFixture(t, spec.Body)) {
		t.Fatalf("wire contract differs\nwant=%s\ngot =%s", mustFixtureJSON(expected), mustFixtureJSON(spec.Body))
	}
}
func mustFixtureJSON(value any) []byte { data, _ := json.Marshal(value); return data }

func TestBailianPluginGoldenRequests(t *testing.T) {
	image := canvasGenerationInput{Mode: "image", Prompt: "fixture", Config: providerConfig{Model: "qwen-image-3.0-pro", Size: "auto", Count: "6"}, Metadata: map[string]any{"seed": 0, "prompt_extend": false, "enable_thinking": false}}
	assertPluginMatchesOracle(t, image)
	image.Config.SystemPrompt = "  consistent style  "
	assertPluginMatchesOracle(t, image)
	image.Config.SystemPrompt = ""
	image.Config.Size = "5:4"
	image.Config.Count = "2"
	image.Metadata = map[string]any{"prompt_extend_mode": "direct"}
	image.ReferenceImages = []providerMedia{{DataURL: qwenImageTestPNG}, {DataURL: qwenImageTestPNG}, {DataURL: qwenImageTestPNG}}
	assertPluginMatchesOracle(t, image)
	for _, key := range []string{"wan2.7-t2v", "wan2.7-t2v-2026-04-25", "wan2.7-t2v-2026-06-12", "wan2.7-i2v", "wan2.7-i2v-2026-04-25", "wan2.7-r2v", "wan2.7-r2v-2026-06-12", "happyhorse-1.1-t2v", "happyhorse-1.1-i2v", "happyhorse-1.1-r2v", "wan3.0-video"} {
		t.Run(key, func(t *testing.T) {
			input := canvasGenerationInput{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: key, VQuality: "720p", VideoSeconds: "5", Size: "16:9"}, Metadata: map[string]any{"seed": 17}}
			if strings.Contains(key, "i2v") || strings.Contains(key, "r2v") {
				input.ReferenceImages = []providerMedia{{ID: "first", URL: "https://example.com/first.png", MimeType: "image/png", Width: 720, Height: 720}}
			}
			assertPluginMatchesOracle(t, input)
			input.Config.VideoWatermark = "false"
			assertPluginMatchesOracle(t, input)
		})
	}
	frame := providerMedia{ID: "first", URL: "https://example.com/first.png", MimeType: "image/png", Width: 720, Height: 720}
	tail := frame
	tail.ID = "last"
	tail.URL = "https://example.com/last.png"
	clip := providerMedia{ID: "clip", URL: "https://example.com/video.mp4", MimeType: "video/mp4", DurationMs: 3000, Width: 720, Height: 720}
	audio := providerMedia{ID: "audio", URL: "https://example.com/audio.wav", MimeType: "audio/wav", DurationMs: 3000}
	cases := []canvasGenerationInput{
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan2.7-i2v", VideoSeconds: "5"}, ReferenceImages: []providerMedia{tail, frame}, ReferenceAudios: []providerMedia{audio}, Metadata: map[string]any{"videoStartFrameNodeId": "first", "videoEndFrameNodeId": "last"}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan2.7-i2v", VideoSeconds: "5"}, ReferenceImages: []providerMedia{tail}, ReferenceVideos: []providerMedia{clip}, Metadata: map[string]any{"videoEndFrameNodeId": "last"}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan2.7-r2v", VideoSeconds: "5"}, ReferenceImages: []providerMedia{frame}, ReferenceVideos: []providerMedia{clip}, ReferenceAudios: []providerMedia{audio}, Metadata: map[string]any{"videoStartFrameNodeId": "first"}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan3.0-video", VideoSeconds: "-1", VideoGenerateAudio: "false"}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan3.0-video", VideoSeconds: "5"}, ReferenceImages: []providerMedia{tail, frame}, Metadata: map[string]any{"videoStartFrameNodeId": "first", "videoEndFrameNodeId": "last", "videoEditOperation": "image_to_video"}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan3.0-video", VideoSeconds: "5"}, ReferenceImages: []providerMedia{frame}, ReferenceVideos: []providerMedia{clip}, ReferenceAudios: []providerMedia{audio}, Metadata: map[string]any{"videoEditOperation": "reference_to_video"}},
	}
	for _, input := range cases {
		assertPluginMatchesOracle(t, input)
	}
}

func TestBailianPluginRejectsUnsupportedRequests(t *testing.T) {
	cases := []canvasGenerationInput{
		{Mode: "image", Prompt: "fixture", Config: providerConfig{Model: "qwen-image-3.0-pro", Size: "4096x4096"}},
		{Mode: "image", Prompt: "fixture", Config: providerConfig{Model: "qwen-image-3.0-pro", Count: "7"}},
		{Mode: "image", Prompt: "fixture", Config: providerConfig{Model: "qwen-image-3.0-pro"}, ReferenceImages: []providerMedia{{DataURL: qwenImageTestPNG}}, Metadata: map[string]any{"prompt_extend_mode": "agent"}},
		{Mode: "image", Prompt: "fixture", Config: providerConfig{Model: "qwen-image-3.0-pro"}, Metadata: map[string]any{"prompt_extend": false, "enable_thinking": true}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan2.7-t2v", VideoSeconds: "16"}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan2.7-t2v", VQuality: "480p"}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "happyhorse-1.1-t2v", VideoSeconds: "2"}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "happyhorse-1.1-i2v"}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan3.0-video-prime"}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan3.0-video", VideoSeconds: "3.5"}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan2.7-t2v", VideoSeconds: "5"}, ReferenceAudios: []providerMedia{{URL: "https://example.com/a.wav", MimeType: "audio/wav", DurationMs: 31000}}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan2.7-r2v", VideoSeconds: "11"}, ReferenceVideos: []providerMedia{{URL: "https://example.com/a.mp4", MimeType: "video/mp4", DurationMs: 3000}}},
		{Mode: "video", Prompt: "fixture", Config: providerConfig{Model: "wan3.0-video", VideoSeconds: "25"}, ReferenceVideos: []providerMedia{{URL: "https://example.com/a.mp4", MimeType: "video/mp4", DurationMs: 8000}}},
	}
	for index, input := range cases {
		id := "dashscope-video"
		if input.Mode == "image" {
			id = "dashscope-image"
		}
		if _, err := migratedBailianAdapter(t, id).BuildCreate(context.Background(), protocol.RequestContext{Request: protocolRequestFromInput(input)}); err == nil {
			t.Fatalf("invalid fixture %d accepted", index)
		}
	}
}

func TestBailianPluginLifecycleSaveAndAdmission(t *testing.T) {
	dir := t.TempDir()
	center, err := newPluginRuntime(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc, db := newChannelModelTestService(t)
	svc.pluginRuntime = center
	admin := &model.User{ID: "admin", Role: model.UserRoleAdmin}
	channel := model.ModelChannel{ID: "bailian", UserID: admin.ID, Scope: model.ChannelScopeSystem, Enabled: true, Name: "fixture", BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1", APIKey: "fixture", APIFormat: "openai", ModelsJSON: `[]`}
	if err := db.Create(&channel).Error; err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"dashscope-image", "dashscope-video"} {
		key, capability, legacy := "qwen-image-3.0-pro", "image", "dashscope-wanx-image"
		if id == "dashscope-video" {
			key, capability, legacy = "wan2.7-t2v", "video", "dashscope-wan-video"
		}
		enabled := true
		req := ChannelModelRequest{ModelKey: key, Capability: capability, Protocol: id, Enabled: &enabled, CapabilityConfig: DefaultModelCapabilityConfigForModel(id, key), PriceTiers: []ChannelModelPriceTierRequest{{BillingMode: "fixed_request", UnitPriceMicrocredits: 1000000, PriceConfigured: true, Enabled: &enabled}}}
		saved, err := svc.SaveAdminChannelModel(admin, channel.ID, "", req)
		if err != nil {
			t.Fatal(err)
		}
		stored, err := svc.repo.ChannelModel(saved.ID)
		if err != nil || string(stored.Protocol) != id {
			t.Fatal("saved protocol not preserved")
		}
		if public, err := svc.sanitizeChannelModel(stored); err != nil || !public.Available {
			t.Fatal("enabled plugin model unavailable")
		}
		req.Protocol = legacy
		if _, err := svc.SaveAdminChannelModel(admin, channel.ID, saved.ID, req); err == nil {
			t.Fatal("old incompatible protocol accepted")
		}
		if _, err := svc.SetPluginEnabled(id, false); err != nil {
			t.Fatal(err)
		}
		if public, err := svc.sanitizeChannelModel(stored); err == nil && public.Available {
			t.Fatal("disabled plugin still advertised")
		}
		if _, err := svc.resolveSystemChannelModelSelection(map[string]any{"config": map[string]any{"channelId": channel.ID, "model": key}}, "canvas_"+capability, ""); err == nil {
			t.Fatal("disabled plugin admits task")
		}
		for _, item := range svc.PluginProviderCatalog("admin.system-channel", "", false) {
			if item.ID == id {
				t.Fatal("disabled plugin still selectable")
			}
		}
	}
	center, err = newPluginRuntime(dir)
	if err != nil {
		t.Fatal(err)
	}
	svc.pluginRuntime = center
	for _, id := range []string{"dashscope-image", "dashscope-video"} {
		if svc.protocolIsSelectable(id) {
			t.Fatal("restart lost disabled state")
		}
		if _, err := svc.SetPluginEnabled(id, true); err != nil {
			t.Fatal(err)
		}
		if !svc.protocolIsSelectable(id) {
			t.Fatal("reenabled plugin unavailable")
		}
	}
}

func TestBailianPluginMockSyncCreateAndDownload(t *testing.T) {
	var server *httptest.Server
	creates, downloads := 0, 0
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/services/aigc/multimodal-generation/generation":
			creates++
			if r.Header.Get("Authorization") != "Bearer fixture" || r.Header.Get("X-DashScope-Async") != "" {
				t.Error("wrong sync headers")
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"output": map[string]any{"choices": []any{map[string]any{"message": map[string]any{"content": []any{map[string]any{"image": server.URL + "/image.png"}}}}}}})
		case "/image.png":
			downloads++
			if r.Header.Get("Authorization") != "" || r.Header.Get("X-Tenant") != "" {
				t.Error("credentials leaked to media download")
			}
			data, _, _ := mediaBytes(providerMedia{DataURL: qwenImageTestPNG})
			w.Header().Set("Content-Type", "image/png")
			w.Write(data)
		default:
			t.Errorf("unexpected endpoint %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	config := providerConfig{Model: "qwen-image-3.0-pro", InterfaceType: "dashscope-image", BaseURL: server.URL + "/compatible-mode/v1", APIKey: "fixture", Size: "auto", Count: "1", AllowLocalChannel: true}
	ctx := withProtocolRegistry(withProviderOutboundPolicy(context.Background(), config), loadOfficialFallbackRegistry())
	result, err := runImageTask(ctx, canvasGenerationInput{Mode: "image", Prompt: "fixture", Config: config})
	if err != nil || result["mode"] != "image" || creates != 1 || downloads != 1 {
		t.Fatalf("mock lifecycle creates=%d downloads=%d err=%v", creates, downloads, err)
	}
}

func TestBailianPluginStrictTaskStates(t *testing.T) {
	adapter := migratedBailianAdapter(t, "dashscope-video")
	for _, status := range []string{"UNKNOWN", "FAILED", "CANCELED"} {
		state, err := adapter.ParsePoll(context.Background(), protocol.PollContext{TaskID: "fixture"}, []byte(`{"output":{"task_status":"`+status+`"}}`))
		if err != nil || (state.Status != protocol.StatusFailed && state.Status != protocol.StatusCancelled) {
			t.Fatalf("terminal state %s: %v %v", status, state, err)
		}
	}
	if _, err := adapter.ParsePoll(context.Background(), protocol.PollContext{TaskID: "fixture"}, []byte(`{"output":{"task_status":"SURPRISE"}}`)); err == nil {
		t.Fatal("unknown state accepted")
	}
	state, err := adapter.ParsePoll(context.Background(), protocol.PollContext{TaskID: "fixture"}, []byte(`{"output":{"task_status":"PENDING","video_url":"https://example.com/partial.mp4"}}`))
	if err != nil || state.Status != protocol.StatusPending {
		t.Fatal("pending result was treated as completed")
	}
}

func TestBailianPluginPreservesRegionalAndProxyBasePaths(t *testing.T) {
	adapter := migratedBailianAdapter(t, "dashscope-video")
	for _, base := range []string{"https://example.com/compatible-mode/v1", "https://example.com/gateway/v1/compatible-mode/v1", "https://example.com/api/v1"} {
		spec, err := adapter.BuildPoll(context.Background(), protocol.PollContext{TaskID: "task/1"})
		if err != nil {
			t.Fatal(err)
		}
		got, err := protocolRequestURL(base, spec)
		if err != nil {
			t.Fatal(err)
		}
		want := strings.TrimSuffix(strings.TrimSuffix(base, "/compatible-mode/v1"), "/api/v1") + "/api/v1/tasks/task%2F1"
		if got != want {
			t.Fatalf("base=%s got=%s want=%s", base, got, want)
		}
	}
}

func TestBailianPluginMockVideoRecoveryDoesNotCreateAgain(t *testing.T) {
	var server *httptest.Server
	creates, polls, downloads := 0, 0, 0
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/services/aigc/video-generation/video-synthesis":
			creates++
			if r.Header.Get("Authorization") != "Bearer fixture" || r.Header.Get("X-DashScope-Async") != "enable" {
				t.Error("invalid async headers")
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"output":{"task_id":"upstream-1","task_status":"PENDING"}}`))
		case "/api/v1/tasks/upstream-1":
			polls++
			if r.Header.Get("Authorization") != "Bearer fixture" {
				t.Error("missing polling credentials")
			}
			w.Header().Set("Content-Type", "application/json")
			if polls == 1 {
				w.Write([]byte(`{"output":{"task_id":"upstream-1","task_status":"RUNNING"}}`))
				return
			}
			json.NewEncoder(w).Encode(map[string]any{"output": map[string]any{"task_id": "upstream-1", "task_status": "SUCCEEDED", "video_url": server.URL + "/clip.mp4"}, "usage": map[string]any{"duration": 5}})
		case "/clip.mp4":
			downloads++
			if r.Header.Get("Authorization") != "" {
				t.Error("result download leaked key")
			}
			w.Header().Set("Content-Type", "video/mp4")
			w.Write([]byte("fixture-video"))
		default:
			t.Errorf("unexpected endpoint %s", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	config := providerConfig{Model: "wan2.7-t2v", InterfaceType: "dashscope-video", BaseURL: server.URL + "/compatible-mode/v1", APIKey: "fixture", VideoSeconds: "5", VQuality: "720p", Size: "16:9", AllowLocalChannel: true}
	input := canvasGenerationInput{Mode: "video", Prompt: "fixture", Config: config}
	ctx := withProtocolRegistry(withProviderOutboundPolicy(context.Background(), config), loadOfficialFallbackRegistry())
	adapter := migratedBailianAdapter(t, "dashscope-video")
	if _, err := runProtocolAdapterTask(ctx, input, adapter); err != nil {
		t.Fatal(err)
	}
	if creates != 1 || polls != 2 || downloads != 1 {
		t.Fatalf("create=%d poll=%d download=%d", creates, polls, downloads)
	}
	recovery := context.WithValue(ctx, providerAnalyticsKey{}, providerAnalyticsContext{ProviderRequestID: "upstream-1"})
	if _, err := runProtocolAdapterTask(recovery, input, adapter); err != nil {
		t.Fatal(err)
	}
	if creates != 1 || polls != 3 || downloads != 2 {
		t.Fatal("recovery created another task or failed to reuse results")
	}
}
