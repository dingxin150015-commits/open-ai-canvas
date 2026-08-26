package service

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
)

const wan30VideoModel = "wan3.0-video"

type wan30VideoRequest struct {
	Model      string               `json:"model"`
	Input      wan30VideoInput      `json:"input"`
	Parameters wan30VideoParameters `json:"parameters"`
}

type wan30VideoInput struct {
	Prompt string            `json:"prompt,omitempty"`
	Media  []wan30VideoMedia `json:"media,omitempty"`
}

type wan30VideoMedia struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type wan30VideoParameters struct {
	Resolution string `json:"resolution"`
	Ratio      string `json:"ratio"`
	Duration   int    `json:"duration"`
	Audio      bool   `json:"audio"`
	Seed       *int64 `json:"seed,omitempty"`
	Watermark  bool   `json:"watermark"`
}

func isWan30VideoModel(value string) bool {
	return strings.EqualFold(strings.TrimPrefix(strings.TrimSpace(value), "models/"), wan30VideoModel)
}

func buildWan30VideoRequest(input canvasGenerationInput) (map[string]interface{}, error) {
	modelName := strings.TrimPrefix(strings.TrimSpace(input.Config.Model), "models/")
	if !isWan30VideoModel(modelName) {
		return nil, fmt.Errorf("Wan 3.0 Adapter 仅支持 %s；模型 %s 尚未完成适配", wan30VideoModel, modelName)
	}
	prompt := strings.TrimSpace(input.Prompt)
	if utf8.RuneCountInString(prompt) > 20000 {
		return nil, errors.New("Wan 3.0 提示词不能超过 20000 个字符")
	}
	duration, err := wan30Duration(input.Config.VideoSeconds)
	if err != nil {
		return nil, err
	}
	resolution, err := wan30Resolution(input.Config.VQuality)
	if err != nil {
		return nil, err
	}
	ratio, err := wan30Ratio(input.Config.Size)
	if err != nil {
		return nil, err
	}
	validationInput := input
	validationInput.Config.VideoSeconds = strconv.Itoa(duration)
	validationInput.Config.VQuality = strings.ToLower(resolution)
	validationInput.Config.Size = ratio
	if err := validateVideoTask(wan30VideoCapabilityConfig(), validationInput); err != nil {
		return nil, err
	}
	media, err := buildWan30VideoMedia(input)
	if err != nil {
		return nil, err
	}
	if prompt == "" && len(media) == 0 {
		return nil, errors.New("Wan 3.0 的提示词和参考素材至少需要提供一项")
	}
	if duration > 0 {
		var inputVideoDurationMs int64
		for _, video := range input.ReferenceVideos {
			inputVideoDurationMs += max(video.DurationMs, 0)
		}
		if inputVideoDurationMs > 0 && inputVideoDurationMs+int64(duration)*1000 > 30000 {
			return nil, errors.New("Wan 3.0 的参考视频总时长与输出时长之和不能超过 30 秒")
		}
	}
	seed, err := dashScopeVideoSeed(input.Metadata)
	if err != nil {
		return nil, err
	}
	request := wan30VideoRequest{
		Model: modelName,
		Input: wan30VideoInput{Prompt: prompt, Media: media},
		Parameters: wan30VideoParameters{
			Resolution: resolution,
			Ratio:      ratio,
			Duration:   duration,
			Audio:      parseBool(input.Config.VideoGenerateAudio, true),
			Seed:       seed,
			Watermark:  parseBool(input.Config.VideoWatermark, false),
		},
	}
	return requestAsMap(request)
}

func buildWan30VideoMedia(input canvasGenerationInput) ([]wan30VideoMedia, error) {
	operation := strings.ToLower(metadataString(input.Metadata, "videoEditOperation"))
	startFrameID := metadataString(input.Metadata, "videoStartFrameNodeId")
	endFrameID := metadataString(input.Metadata, "videoEndFrameNodeId")
	if endFrameID != "" && startFrameID == "" {
		return nil, errors.New("Wan 3.0 首尾帧模式必须同时配置首帧")
	}
	if startFrameID != "" || endFrameID != "" {
		if operation != "" && operation != "image_to_video" {
			return nil, errors.New("Wan 3.0 首帧/首尾帧素材只能用于图生视频模式")
		}
		if len(input.ReferenceVideos) > 0 || len(input.ReferenceAudios) > 0 {
			return nil, errors.New("Wan 3.0 的首帧/首尾帧不能与参考视频或参考音频混用")
		}
		return wan30ExplicitFrameMedia(input, startFrameID, endFrameID)
	}
	if operation == "" {
		switch {
		case len(input.ReferenceVideos) > 0 || len(input.ReferenceAudios) > 0 || len(input.ReferenceImages) > 2:
			operation = "reference_to_video"
		case len(input.ReferenceImages) > 0:
			operation = "image_to_video"
		default:
			operation = "text_to_video"
		}
	}
	switch operation {
	case "text_to_video":
		if len(input.ReferenceImages)+len(input.ReferenceVideos)+len(input.ReferenceAudios) > 0 {
			return nil, errors.New("Wan 3.0 文生视频模式不能携带参考素材")
		}
		return nil, nil
	case "image_to_video":
		if len(input.ReferenceVideos) > 0 || len(input.ReferenceAudios) > 0 || len(input.ReferenceImages) < 1 || len(input.ReferenceImages) > 2 {
			return nil, errors.New("Wan 3.0 图生视频模式仅支持 1 张首帧或 2 张首尾帧图片")
		}
		media := make([]wan30VideoMedia, 0, len(input.ReferenceImages))
		for index, image := range input.ReferenceImages {
			mediaType := "first_frame"
			if index == 1 {
				mediaType = "last_frame"
			}
			value, err := wan30MediaURL(image)
			if err != nil {
				return nil, fmt.Errorf("读取 Wan 3.0 %s地址失败：%w", mediaType, err)
			}
			media = append(media, wan30VideoMedia{Type: mediaType, URL: value})
		}
		return media, nil
	case "reference_to_video", "audio_to_video":
		if len(input.ReferenceImages)+len(input.ReferenceVideos)+len(input.ReferenceAudios) == 0 {
			return nil, errors.New("Wan 3.0 参考生视频模式至少需要一个参考素材")
		}
		return wan30ReferenceMedia(input)
	default:
		return nil, fmt.Errorf("Wan 3.0 不支持生成模式 %s", operation)
	}
}

func wan30ExplicitFrameMedia(input canvasGenerationInput, startFrameID string, endFrameID string) ([]wan30VideoMedia, error) {
	if startFrameID == endFrameID && endFrameID != "" {
		return nil, errors.New("Wan 3.0 首帧和尾帧不能使用同一个素材")
	}
	byID := make(map[string]providerMedia, len(input.ReferenceImages))
	for _, image := range input.ReferenceImages {
		if _, exists := byID[image.ID]; image.ID == "" || exists {
			return nil, errors.New("Wan 3.0 帧素材必须具有唯一且非空的 ID")
		}
		byID[image.ID] = image
	}
	startFrame, ok := byID[startFrameID]
	if !ok {
		return nil, errors.New("已配置的 Wan 3.0 首帧参考图未包含在视频请求中")
	}
	allowed := 1
	if endFrameID != "" {
		allowed++
		if _, ok := byID[endFrameID]; !ok {
			return nil, errors.New("已配置的 Wan 3.0 尾帧参考图未包含在视频请求中")
		}
	}
	if len(input.ReferenceImages) != allowed {
		return nil, errors.New("Wan 3.0 首帧/首尾帧模式不能混入普通参考图片")
	}
	startURL, err := wan30MediaURL(startFrame)
	if err != nil {
		return nil, fmt.Errorf("读取 Wan 3.0 首帧地址失败：%w", err)
	}
	media := []wan30VideoMedia{{Type: "first_frame", URL: startURL}}
	if endFrameID != "" {
		endURL, err := wan30MediaURL(byID[endFrameID])
		if err != nil {
			return nil, fmt.Errorf("读取 Wan 3.0 尾帧地址失败：%w", err)
		}
		media = append(media, wan30VideoMedia{Type: "last_frame", URL: endURL})
	}
	return media, nil
}

func wan30ReferenceMedia(input canvasGenerationInput) ([]wan30VideoMedia, error) {
	media := make([]wan30VideoMedia, 0, len(input.ReferenceImages)+len(input.ReferenceVideos)+len(input.ReferenceAudios))
	appendMedia := func(mediaType string, values []providerMedia) error {
		for index, item := range values {
			value, err := wan30MediaURL(item)
			if err != nil {
				return fmt.Errorf("读取第 %d 个 %s 地址失败：%w", index+1, mediaType, err)
			}
			media = append(media, wan30VideoMedia{Type: mediaType, URL: value})
		}
		return nil
	}
	if err := appendMedia("reference_image", input.ReferenceImages); err != nil {
		return nil, err
	}
	if err := appendMedia("reference_video", input.ReferenceVideos); err != nil {
		return nil, err
	}
	if err := appendMedia("reference_audio", input.ReferenceAudios); err != nil {
		return nil, err
	}
	return media, nil
}

func wan30MediaURL(media providerMedia) (string, error) {
	value := strings.TrimSpace(media.URL)
	if value == "" || strings.HasPrefix(value, "asset://") {
		value = strings.TrimSpace(media.DataURL)
	}
	if strings.HasPrefix(value, "data:") {
		if !strings.Contains(value, ";base64,") {
			return "", errors.New("data URL 必须使用 base64 编码")
		}
		return value, nil
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" {
		return "", errors.New("参考素材需要公网 HTTP(S)、DashScope OSS 临时地址或 data URL")
	}
	if parsed.Scheme == "oss" && parsed.Host == "dashscope-instant" {
		return value, nil
	}
	if !isPublicMediaURL(value) {
		return "", errors.New("参考素材需要公网 HTTP(S)、DashScope OSS 临时地址或 data URL")
	}
	if _, err := ValidateOutboundURL(value); err != nil {
		return "", err
	}
	return value, nil
}

func wan30Duration(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 5, nil
	}
	duration, err := strconv.Atoi(value)
	if err != nil || (duration != -1 && (duration < 2 || duration > 30)) {
		return 0, errors.New("Wan 3.0 视频时长仅支持 -1（智能）或 2-30 秒整数")
	}
	return duration, nil
}

func wan30Resolution(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return "1080P", nil
	}
	if value == "480" || value == "720" || value == "1080" {
		value += "P"
	}
	if value != "480P" && value != "720P" && value != "1080P" {
		return "", errors.New("Wan 3.0 输出分辨率仅支持 480P、720P 或 1080P")
	}
	return value, nil
}

func wan30Ratio(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return "adaptive", nil
	}
	switch value {
	case "adaptive", "16:9", "4:3", "1:1", "3:4", "9:16":
		return value, nil
	default:
		return "", errors.New("Wan 3.0 画面比例不受支持")
	}
}
