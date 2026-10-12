package app

func wan30VideoCapabilityConfig() *VideoCapabilityConfig {
	return &VideoCapabilityConfig{
		// 官方页未声明媒体数量上限；这里采用产品侧保守上传上限，Adapter 不会静默截断。
		References: VideoReferenceConfig{
			PromptMaxChars:   20000,
			MaxImages:        9,
			MaxImageBytes:    30 * 1024 * 1024,
			MaxVideos:        3,
			MaxVideoBytes:    50 * 1024 * 1024,
			MaxVideoDuration: 30,
			MaxAudios:        1,
			MaxAudioBytes:    15 * 1024 * 1024,
			MaxAudioDuration: 30,
		},
		Duration:          VideoDurationConfig{Selection: "range", Min: 2, Max: 30, Step: 1, Default: 5, SmartSupported: true},
		Ratios:            []string{"adaptive", "16:9", "4:3", "1:1", "3:4", "9:16"},
		DefaultRatio:      "adaptive",
		Resolutions:       []string{"480p", "720p", "1080p"},
		DefaultResolution: "1080p",
		GenerateAudio:     VideoBooleanConfig{Supported: true, Default: true},
		Watermark:         VideoBooleanConfig{Supported: true, Default: false},
		Operations:        []string{"text_to_video", "image_to_video", "reference_to_video", "audio_to_video"},
		DefaultOperation:  "text_to_video",
	}
}

func wan27VideoCapabilityConfig(kind string) *VideoCapabilityConfig {
	profile := &VideoCapabilityConfig{
		References:        VideoReferenceConfig{PromptMaxChars: 5000},
		Duration:          VideoDurationConfig{Selection: "range", Min: 2, Max: 15, Step: 1, Default: 5},
		Ratios:            []string{"16:9", "9:16", "1:1", "4:3", "3:4"},
		DefaultRatio:      "16:9",
		Resolutions:       []string{"720p", "1080p"},
		DefaultResolution: "1080p",
		GenerateAudio:     VideoBooleanConfig{Supported: false, Default: false},
		Watermark:         VideoBooleanConfig{Supported: true, Default: false},
	}
	switch kind {
	case "t2v":
		profile.References.MaxAudios = 1
		profile.References.MaxAudioBytes = 15 * 1024 * 1024
		profile.References.MaxAudioDuration = 30
		profile.References.AudioMustFitOutput = false
		profile.Operations = []string{"text_to_video", "audio_to_video"}
		profile.DefaultOperation = "text_to_video"
	case "i2v":
		profile.References.MaxImages = 2
		profile.References.MinVisuals = 1
		profile.References.MaxVisuals = 2
		profile.References.MaxImageBytes = 20 * 1024 * 1024
		profile.References.MaxVideos = 1
		profile.References.MaxVideoBytes = 100 * 1024 * 1024
		profile.References.MaxVideoDuration = 10
		profile.References.MaxAudios = 1
		profile.References.MaxAudioBytes = 15 * 1024 * 1024
		profile.References.MaxAudioDuration = 30
		profile.References.AudioMustFitOutput = false
		profile.Ratios = []string{"adaptive"}
		profile.DefaultRatio = "adaptive"
		profile.Operations = []string{"image_to_video", "extend", "reference_to_video", "audio_to_video"}
		profile.DefaultOperation = "image_to_video"
	case "r2v":
		profile.References.MaxImages = 5
		profile.References.MinVisuals = 1
		profile.References.MaxVisuals = 5
		profile.References.MaxImageBytes = 20 * 1024 * 1024
		profile.References.MaxVideos = 5
		profile.References.MaxVideoBytes = 100 * 1024 * 1024
		profile.References.MaxVideoDuration = 30
		profile.References.MaxOutputDurationWithVideo = 10
		profile.References.MaxAudios = 1
		profile.References.MaxAudioBytes = 15 * 1024 * 1024
		profile.References.MaxAudioDuration = 10
		profile.Operations = []string{"reference_to_video", "image_to_video", "audio_to_video"}
		profile.DefaultOperation = "reference_to_video"
	}
	return profile
}

func happyHorseVideoCapabilityConfig(kind string) *VideoCapabilityConfig {
	profile := &VideoCapabilityConfig{
		References:        VideoReferenceConfig{PromptMaxChars: 5000},
		Duration:          VideoDurationConfig{Selection: "range", Min: 3, Max: 15, Step: 1, Default: 5},
		Ratios:            []string{"16:9", "9:16", "1:1", "4:3", "3:4", "4:5", "5:4", "9:21", "21:9"},
		DefaultRatio:      "16:9",
		Resolutions:       []string{"480p", "720p", "1080p"},
		DefaultResolution: "1080p",
		GenerateAudio:     VideoBooleanConfig{Supported: false, Default: false},
		Watermark:         VideoBooleanConfig{Supported: true, Default: true},
	}
	switch kind {
	case "t2v":
		profile.Operations = []string{"text_to_video"}
		profile.DefaultOperation = "text_to_video"
	case "i2v":
		profile.References.MinImages = 1
		profile.References.MinVisuals = 1
		profile.References.MaxVisuals = 1
		profile.References.MaxImages = 1
		profile.References.MaxImageBytes = 20 * 1024 * 1024
		profile.Ratios = []string{"adaptive"}
		profile.DefaultRatio = "adaptive"
		profile.Operations = []string{"image_to_video"}
		profile.DefaultOperation = "image_to_video"
	case "r2v":
		profile.References.MinImages = 1
		profile.References.MinVisuals = 1
		profile.References.MaxVisuals = 9
		profile.References.MaxImages = 9
		profile.References.MaxImageBytes = 20 * 1024 * 1024
		profile.Operations = []string{"reference_to_video", "image_to_video"}
		profile.DefaultOperation = "reference_to_video"
	}
	return profile
}
