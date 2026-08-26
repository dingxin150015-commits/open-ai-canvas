package service

import (
	"reflect"
	"strings"
	"testing"
)

func TestBuildWan27T2VRequestMapsOfficialFields(t *testing.T) {
	body, err := buildWan27VideoRequest(canvasGenerationInput{
		Prompt:          "cinematic cat",
		Config:          providerConfig{Model: "wan2.7-t2v", VideoSeconds: "15", VQuality: "1080p", Size: "3:4", VideoWatermark: "true"},
		ReferenceAudios: []providerMedia{{URL: "https://example.com/voice.mp3", Bytes: 1024, DurationMs: 10000}},
		Metadata:        map[string]interface{}{"videoEditOperation": "audio_to_video", "negative_prompt": "blurry", "prompt_extend": false, "seed": 7},
	})
	if err != nil {
		t.Fatal(err)
	}
	input := body["input"].(map[string]interface{})
	if input["audio_url"] != "https://example.com/voice.mp3" || input["negative_prompt"] != "blurry" {
		t.Fatalf("input = %#v", input)
	}
	parameters := body["parameters"].(map[string]interface{})
	if parameters["resolution"] != "1080P" || parameters["ratio"] != "3:4" || parameters["duration"] != float64(15) || parameters["prompt_extend"] != false || parameters["watermark"] != true || parameters["seed"] != float64(7) {
		t.Fatalf("parameters = %#v", parameters)
	}
}

func TestBuildWan27I2VRequestSupportsFramesAudioAndContinuation(t *testing.T) {
	frames, err := buildWan27VideoRequest(canvasGenerationInput{
		Prompt:          "transition",
		Config:          providerConfig{Model: "wan2.7-i2v", VideoSeconds: "10", VQuality: "720p", Size: "adaptive"},
		ReferenceImages: []providerMedia{{ID: "last", URL: "https://example.com/last.png", Width: 720, Height: 1280}, {ID: "first", URL: "https://example.com/first.png", Width: 720, Height: 1280}},
		ReferenceAudios: []providerMedia{{URL: "https://example.com/audio.mp3", DurationMs: 5000}},
		Metadata:        map[string]interface{}{"videoEditOperation": "audio_to_video", "videoStartFrameNodeId": "first", "videoEndFrameNodeId": "last"},
	})
	if err != nil {
		t.Fatal(err)
	}
	media := frames["input"].(map[string]interface{})["media"].([]interface{})
	gotTypes := []string{}
	for _, item := range media {
		gotTypes = append(gotTypes, item.(map[string]interface{})["type"].(string))
	}
	if want := []string{"first_frame", "last_frame", "driving_audio"}; !reflect.DeepEqual(gotTypes, want) {
		t.Fatalf("media types = %#v", gotTypes)
	}
	if _, exists := frames["parameters"].(map[string]interface{})["ratio"]; exists {
		t.Fatalf("Wan 2.7 i2v must not serialize ratio: %#v", frames)
	}

	continuation, err := buildWan27VideoRequest(canvasGenerationInput{
		Prompt:          "continue",
		Config:          providerConfig{Model: "wan2.7-i2v-2026-04-25", VideoSeconds: "12", VQuality: "1080p"},
		ReferenceVideos: []providerMedia{{URL: "https://example.com/clip.mp4", DurationMs: 5000}},
		ReferenceImages: []providerMedia{{ID: "last", URL: "https://example.com/last.png", Width: 720, Height: 1280}},
		Metadata:        map[string]interface{}{"videoEditOperation": "extend", "videoEndFrameNodeId": "last"},
	})
	if err != nil {
		t.Fatal(err)
	}
	media = continuation["input"].(map[string]interface{})["media"].([]interface{})
	gotTypes = gotTypes[:0]
	for _, item := range media {
		gotTypes = append(gotTypes, item.(map[string]interface{})["type"].(string))
	}
	if want := []string{"first_clip", "last_frame"}; !reflect.DeepEqual(gotTypes, want) {
		t.Fatalf("continuation media types = %#v", gotTypes)
	}
}

func TestBuildWan27R2VRequestMapsReferencesAndVoice(t *testing.T) {
	body, err := buildWan27VideoRequest(canvasGenerationInput{
		Prompt:          "Image 1 meets Video 1",
		Config:          providerConfig{Model: "wan2.7-r2v-2026-06-12", VideoSeconds: "10", VQuality: "720p", Size: "9:16"},
		ReferenceImages: []providerMedia{{ID: "frame", URL: "https://example.com/frame.png", Width: 720, Height: 1280}, {ID: "character", URL: "https://example.com/character.png", Width: 720, Height: 1280}},
		ReferenceVideos: []providerMedia{{URL: "https://example.com/reference.mp4", DurationMs: 5000, Width: 720, Height: 1280}},
		ReferenceAudios: []providerMedia{{URL: "https://example.com/voice.mp3", DurationMs: 5000}},
		Metadata:        map[string]interface{}{"videoEditOperation": "audio_to_video", "videoStartFrameNodeId": "frame"},
	})
	if err != nil {
		t.Fatal(err)
	}
	input := body["input"].(map[string]interface{})
	if input["reference_voice"] != "https://example.com/voice.mp3" {
		t.Fatalf("input = %#v", input)
	}
	media := input["media"].([]interface{})
	gotTypes := []string{}
	for _, item := range media {
		gotTypes = append(gotTypes, item.(map[string]interface{})["type"].(string))
	}
	if want := []string{"first_frame", "reference_image", "reference_video"}; !reflect.DeepEqual(gotTypes, want) {
		t.Fatalf("media types = %#v", gotTypes)
	}
	if _, exists := body["parameters"].(map[string]interface{})["ratio"]; exists {
		t.Fatalf("ratio must be omitted when first_frame is present: %#v", body)
	}
}

func TestBuildWan27RequestRejectsInvalidContracts(t *testing.T) {
	valid := canvasGenerationInput{Prompt: "test", Config: providerConfig{Model: "wan2.7-t2v", VideoSeconds: "5", VQuality: "720p", Size: "16:9"}, Metadata: map[string]interface{}{"videoEditOperation": "text_to_video"}}
	tests := []struct {
		name string
		edit func(*canvasGenerationInput)
		want string
	}{
		{name: "480p", edit: func(input *canvasGenerationInput) { input.Config.VQuality = "480p" }, want: "720P 或 1080P"},
		{name: "wrong operation", edit: func(input *canvasGenerationInput) {
			input.Config.Model = "wan2.7-i2v"
			input.ReferenceImages = []providerMedia{{URL: "https://example.com/a.png"}}
			input.Metadata = map[string]interface{}{"videoEditOperation": "text_to_video"}
		}, want: "不支持生成模式"},
		{name: "16 seconds", edit: func(input *canvasGenerationInput) { input.Config.VideoSeconds = "16" }, want: "2-15"},
		{name: "long prompt", edit: func(input *canvasGenerationInput) { input.Prompt = strings.Repeat("界", 5001) }, want: "5000"},
		{name: "t2v image", edit: func(input *canvasGenerationInput) {
			input.ReferenceImages = []providerMedia{{URL: "https://example.com/a.png"}}
		}, want: "不支持图片"},
		{name: "t2v audio would truncate", edit: func(input *canvasGenerationInput) {
			input.ReferenceAudios = []providerMedia{{URL: "https://example.com/a.mp3", DurationMs: 6000}}
		}, want: "不能长于输出视频"},
		{name: "i2v clip and audio", edit: func(input *canvasGenerationInput) {
			input.Config.Model = "wan2.7-i2v"
			input.Metadata = map[string]interface{}{"videoEditOperation": "extend"}
			input.ReferenceVideos = []providerMedia{{URL: "https://example.com/a.mp4", DurationMs: 5000}}
			input.ReferenceAudios = []providerMedia{{URL: "https://example.com/a.mp3", DurationMs: 5000}}
		}, want: "不能混用"},
		{name: "i2v gif frame", edit: func(input *canvasGenerationInput) {
			input.Config.Model = "wan2.7-i2v"
			input.Metadata = map[string]interface{}{"videoEditOperation": "image_to_video"}
			input.ReferenceImages = []providerMedia{{URL: "https://example.com/a.gif", MimeType: "image/gif", Width: 720, Height: 1280}}
		}, want: "格式"},
		{name: "r2v over five", edit: func(input *canvasGenerationInput) {
			input.Config.Model = "wan2.7-r2v"
			input.Metadata = map[string]interface{}{"videoEditOperation": "reference_to_video"}
			input.ReferenceImages = make([]providerMedia, 6)
			for index := range input.ReferenceImages {
				input.ReferenceImages[index] = providerMedia{URL: "https://example.com/a.png"}
			}
		}, want: "1-5"},
		{name: "r2v video max duration", edit: func(input *canvasGenerationInput) {
			input.Config.Model = "wan2.7-r2v"
			input.Metadata = map[string]interface{}{"videoEditOperation": "reference_to_video"}
			input.Config.VideoSeconds = "11"
			input.ReferenceVideos = []providerMedia{{URL: "https://example.com/a.mp4", DurationMs: 5000}}
		}, want: "2-10"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := valid
			input.Metadata = map[string]interface{}{"videoEditOperation": "text_to_video"}
			test.edit(&input)
			_, err := buildWan27VideoRequest(input)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func TestWan27CapabilityProfilesAreModelSpecific(t *testing.T) {
	t2v := DefaultModelCapabilityConfigForModel("dashscope-video", "wan2.7-t2v").Video
	i2v := DefaultModelCapabilityConfigForModel("dashscope-video", "wan2.7-i2v").Video
	r2v := DefaultModelCapabilityConfigForModel("dashscope-video", "wan2.7-r2v").Video
	if t2v.References.MaxImages != 0 || t2v.References.MaxAudios != 1 || t2v.DefaultResolution != "1080p" {
		t.Fatalf("t2v = %#v", t2v)
	}
	if i2v.References.MaxImages != 2 || i2v.References.MaxVideos != 1 || !reflect.DeepEqual(i2v.Ratios, []string{"adaptive"}) {
		t.Fatalf("i2v = %#v", i2v)
	}
	if r2v.References.MaxImages != 5 || r2v.References.MaxVideos != 5 || r2v.References.MaxAudioDuration != 10 || r2v.References.MaxOutputDurationWithVideo != 10 {
		t.Fatalf("r2v = %#v", r2v)
	}
	input := canvasGenerationInput{Config: providerConfig{Model: "wan2.7-r2v", VideoSeconds: "11", VQuality: "1080p", Size: "16:9"}, ReferenceVideos: []providerMedia{{DurationMs: 5000}}, Metadata: map[string]interface{}{"videoEditOperation": "reference_to_video"}}
	if err := validateVideoTask(r2v, input); err == nil || !strings.Contains(err.Error(), "最多为 10 秒") {
		t.Fatalf("validateVideoTask() error = %v", err)
	}
	spec, err := CapabilitySpecFromModelCapabilityConfig(&ModelCapabilityConfig{Version: 1, Video: r2v}, "video")
	durationWithVideo := spec.Options["videoSecondsWithReferenceVideo"]
	if err != nil || durationWithVideo.Max == nil || *durationWithVideo.Max != 10 || spec.Inputs["visual"].Max != 5 {
		t.Fatalf("capability spec = %#v, error = %v", spec, err)
	}
	match := MatchCapability(spec, ModelRequestIntent{Capability: "video", Operation: "reference_to_video", Inputs: map[string]int{"image": 3, "video": 3}, Options: map[string]any{"videoSeconds": 5, "size": "16:9", "vquality": "1080p", "videoGenerateAudio": false, "videoWatermark": false}})
	if match.Matched || !strings.Contains(strings.Join(match.Reasons, ","), "图片/视频视觉素材") {
		t.Fatalf("combined visual match = %#v", match)
	}
}
