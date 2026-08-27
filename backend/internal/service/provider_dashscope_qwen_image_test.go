package service

import (
	"strings"
	"testing"
)

const qwenImageTestPNG = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="

func TestBuildQwenImage30TextRequestOmitsAutoSize(t *testing.T) {
	body, err := buildQwenImage30Request(canvasGenerationInput{
		Prompt: "画一只猫",
		Config: providerConfig{
			Model: "qwen-image-3.0-pro", Size: "auto", Count: "6", Quality: "auto",
		},
		Metadata: map[string]interface{}{
			"negative_prompt":    "模糊",
			"prompt_extend":      true,
			"prompt_extend_mode": "agent",
			"enable_thinking":    true,
			"watermark":          true,
			"seed":               42,
		},
	})
	if err != nil {
		t.Fatalf("buildQwenImage30Request() error = %v", err)
	}
	messages := body["input"].(map[string]interface{})["messages"].([]map[string]interface{})
	content := messages[0]["content"].([]map[string]interface{})
	if len(content) != 1 || content[0]["text"] != "画一只猫" {
		t.Fatalf("text content = %#v", content)
	}
	parameters := body["parameters"].(map[string]interface{})
	if _, exists := parameters["size"]; exists {
		t.Fatalf("auto must omit size: %#v", parameters)
	}
	if parameters["n"] != 6 || parameters["prompt_extend_mode"] != "agent" || parameters["enable_thinking"] != true || parameters["watermark"] != true || parameters["seed"] != int64(42) || parameters["negative_prompt"] != "模糊" {
		t.Fatalf("parameters = %#v", parameters)
	}
}

func TestBuildQwenImage30EditRequestOrdersImagesBeforeText(t *testing.T) {
	body, err := buildQwenImage30Request(canvasGenerationInput{
		Prompt: "融合三张图片",
		Config: providerConfig{Model: "qwen-image-3.0-pro", Size: "16:9", Count: "2", Quality: "auto"},
		ReferenceImages: []providerMedia{
			{DataURL: qwenImageTestPNG},
			{DataURL: qwenImageTestPNG},
			{DataURL: qwenImageTestPNG},
		},
	})
	if err != nil {
		t.Fatalf("buildQwenImage30Request() error = %v", err)
	}
	content := body["input"].(map[string]interface{})["messages"].([]map[string]interface{})[0]["content"].([]map[string]interface{})
	if len(content) != 4 || content[0]["image"] == nil || content[1]["image"] == nil || content[2]["image"] == nil || content[3]["text"] != "融合三张图片" {
		t.Fatalf("edit content = %#v", content)
	}
	parameters := body["parameters"].(map[string]interface{})
	if parameters["size"] != "1280*720" || parameters["n"] != 2 || parameters["prompt_extend_mode"] != "direct" {
		t.Fatalf("parameters = %#v", parameters)
	}
}

func TestQwenImage30SizeContract(t *testing.T) {
	tests := []struct {
		value string
		want  string
		fail  bool
	}{
		{value: "auto", want: ""},
		{value: "1024x1536", want: "1024*1536"},
		{value: "16:9", want: "1280*720"},
		{value: "5:4", want: "1145*916"},
		{value: "100x100", fail: true},
		{value: "4096x4096", fail: true},
		{value: "9:1", fail: true},
		{value: "bad", fail: true},
	}
	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			got, err := qwenImage30Size(test.value)
			if test.fail {
				if err == nil {
					t.Fatalf("qwenImage30Size(%q) = %q, want error", test.value, got)
				}
				return
			}
			if err != nil || got != test.want {
				t.Fatalf("qwenImage30Size(%q) = %q, %v; want %q", test.value, got, err, test.want)
			}
		})
	}
}

func TestBuildQwenImage30RejectsUnsupportedInputs(t *testing.T) {
	base := canvasGenerationInput{
		Prompt: "test",
		Config: providerConfig{Model: "qwen-image-3.0-pro", Size: "1024x1024", Count: "1", Quality: "auto"},
	}
	tests := []struct {
		name string
		edit func(*canvasGenerationInput)
		want string
	}{
		{name: "wrong model", edit: func(input *canvasGenerationInput) { input.Config.Model = "wan2.7-image-pro" }, want: "不属于"},
		{name: "mask", edit: func(input *canvasGenerationInput) { input.Mask = &providerMedia{DataURL: qwenImageTestPNG} }, want: "蒙版"},
		{name: "four images", edit: func(input *canvasGenerationInput) {
			input.ReferenceImages = []providerMedia{{DataURL: qwenImageTestPNG}, {DataURL: qwenImageTestPNG}, {DataURL: qwenImageTestPNG}, {DataURL: qwenImageTestPNG}}
		}, want: "最多支持 3 张"},
		{name: "quality", edit: func(input *canvasGenerationInput) { input.Config.Quality = "high" }, want: "质量"},
		{name: "transparent", edit: func(input *canvasGenerationInput) { input.Config.TransparentBackground = "true" }, want: "透明"},
		{name: "count", edit: func(input *canvasGenerationInput) { input.Config.Count = "7" }, want: "1-6"},
		{name: "agent edit", edit: func(input *canvasGenerationInput) {
			input.ReferenceImages = []providerMedia{{DataURL: qwenImageTestPNG}}
			input.Metadata = map[string]interface{}{"prompt_extend_mode": "agent"}
		}, want: "编辑不支持 agent"},
		{name: "thinking without extend", edit: func(input *canvasGenerationInput) {
			input.Metadata = map[string]interface{}{"prompt_extend": false, "enable_thinking": true}
		}, want: "不能开启思考"},
		{name: "negative prompt", edit: func(input *canvasGenerationInput) {
			input.Metadata = map[string]interface{}{"negative_prompt": strings.Repeat("界", 501)}
		}, want: "500"},
		{name: "bad seed", edit: func(input *canvasGenerationInput) { input.Metadata = map[string]interface{}{"seed": -1} }, want: "0-2147483647"},
		{name: "bad size", edit: func(input *canvasGenerationInput) { input.Config.Size = "100x100" }, want: "总像素"},
		{name: "bad mime", edit: func(input *canvasGenerationInput) {
			input.ReferenceImages = []providerMedia{{DataURL: "data:text/plain;base64,dGVzdA=="}}
		}, want: "格式"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := base
			test.edit(&input)
			if _, err := buildQwenImage30Request(input); err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestQwenImage30CapabilityIsConservative(t *testing.T) {
	profile := DefaultImageCapabilityConfig("dashscope-image", "qwen-image-3.0-pro")
	if profile.References.MaxImages != 3 || profile.References.MaxImageBytes != qwenImageMaxBytes || profile.References.MaskSupported {
		t.Fatalf("references = %#v", profile.References)
	}
	if profile.Size.Default != "auto" || !profile.Size.AllowCustom || profile.Quality.Supported || profile.TransparentBackground.Supported || profile.ResponseFormat.Supported || profile.OutputFormat.Supported || profile.MaxOutputs != 6 {
		t.Fatalf("profile = %#v", profile)
	}
}

func TestValidateQwenImage30TaskRejectsInvalidProductValues(t *testing.T) {
	profile := qwenImage30CapabilityConfig()
	base := canvasGenerationInput{
		Prompt: "test",
		Config: providerConfig{
			Model: "qwen-image-3.0-pro", Size: "1024x1024", Count: "1", Quality: "auto",
		},
	}
	tests := []struct {
		name string
		edit func(*canvasGenerationInput)
		want string
	}{
		{name: "invalid size", edit: func(input *canvasGenerationInput) { input.Config.Size = "100x100" }, want: "总像素"},
		{name: "invalid count", edit: func(input *canvasGenerationInput) { input.Config.Count = "many" }, want: "正整数"},
		{name: "too many outputs", edit: func(input *canvasGenerationInput) { input.Config.Count = "7" }, want: "最多生成 6"},
		{name: "quality", edit: func(input *canvasGenerationInput) { input.Config.Quality = "high" }, want: "质量"},
		{name: "transparent", edit: func(input *canvasGenerationInput) { input.Config.TransparentBackground = "true" }, want: "透明"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := base
			test.edit(&input)
			if err := validateImageTask(profile, input); err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}
