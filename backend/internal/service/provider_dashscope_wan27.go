package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type wan27VideoInput struct {
	Prompt         string                   `json:"prompt,omitempty"`
	NegativePrompt string                   `json:"negative_prompt,omitempty"`
	AudioURL       string                   `json:"audio_url,omitempty"`
	Media          []dashScopeContractMedia `json:"media,omitempty"`
	ReferenceVoice string                   `json:"reference_voice,omitempty"`
}

type dashScopeContractMedia struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type wan27VideoParameters struct {
	Resolution   string `json:"resolution"`
	Ratio        string `json:"ratio,omitempty"`
	Duration     int    `json:"duration"`
	PromptExtend bool   `json:"prompt_extend"`
	Watermark    bool   `json:"watermark"`
	Seed         *int64 `json:"seed,omitempty"`
}

type wan27VideoRequest struct {
	Model      string               `json:"model"`
	Input      wan27VideoInput      `json:"input"`
	Parameters wan27VideoParameters `json:"parameters"`
}

func wan27VideoKind(value string) string {
	value = strings.ToLower(strings.TrimPrefix(strings.TrimSpace(value), "models/"))
	switch value {
	case "wan2.7-t2v", "wan2.7-t2v-2026-06-12", "wan2.7-t2v-2026-04-25":
		return "t2v"
	case "wan2.7-i2v", "wan2.7-i2v-2026-04-25":
		return "i2v"
	case "wan2.7-r2v", "wan2.7-r2v-2026-06-12":
		return "r2v"
	default:
		return ""
	}
}

func isWan27ReadyVideoModel(value string) bool { return wan27VideoKind(value) != "" }

func buildWan27VideoRequest(input canvasGenerationInput) (map[string]interface{}, error) {
	modelName := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(input.Config.Model), "models/"))
	kind := wan27VideoKind(modelName)
	if kind == "" {
		return nil, fmt.Errorf("Wan 2.7 Adapter 不支持模型 %s", modelName)
	}
	prompt := strings.TrimSpace(input.Prompt)
	if err := validatePromptRunes(prompt, 5000, "Wan 2.7", kind != "i2v"); err != nil {
		return nil, err
	}
	negativePrompt := metadataString(input.Metadata, "negative_prompt")
	if err := validateNegativePrompt(negativePrompt, 500, "Wan 2.7"); err != nil {
		return nil, err
	}
	duration, err := wan27Duration(input.Config.VideoSeconds, 15)
	if err != nil {
		return nil, err
	}
	resolution, err := wan27Resolution(input.Config.VQuality)
	if err != nil {
		return nil, err
	}
	promptExtend, err := dashScopeMetadataBool(input.Metadata, "prompt_extend", true)
	if err != nil {
		return nil, err
	}
	seed, err := dashScopeVideoSeed(input.Metadata)
	if err != nil {
		return nil, err
	}
	request := wan27VideoRequest{
		Model: modelName,
		Input: wan27VideoInput{Prompt: prompt, NegativePrompt: negativePrompt},
		Parameters: wan27VideoParameters{
			Resolution: resolution, Duration: duration, PromptExtend: promptExtend,
			Watermark: parseBool(input.Config.VideoWatermark, false), Seed: seed,
		},
	}
	switch kind {
	case "t2v":
		if err := buildWan27T2VInput(input, duration, &request); err != nil {
			return nil, err
		}
		ratio, err := wan27Ratio(input.Config.Size)
		if err != nil {
			return nil, err
		}
		request.Parameters.Ratio = ratio
	case "i2v":
		if err := validateDashScopeVideoOperation(input, "Wan 2.7 图生视频", "image_to_video", "extend", "reference_to_video", "audio_to_video"); err != nil {
			return nil, err
		}
		media, err := buildWan27I2VMedia(input, duration)
		if err != nil {
			return nil, err
		}
		request.Input.Media = media
	case "r2v":
		if err := validateDashScopeVideoOperation(input, "Wan 2.7 参考生视频", "reference_to_video", "image_to_video", "audio_to_video"); err != nil {
			return nil, err
		}
		media, referenceVoice, hasFirstFrame, err := buildWan27R2VMedia(input, duration)
		if err != nil {
			return nil, err
		}
		request.Input.Media = media
		request.Input.ReferenceVoice = referenceVoice
		if !hasFirstFrame {
			ratio, err := wan27Ratio(input.Config.Size)
			if err != nil {
				return nil, err
			}
			request.Parameters.Ratio = ratio
		}
	}
	return requestAsMap(request)
}

func buildWan27T2VInput(input canvasGenerationInput, outputDuration int, request *wan27VideoRequest) error {
	if len(input.ReferenceImages) > 0 || len(input.ReferenceVideos) > 0 || len(input.ReferenceAudios) > 1 {
		return errors.New("Wan 2.7 文生视频只支持最多 1 个驱动音频，不支持图片或视频输入")
	}
	operation := strings.ToLower(metadataString(input.Metadata, "videoEditOperation"))
	if operation != "" && operation != "text_to_video" && operation != "audio_to_video" {
		return fmt.Errorf("Wan 2.7 文生视频不支持生成模式 %s", operation)
	}
	if len(input.ReferenceAudios) == 0 {
		return nil
	}
	audio := input.ReferenceAudios[0]
	if err := validateProviderMediaMIME(audio, []string{"audio/wav", "audio/x-wav", "audio/mpeg", "audio/mp3"}, "Wan 2.7 驱动音频"); err != nil {
		return err
	}
	if err := validateProviderMediaSize(audio, 15*1024*1024, "Wan 2.7 驱动音频"); err != nil {
		return err
	}
	if err := validateProviderMediaDuration(audio, 2, 30, "Wan 2.7 驱动音频"); err != nil {
		return err
	}
	if audio.DurationMs > int64(outputDuration)*1000 {
		return errors.New("Wan 2.7 驱动音频不能长于输出视频，否则上游会截断音频")
	}
	value, err := dashScopePublicMediaURL(audio, false)
	if err != nil {
		return fmt.Errorf("读取 Wan 2.7 驱动音频地址失败：%w", err)
	}
	request.Input.AudioURL = value
	return nil
}

func buildWan27I2VMedia(input canvasGenerationInput, outputDuration int) ([]dashScopeContractMedia, error) {
	if len(input.ReferenceVideos) > 1 || len(input.ReferenceAudios) > 1 || len(input.ReferenceImages) > 2 {
		return nil, errors.New("Wan 2.7 图生视频最多支持 2 张帧图、1 个续写视频和 1 个驱动音频")
	}
	if len(input.ReferenceVideos) == 1 {
		if len(input.ReferenceAudios) > 0 || len(input.ReferenceImages) > 1 || metadataString(input.Metadata, "videoStartFrameNodeId") != "" {
			return nil, errors.New("Wan 2.7 视频续写只支持 first_clip，或 first_clip 加 1 张尾帧，不能混用首帧或驱动音频")
		}
		clip := input.ReferenceVideos[0]
		if err := validateProviderMediaMIME(clip, []string{"video/mp4", "video/quicktime"}, "Wan 2.7 续写视频"); err != nil {
			return nil, err
		}
		if err := validateProviderMediaSize(clip, 100*1024*1024, "Wan 2.7 续写视频"); err != nil {
			return nil, err
		}
		if err := validateProviderMediaDuration(clip, 2, 10, "Wan 2.7 续写视频"); err != nil {
			return nil, err
		}
		if clip.DurationMs > 0 && int64(outputDuration)*1000 <= clip.DurationMs {
			return nil, errors.New("Wan 2.7 视频续写的总输出时长必须大于输入视频时长")
		}
		clipURL, err := dashScopePublicMediaURL(clip, false)
		if err != nil {
			return nil, fmt.Errorf("读取 Wan 2.7 续写视频地址失败：%w", err)
		}
		media := []dashScopeContractMedia{{Type: "first_clip", URL: clipURL}}
		if len(input.ReferenceImages) == 1 {
			if endID := metadataString(input.Metadata, "videoEndFrameNodeId"); endID != "" && input.ReferenceImages[0].ID != endID {
				return nil, errors.New("已配置的 Wan 2.7 尾帧未包含在请求中")
			}
			lastURL, err := wan27FrameURL(input.ReferenceImages[0], "尾帧")
			if err != nil {
				return nil, err
			}
			media = append(media, dashScopeContractMedia{Type: "last_frame", URL: lastURL})
		}
		return media, nil
	}
	if len(input.ReferenceImages) < 1 {
		return nil, errors.New("Wan 2.7 图生视频必须提供首帧图片或续写视频")
	}
	startID := metadataString(input.Metadata, "videoStartFrameNodeId")
	endID := metadataString(input.Metadata, "videoEndFrameNodeId")
	ordered, err := orderedFrameMedia(input.ReferenceImages, startID, endID, "Wan 2.7")
	if err != nil {
		return nil, err
	}
	media := make([]dashScopeContractMedia, 0, len(ordered)+len(input.ReferenceAudios))
	for index, image := range ordered {
		mediaType := "first_frame"
		label := "首帧"
		if index == 1 {
			mediaType, label = "last_frame", "尾帧"
		}
		value, err := wan27FrameURL(image, label)
		if err != nil {
			return nil, err
		}
		media = append(media, dashScopeContractMedia{Type: mediaType, URL: value})
	}
	if len(input.ReferenceAudios) == 1 {
		audio := input.ReferenceAudios[0]
		if err := validateProviderMediaMIME(audio, []string{"audio/wav", "audio/x-wav", "audio/mpeg", "audio/mp3"}, "Wan 2.7 驱动音频"); err != nil {
			return nil, err
		}
		if err := validateProviderMediaSize(audio, 15*1024*1024, "Wan 2.7 驱动音频"); err != nil {
			return nil, err
		}
		if err := validateProviderMediaDuration(audio, 2, 30, "Wan 2.7 驱动音频"); err != nil {
			return nil, err
		}
		if audio.DurationMs > int64(outputDuration)*1000 {
			return nil, errors.New("Wan 2.7 驱动音频不能长于输出视频，否则上游会截断音频")
		}
		value, err := dashScopePublicMediaURL(audio, false)
		if err != nil {
			return nil, fmt.Errorf("读取 Wan 2.7 驱动音频地址失败：%w", err)
		}
		media = append(media, dashScopeContractMedia{Type: "driving_audio", URL: value})
	}
	return media, nil
}

func buildWan27R2VMedia(input canvasGenerationInput, duration int) ([]dashScopeContractMedia, string, bool, error) {
	if metadataString(input.Metadata, "videoEndFrameNodeId") != "" {
		return nil, "", false, errors.New("Wan 2.7 参考生视频不支持尾帧")
	}
	visualCount := len(input.ReferenceImages) + len(input.ReferenceVideos)
	if visualCount < 1 || visualCount > 5 || len(input.ReferenceAudios) > 1 {
		return nil, "", false, errors.New("Wan 2.7 参考生视频需要 1-5 个图片/视频素材，并且最多 1 个参考音频")
	}
	if len(input.ReferenceVideos) > 0 && duration > 10 {
		return nil, "", false, errors.New("Wan 2.7 参考生视频包含视频素材时，输出时长只能为 2-10 秒")
	}
	startID := metadataString(input.Metadata, "videoStartFrameNodeId")
	foundStart := startID == ""
	media := make([]dashScopeContractMedia, 0, visualCount)
	for _, image := range input.ReferenceImages {
		mediaType := "reference_image"
		if startID != "" && image.ID == startID {
			mediaType = "first_frame"
			foundStart = true
		}
		if err := validateWan27Image(image, "Wan 2.7 参考图片"); err != nil {
			return nil, "", false, err
		}
		value, err := dashScopePublicMediaURL(image, false)
		if err != nil {
			return nil, "", false, err
		}
		media = append(media, dashScopeContractMedia{Type: mediaType, URL: value})
	}
	if !foundStart {
		return nil, "", false, errors.New("已配置的 Wan 2.7 首帧未包含在请求中")
	}
	for _, video := range input.ReferenceVideos {
		if err := validateProviderMediaMIME(video, []string{"video/mp4", "video/quicktime"}, "Wan 2.7 参考视频"); err != nil {
			return nil, "", false, err
		}
		if err := validateProviderMediaSize(video, 100*1024*1024, "Wan 2.7 参考视频"); err != nil {
			return nil, "", false, err
		}
		if err := validateProviderMediaDuration(video, 1, 30, "Wan 2.7 参考视频"); err != nil {
			return nil, "", false, err
		}
		if err := validateMediaDimensions(video, 240, 4096, 8, "Wan 2.7 参考视频"); err != nil {
			return nil, "", false, err
		}
		value, err := dashScopePublicMediaURL(video, false)
		if err != nil {
			return nil, "", false, err
		}
		media = append(media, dashScopeContractMedia{Type: "reference_video", URL: value})
	}
	referenceVoice := ""
	if len(input.ReferenceAudios) == 1 {
		audio := input.ReferenceAudios[0]
		if err := validateProviderMediaMIME(audio, []string{"audio/wav", "audio/x-wav", "audio/mpeg", "audio/mp3"}, "Wan 2.7 参考音频"); err != nil {
			return nil, "", false, err
		}
		if err := validateProviderMediaSize(audio, 15*1024*1024, "Wan 2.7 参考音频"); err != nil {
			return nil, "", false, err
		}
		if err := validateProviderMediaDuration(audio, 1, 10, "Wan 2.7 参考音频"); err != nil {
			return nil, "", false, err
		}
		value, err := dashScopePublicMediaURL(audio, false)
		if err != nil {
			return nil, "", false, err
		}
		referenceVoice = value
	}
	return media, referenceVoice, startID != "", nil
}

func orderedFrameMedia(images []providerMedia, startID string, endID string, label string) ([]providerMedia, error) {
	if startID != "" && startID == endID {
		return nil, fmt.Errorf("%s 首帧和尾帧不能使用同一个素材", label)
	}
	if startID == "" && endID != "" {
		return nil, fmt.Errorf("%s 配置尾帧时必须同时配置首帧", label)
	}
	if startID == "" {
		return append([]providerMedia(nil), images...), nil
	}
	byID := make(map[string]providerMedia, len(images))
	for _, image := range images {
		if image.ID == "" || byID[image.ID].ID != "" {
			return nil, fmt.Errorf("%s 帧素材必须具有唯一且非空的 ID", label)
		}
		byID[image.ID] = image
	}
	start, ok := byID[startID]
	if !ok {
		return nil, fmt.Errorf("已配置的%s首帧未包含在请求中", label)
	}
	result := []providerMedia{start}
	if endID != "" {
		end, ok := byID[endID]
		if !ok {
			return nil, fmt.Errorf("已配置的%s尾帧未包含在请求中", label)
		}
		result = append(result, end)
	}
	if len(result) != len(images) {
		return nil, fmt.Errorf("%s 帧模式不能混入普通参考图片", label)
	}
	return result, nil
}

func wan27FrameURL(image providerMedia, label string) (string, error) {
	if err := validateWan27Image(image, "Wan 2.7 "+label); err != nil {
		return "", err
	}
	value, err := dashScopePublicMediaURL(image, false)
	if err != nil {
		return "", fmt.Errorf("读取 Wan 2.7 %s地址失败：%w", label, err)
	}
	return value, nil
}

func validateWan27Image(image providerMedia, label string) error {
	if err := validateProviderMediaMIME(image, []string{"image/jpeg", "image/jpg", "image/png", "image/bmp", "image/webp"}, label); err != nil {
		return err
	}
	if err := validateProviderMediaSize(image, 20*1024*1024, label); err != nil {
		return err
	}
	return validateMediaDimensions(image, 240, 8000, 8, label)
}

func validateMediaDimensions(media providerMedia, minSide int, maxSide int, maxAspect float64, label string) error {
	if media.Width <= 0 || media.Height <= 0 {
		return nil
	}
	if media.Width < minSide || media.Height < minSide || media.Width > maxSide || media.Height > maxSide {
		return fmt.Errorf("%s尺寸必须在 %d-%d 像素之间", label, minSide, maxSide)
	}
	ratio := float64(max(media.Width, media.Height)) / float64(min(media.Width, media.Height))
	if ratio > maxAspect {
		return fmt.Errorf("%s宽高比不能超过 %.1f:1", label, maxAspect)
	}
	return nil
}

func wan27Duration(value string, maxSeconds int) (int, error) {
	if strings.TrimSpace(value) == "" {
		return 5, nil
	}
	duration, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || duration < 2 || duration > maxSeconds {
		return 0, fmt.Errorf("Wan 2.7 视频时长仅支持 2-%d 秒整数", maxSeconds)
	}
	return duration, nil
}

func wan27Resolution(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return "1080P", nil
	}
	if value == "720" || value == "1080" {
		value += "P"
	}
	if value != "720P" && value != "1080P" {
		return "", errors.New("Wan 2.7 输出分辨率仅支持 720P 或 1080P")
	}
	return value, nil
}

func wan27Ratio(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "16:9", nil
	}
	switch value {
	case "16:9", "9:16", "1:1", "4:3", "3:4":
		return value, nil
	default:
		return "", errors.New("Wan 2.7 画面比例不受支持")
	}
}
