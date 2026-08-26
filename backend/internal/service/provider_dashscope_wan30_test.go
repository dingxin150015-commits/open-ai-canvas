package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestBuildWan30VideoRequestTextToVideo(t *testing.T) {
	body, err := buildWan30VideoRequest(canvasGenerationInput{
		Prompt: "月光下奔跑的小猫",
		Config: providerConfig{
			Model: "wan3.0-video", VQuality: "1080p", Size: "adaptive", VideoSeconds: "30",
			VideoGenerateAudio: "false", VideoWatermark: "true",
		},
		Metadata: map[string]interface{}{"videoEditOperation": "text_to_video", "seed": float64(42)},
	})
	if err != nil {
		t.Fatal(err)
	}
	input := body["input"].(map[string]interface{})
	if input["prompt"] != "月光下奔跑的小猫" || input["media"] != nil {
		t.Fatalf("input = %#v", input)
	}
	parameters := body["parameters"].(map[string]interface{})
	if parameters["resolution"] != "1080P" || parameters["ratio"] != "adaptive" || parameters["duration"] != float64(30) || parameters["audio"] != false || parameters["watermark"] != true || parameters["seed"] != float64(42) {
		t.Fatalf("parameters = %#v", parameters)
	}
	if _, exists := parameters["prompt_extend"]; exists {
		t.Fatalf("official Wan 3.0 schema does not define prompt_extend: %#v", parameters)
	}
}

func TestBuildWan30VideoRequestFramesUseExplicitRoles(t *testing.T) {
	body, err := buildWan30VideoRequest(canvasGenerationInput{
		Prompt: "从微笑过渡到大笑",
		Config: providerConfig{Model: "wan3.0-video", VideoSeconds: "5", VQuality: "480", Size: "16:9"},
		ReferenceImages: []providerMedia{
			{ID: "other-order-last", URL: "https://example.com/last.png"},
			{ID: "other-order-first", URL: "https://example.com/first.png"},
		},
		Metadata: map[string]interface{}{
			"videoEditOperation": "image_to_video", "videoStartFrameNodeId": "other-order-first", "videoEndFrameNodeId": "other-order-last",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	media := body["input"].(map[string]interface{})["media"].([]interface{})
	want := []struct{ mediaType, mediaURL string }{{"first_frame", "https://example.com/first.png"}, {"last_frame", "https://example.com/last.png"}}
	for index, expected := range want {
		item := media[index].(map[string]interface{})
		if item["type"] != expected.mediaType || item["url"] != expected.mediaURL {
			t.Fatalf("media[%d] = %#v", index, item)
		}
	}
}

func TestBuildWan30VideoRequestReferenceMedia(t *testing.T) {
	body, err := buildWan30VideoRequest(canvasGenerationInput{
		Prompt:          "让图1和视频1中的人物对话",
		Config:          providerConfig{Model: "wan3.0-video", VideoSeconds: "-1", VQuality: "720P", Size: "9:16"},
		ReferenceImages: []providerMedia{{ID: "image", DataURL: "data:image/png;base64,aW1hZ2U="}},
		ReferenceVideos: []providerMedia{{ID: "video", URL: "https://example.com/video.mp4", DurationMs: 12000}},
		ReferenceAudios: []providerMedia{{ID: "audio", URL: "oss://dashscope-instant/audio/reference.wav"}},
		Metadata:        map[string]interface{}{"videoEditOperation": "reference_to_video"},
	})
	if err != nil {
		t.Fatal(err)
	}
	media := body["input"].(map[string]interface{})["media"].([]interface{})
	gotTypes := make([]string, 0, len(media))
	for _, raw := range media {
		gotTypes = append(gotTypes, raw.(map[string]interface{})["type"].(string))
	}
	if want := []string{"reference_image", "reference_video", "reference_audio"}; !reflect.DeepEqual(gotTypes, want) {
		t.Fatalf("media types = %#v, want %#v", gotTypes, want)
	}
	parameters := body["parameters"].(map[string]interface{})
	if parameters["duration"] != float64(-1) || parameters["audio"] != true {
		t.Fatalf("parameters = %#v", parameters)
	}
}

func TestBuildWan30VideoRequestRejectsInvalidContracts(t *testing.T) {
	valid := canvasGenerationInput{Prompt: "test", Config: providerConfig{Model: "wan3.0-video", VideoSeconds: "5", VQuality: "720p", Size: "adaptive"}, Metadata: map[string]interface{}{"videoEditOperation": "text_to_video"}}
	tests := []struct {
		name string
		edit func(*canvasGenerationInput)
		want string
	}{
		{name: "prime remains planned", edit: func(input *canvasGenerationInput) { input.Config.Model = "wan3.0-video-prime" }, want: "尚未完成适配"},
		{name: "duration below range", edit: func(input *canvasGenerationInput) { input.Config.VideoSeconds = "1" }, want: "2-30"},
		{name: "duration above range", edit: func(input *canvasGenerationInput) { input.Config.VideoSeconds = "31" }, want: "2-30"},
		{name: "unsupported resolution", edit: func(input *canvasGenerationInput) { input.Config.VQuality = "1440p" }, want: "分辨率"},
		{name: "unsupported ratio", edit: func(input *canvasGenerationInput) { input.Config.Size = "21:9" }, want: "比例"},
		{name: "prompt too long", edit: func(input *canvasGenerationInput) { input.Prompt = strings.Repeat("界", 20001) }, want: "20000"},
		{name: "text mode with media", edit: func(input *canvasGenerationInput) {
			input.ReferenceImages = []providerMedia{{URL: "https://example.com/a.png"}}
		}, want: "不能携带"},
		{name: "last frame without first", edit: func(input *canvasGenerationInput) {
			input.Metadata = map[string]interface{}{"videoEditOperation": "image_to_video", "videoEndFrameNodeId": "last"}
			input.ReferenceImages = []providerMedia{{ID: "last", URL: "https://example.com/last.png"}}
		}, want: "必须同时配置首帧"},
		{name: "mixed frame and reference video", edit: func(input *canvasGenerationInput) {
			input.Metadata = map[string]interface{}{"videoEditOperation": "image_to_video", "videoStartFrameNodeId": "first"}
			input.ReferenceImages = []providerMedia{{ID: "first", URL: "https://example.com/first.png"}}
			input.ReferenceVideos = []providerMedia{{URL: "https://example.com/video.mp4"}}
		}, want: "不能与参考视频"},
		{name: "video duration budget", edit: func(input *canvasGenerationInput) {
			input.Metadata = map[string]interface{}{"videoEditOperation": "reference_to_video"}
			input.Config.VideoSeconds = "20"
			input.ReferenceVideos = []providerMedia{{URL: "https://example.com/video.mp4", DurationMs: 11000}}
		}, want: "不能超过 30 秒"},
		{name: "asset URL not hydrated", edit: func(input *canvasGenerationInput) {
			input.Metadata = map[string]interface{}{"videoEditOperation": "reference_to_video"}
			input.ReferenceImages = []providerMedia{{URL: "asset://resource-1"}}
		}, want: "公网 HTTP"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := valid
			input.Metadata = map[string]interface{}{"videoEditOperation": "text_to_video"}
			test.edit(&input)
			_, err := buildWan30VideoRequest(input)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func TestWan30CapabilityMatchesOfficialContract(t *testing.T) {
	profile := DefaultModelCapabilityConfigForModel("dashscope-video", "wan3.0-video").Video
	if profile.Duration.Min != 2 || profile.Duration.Max != 30 || !profile.Duration.SmartSupported || profile.DefaultResolution != "1080p" || profile.DefaultRatio != "adaptive" {
		t.Fatalf("Wan 3.0 profile = %#v", profile)
	}
	for _, seconds := range []string{"-1", "2", "30"} {
		input := canvasGenerationInput{Config: providerConfig{Model: "wan3.0-video", VideoSeconds: seconds, VQuality: "1080p", Size: "adaptive"}, Metadata: map[string]interface{}{"videoEditOperation": "text_to_video"}}
		if err := validateVideoTask(profile, input); err != nil {
			t.Fatalf("validateVideoTask(%s) error = %v", seconds, err)
		}
	}
	input := canvasGenerationInput{Config: providerConfig{Model: "wan3.0-video", VideoSeconds: "31", VQuality: "1080p", Size: "adaptive"}, Metadata: map[string]interface{}{"videoEditOperation": "text_to_video"}}
	if err := validateVideoTask(profile, input); err == nil {
		t.Fatal("31 seconds must be rejected")
	}
	normalized, err := NormalizeModelCapabilityConfig("video", "dashscope-video", &ModelCapabilityConfig{Version: 1, Video: profile})
	if err != nil {
		t.Fatalf("NormalizeModelCapabilityConfig() error = %v", err)
	}
	spec, err := CapabilitySpecFromModelCapabilityConfig(normalized, "video")
	if err != nil {
		t.Fatalf("CapabilitySpecFromModelCapabilityConfig() error = %v", err)
	}
	values := spec.Options["videoSeconds"].Values
	if len(values) != 30 || values[0] != -1 || values[len(values)-1] != 30 {
		t.Fatalf("videoSeconds capability values = %#v", values)
	}
}

func TestDashScopeDispatcherAndSubmitUseWan30NativeContract(t *testing.T) {
	t.Setenv("CANVAS_ALLOW_PRIVATE_UPSTREAMS", "true")
	var received map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost || request.URL.Path != "/api/v1/services/aigc/video-generation/video-synthesis" {
			http.NotFound(w, request)
			return
		}
		if request.Header.Get("X-DashScope-Async") != "enable" || request.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("headers = %#v", request.Header)
		}
		if err := json.NewDecoder(request.Body).Decode(&received); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"request_id":"request-1","output":{"task_id":"task-1","task_status":"PENDING"}}`))
	}))
	defer server.Close()

	input := canvasGenerationInput{
		Prompt:   "test",
		Config:   providerConfig{BaseURL: server.URL, AllowLocalChannel: true, APIKey: "test-key", Model: "wan3.0-video", VideoSeconds: "30", VQuality: "480p", Size: "adaptive"},
		Metadata: map[string]interface{}{"videoEditOperation": "text_to_video"},
	}
	body, err := buildDashScopeVideoRequest(input)
	if err != nil {
		t.Fatal(err)
	}
	parameters := body["parameters"].(map[string]interface{})
	if parameters["duration"] != float64(30) || parameters["resolution"] != "480P" || parameters["audio"] != true {
		t.Fatalf("dispatcher body = %#v", body)
	}
	taskID, err := submitDashScopeVideoTask(context.Background(), input.Config, body)
	if err != nil {
		t.Fatal(err)
	}
	if taskID != "task-1" || received["model"] != "wan3.0-video" {
		t.Fatalf("taskID=%q received=%#v", taskID, received)
	}
}

func TestDownloadDashScopeVideoDoesNotLeakUpstreamCredentials(t *testing.T) {
	t.Setenv("CANVAS_ALLOW_PRIVATE_UPSTREAMS", "true")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if got := request.Header.Get("Authorization"); got != "" {
			t.Errorf("result download Authorization = %q, want empty", got)
		}
		if got := request.Header.Get("X-Gateway-Tenant"); got != "" {
			t.Errorf("result download X-Gateway-Tenant = %q, want empty", got)
		}
		w.Header().Set("Content-Type", "video/mp4")
		_, _ = w.Write([]byte("video-bytes"))
	}))
	defer server.Close()

	dataURL, mimeType, err := downloadDashScopeVideo(context.Background(), providerConfig{
		APIKey:  "secret-key",
		Headers: []OutboundHeader{{Name: "X-Gateway-Tenant", Value: "secret-tenant"}},
	}, server.URL+"/result.mp4")
	if err != nil {
		t.Fatal(err)
	}
	if mimeType != "video/mp4" || !strings.HasPrefix(dataURL, "data:video/mp4;base64,") {
		t.Fatalf("download result mimeType=%q dataURL=%q", mimeType, dataURL)
	}
}

func TestDashScopeVideoTaskOutcomeCoversOfficialTerminalStates(t *testing.T) {
	succeeded := &dashScopeVideoTaskResult{}
	succeeded.Output.TaskStatus = "SUCCEEDED"
	succeeded.Output.VideoURL = "https://example.com/result.mp4"
	if done, err := dashScopeVideoTaskOutcome("wan3.0-video", succeeded); !done || err != nil {
		t.Fatalf("SUCCEEDED done=%v err=%v", done, err)
	}
	for _, test := range []struct{ status, want string }{
		{status: "CANCELED", want: "已取消"},
		{status: "UNKNOWN", want: "已过期"},
		{status: "FAILED", want: "视频生成失败"},
	} {
		result := &dashScopeVideoTaskResult{}
		result.Output.TaskStatus = test.status
		if done, err := dashScopeVideoTaskOutcome("wan3.0-video", result); done || err == nil || !strings.Contains(err.Error(), test.want) {
			t.Fatalf("%s done=%v err=%v", test.status, done, err)
		}
	}
	for _, status := range []string{"PENDING", "RUNNING"} {
		result := &dashScopeVideoTaskResult{}
		result.Output.TaskStatus = status
		if done, err := dashScopeVideoTaskOutcome("wan3.0-video", result); done || err != nil {
			t.Fatalf("%s done=%v err=%v", status, done, err)
		}
	}
}
