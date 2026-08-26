package service

import (
	"reflect"
	"strings"
	"testing"
)

func TestBuildHappyHorseRequestsMatchOfficialContracts(t *testing.T) {
	t2v, err := buildHappyHorseVideoRequest(canvasGenerationInput{
		Prompt:   "cardboard city",
		Config:   providerConfig{Model: "happyhorse-1.1-t2v", VideoSeconds: "15", VQuality: "480p", Size: "9:21"},
		Metadata: map[string]interface{}{"videoEditOperation": "text_to_video", "seed": 9},
	})
	if err != nil {
		t.Fatal(err)
	}
	parameters := t2v["parameters"].(map[string]interface{})
	if parameters["resolution"] != "480P" || parameters["ratio"] != "9:21" || parameters["duration"] != float64(15) || parameters["watermark"] != true || parameters["seed"] != float64(9) {
		t.Fatalf("t2v parameters = %#v", parameters)
	}

	i2v, err := buildHappyHorseVideoRequest(canvasGenerationInput{
		Config:          providerConfig{Model: "happyhorse-1.1-i2v", VideoSeconds: "3", VQuality: "1080p", Size: "adaptive", VideoWatermark: "false"},
		ReferenceImages: []providerMedia{{ID: "first", DataURL: "data:image/png;base64,aW1hZ2U=", Width: 600, Height: 1000}},
		Metadata:        map[string]interface{}{"videoEditOperation": "image_to_video", "videoStartFrameNodeId": "first"},
	})
	if err != nil {
		t.Fatal(err)
	}
	media := i2v["input"].(map[string]interface{})["media"].([]interface{})
	if len(media) != 1 || media[0].(map[string]interface{})["type"] != "first_frame" {
		t.Fatalf("i2v media = %#v", media)
	}
	if _, exists := i2v["parameters"].(map[string]interface{})["ratio"]; exists {
		t.Fatalf("i2v ratio must follow first frame: %#v", i2v)
	}

	images := make([]providerMedia, 9)
	for index := range images {
		images[index] = providerMedia{DataURL: "data:image/png;base64,aW1hZ2U=", Width: 720, Height: 1280}
	}
	r2v, err := buildHappyHorseVideoRequest(canvasGenerationInput{
		Prompt:          "[Image 1] meets [Image 2]",
		Config:          providerConfig{Model: "happyhorse-1.1-r2v", VideoSeconds: "5", VQuality: "720p", Size: "4:5"},
		ReferenceImages: images,
		Metadata:        map[string]interface{}{"videoEditOperation": "reference_to_video"},
	})
	if err != nil {
		t.Fatal(err)
	}
	media = r2v["input"].(map[string]interface{})["media"].([]interface{})
	if len(media) != 9 || media[0].(map[string]interface{})["type"] != "reference_image" {
		t.Fatalf("r2v media = %#v", media)
	}
}

func TestBuildHappyHorseRequestRejectsInvalidContracts(t *testing.T) {
	valid := canvasGenerationInput{Prompt: "test", Config: providerConfig{Model: "happyhorse-1.1-t2v", VideoSeconds: "5", VQuality: "720p", Size: "16:9"}, Metadata: map[string]interface{}{"videoEditOperation": "text_to_video"}}
	tests := []struct {
		name string
		edit func(*canvasGenerationInput)
		want string
	}{
		{name: "old model planned", edit: func(input *canvasGenerationInput) { input.Config.Model = "happyhorse-1.0-t2v" }, want: "不支持模型"},
		{name: "wrong operation", edit: func(input *canvasGenerationInput) {
			input.Metadata = map[string]interface{}{"videoEditOperation": "image_to_video"}
		}, want: "不支持生成模式"},
		{name: "duration 2", edit: func(input *canvasGenerationInput) { input.Config.VideoSeconds = "2" }, want: "3-15"},
		{name: "bad ratio", edit: func(input *canvasGenerationInput) { input.Config.Size = "2:1" }, want: "比例"},
		{name: "weighted prompt", edit: func(input *canvasGenerationInput) { input.Prompt = strings.Repeat("界", 2501) }, want: "2500"},
		{name: "i2v two images", edit: func(input *canvasGenerationInput) {
			input.Config.Model = "happyhorse-1.1-i2v"
			input.Metadata = map[string]interface{}{"videoEditOperation": "image_to_video"}
			input.ReferenceImages = []providerMedia{{DataURL: "data:image/png;base64,YQ=="}, {DataURL: "data:image/png;base64,Yg=="}}
		}, want: "有且仅支持 1 张"},
		{name: "i2v narrow image", edit: func(input *canvasGenerationInput) {
			input.Config.Model = "happyhorse-1.1-i2v"
			input.Metadata = map[string]interface{}{"videoEditOperation": "image_to_video"}
			input.ReferenceImages = []providerMedia{{DataURL: "data:image/png;base64,YQ==", Width: 300, Height: 900}}
		}, want: "1:2.5"},
		{name: "i2v bmp image", edit: func(input *canvasGenerationInput) {
			input.Config.Model = "happyhorse-1.1-i2v"
			input.Metadata = map[string]interface{}{"videoEditOperation": "image_to_video"}
			input.ReferenceImages = []providerMedia{{DataURL: "data:image/bmp;base64,YQ==", MimeType: "image/bmp", Width: 600, Height: 600}}
		}, want: "格式"},
		{name: "r2v video", edit: func(input *canvasGenerationInput) {
			input.Config.Model = "happyhorse-1.1-r2v"
			input.Metadata = map[string]interface{}{"videoEditOperation": "reference_to_video"}
			input.ReferenceImages = []providerMedia{{DataURL: "data:image/png;base64,YQ==", Width: 720, Height: 1280}}
			input.ReferenceVideos = []providerMedia{{URL: "https://example.com/a.mp4"}}
		}, want: "仅支持 1-9 张"},
		{name: "r2v small image", edit: func(input *canvasGenerationInput) {
			input.Config.Model = "happyhorse-1.1-r2v"
			input.Metadata = map[string]interface{}{"videoEditOperation": "reference_to_video"}
			input.ReferenceImages = []providerMedia{{DataURL: "data:image/png;base64,YQ==", Width: 399, Height: 720}}
		}, want: "短边不能小于 400"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := valid
			input.Metadata = map[string]interface{}{"videoEditOperation": "text_to_video"}
			test.edit(&input)
			_, err := buildHappyHorseVideoRequest(input)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want containing %q", err, test.want)
			}
		})
	}
}

func TestHappyHorseCapabilityProfilesAreModelSpecific(t *testing.T) {
	t2v := DefaultModelCapabilityConfigForModel("dashscope-video", "happyhorse-1.1-t2v").Video
	i2v := DefaultModelCapabilityConfigForModel("dashscope-video", "happyhorse-1.1-i2v").Video
	r2v := DefaultModelCapabilityConfigForModel("dashscope-video", "happyhorse-1.1-r2v").Video
	if t2v.Duration.Min != 3 || t2v.Watermark.Default != true || len(t2v.Ratios) != 9 {
		t.Fatalf("t2v = %#v", t2v)
	}
	if i2v.References.MinImages != 1 || i2v.References.MaxImages != 1 || !reflect.DeepEqual(i2v.Ratios, []string{"adaptive"}) {
		t.Fatalf("i2v = %#v", i2v)
	}
	if r2v.References.MinImages != 1 || r2v.References.MaxImages != 9 || r2v.References.MaxVideos != 0 {
		t.Fatalf("r2v = %#v", r2v)
	}
}
