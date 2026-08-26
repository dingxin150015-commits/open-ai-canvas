package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type happyHorseVideoInput struct {
	Prompt string                   `json:"prompt,omitempty"`
	Media  []dashScopeContractMedia `json:"media,omitempty"`
}

type happyHorseVideoParameters struct {
	Resolution string `json:"resolution"`
	Ratio      string `json:"ratio,omitempty"`
	Duration   int    `json:"duration"`
	Watermark  bool   `json:"watermark"`
	Seed       *int64 `json:"seed,omitempty"`
}

type happyHorseVideoRequest struct {
	Model      string                    `json:"model"`
	Input      happyHorseVideoInput      `json:"input"`
	Parameters happyHorseVideoParameters `json:"parameters"`
}

func happyHorse11VideoKind(value string) string {
	value = strings.ToLower(strings.TrimPrefix(strings.TrimSpace(value), "models/"))
	switch value {
	case "happyhorse-1.1-t2v":
		return "t2v"
	case "happyhorse-1.1-i2v":
		return "i2v"
	case "happyhorse-1.1-r2v":
		return "r2v"
	default:
		return ""
	}
}

func isHappyHorse11ReadyVideoModel(value string) bool { return happyHorse11VideoKind(value) != "" }

func buildHappyHorseVideoRequest(input canvasGenerationInput) (map[string]interface{}, error) {
	modelName := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(input.Config.Model), "models/"))
	kind := happyHorse11VideoKind(modelName)
	if kind == "" {
		return nil, fmt.Errorf("HappyHorse 1.1 Adapter 不支持模型 %s", modelName)
	}
	prompt := strings.TrimSpace(input.Prompt)
	if err := validateHappyHorsePrompt(prompt, kind != "i2v"); err != nil {
		return nil, err
	}
	duration, err := happyHorseDuration(input.Config.VideoSeconds)
	if err != nil {
		return nil, err
	}
	resolution, err := happyHorseResolution(input.Config.VQuality)
	if err != nil {
		return nil, err
	}
	seed, err := dashScopeVideoSeed(input.Metadata)
	if err != nil {
		return nil, err
	}
	request := happyHorseVideoRequest{
		Model: modelName,
		Input: happyHorseVideoInput{Prompt: prompt},
		Parameters: happyHorseVideoParameters{
			Resolution: resolution,
			Duration:   duration,
			Watermark:  parseBool(input.Config.VideoWatermark, true),
			Seed:       seed,
		},
	}
	switch kind {
	case "t2v":
		if err := validateDashScopeVideoOperation(input, "HappyHorse 文生视频", "text_to_video"); err != nil {
			return nil, err
		}
		if len(input.ReferenceImages)+len(input.ReferenceVideos)+len(input.ReferenceAudios) > 0 {
			return nil, errors.New("HappyHorse 文生视频不支持参考素材")
		}
		ratio, err := happyHorseRatio(input.Config.Size)
		if err != nil {
			return nil, err
		}
		request.Parameters.Ratio = ratio
	case "i2v":
		if err := validateDashScopeVideoOperation(input, "HappyHorse 图生视频", "image_to_video"); err != nil {
			return nil, err
		}
		media, err := buildHappyHorseI2VMedia(input)
		if err != nil {
			return nil, err
		}
		request.Input.Media = media
	case "r2v":
		if err := validateDashScopeVideoOperation(input, "HappyHorse 参考生视频", "reference_to_video", "image_to_video"); err != nil {
			return nil, err
		}
		media, err := buildHappyHorseR2VMedia(input)
		if err != nil {
			return nil, err
		}
		request.Input.Media = media
		ratio, err := happyHorseRatio(input.Config.Size)
		if err != nil {
			return nil, err
		}
		request.Parameters.Ratio = ratio
	}
	return requestAsMap(request)
}

func buildHappyHorseI2VMedia(input canvasGenerationInput) ([]dashScopeContractMedia, error) {
	if len(input.ReferenceImages) != 1 || len(input.ReferenceVideos) > 0 || len(input.ReferenceAudios) > 0 {
		return nil, errors.New("HappyHorse 图生视频有且仅支持 1 张首帧图片")
	}
	if metadataString(input.Metadata, "videoEndFrameNodeId") != "" {
		return nil, errors.New("HappyHorse 图生视频不支持尾帧")
	}
	image := input.ReferenceImages[0]
	if err := validateProviderMediaMIME(image, []string{"image/jpeg", "image/jpg", "image/png", "image/webp"}, "HappyHorse 首帧"); err != nil {
		return nil, err
	}
	if startID := metadataString(input.Metadata, "videoStartFrameNodeId"); startID != "" && image.ID != startID {
		return nil, errors.New("已配置的 HappyHorse 首帧未包含在请求中")
	}
	if err := validateProviderMediaSize(image, 20*1024*1024, "HappyHorse 首帧"); err != nil {
		return nil, err
	}
	if image.Width > 0 && image.Height > 0 {
		if image.Width < 300 || image.Height < 300 {
			return nil, errors.New("HappyHorse 首帧宽和高均不能小于 300 像素")
		}
		ratio := float64(max(image.Width, image.Height)) / float64(min(image.Width, image.Height))
		if ratio > 2.5 {
			return nil, errors.New("HappyHorse 首帧宽高比必须在 1:2.5 到 2.5:1 之间")
		}
	}
	value, err := dashScopePublicMediaURL(image, true)
	if err != nil {
		return nil, fmt.Errorf("读取 HappyHorse 首帧地址失败：%w", err)
	}
	return []dashScopeContractMedia{{Type: "first_frame", URL: value}}, nil
}

func buildHappyHorseR2VMedia(input canvasGenerationInput) ([]dashScopeContractMedia, error) {
	if len(input.ReferenceImages) < 1 || len(input.ReferenceImages) > 9 || len(input.ReferenceVideos) > 0 || len(input.ReferenceAudios) > 0 {
		return nil, errors.New("HappyHorse 参考生视频仅支持 1-9 张参考图片")
	}
	if metadataString(input.Metadata, "videoStartFrameNodeId") != "" || metadataString(input.Metadata, "videoEndFrameNodeId") != "" {
		return nil, errors.New("HappyHorse 参考生视频不支持首帧或尾帧角色")
	}
	media := make([]dashScopeContractMedia, 0, len(input.ReferenceImages))
	for index, image := range input.ReferenceImages {
		if err := validateProviderMediaMIME(image, []string{"image/jpeg", "image/jpg", "image/png", "image/webp"}, "HappyHorse 参考图片"); err != nil {
			return nil, err
		}
		if err := validateProviderMediaSize(image, 20*1024*1024, "HappyHorse 参考图片"); err != nil {
			return nil, err
		}
		if image.Width > 0 && image.Height > 0 && min(image.Width, image.Height) < 400 {
			return nil, fmt.Errorf("HappyHorse 第 %d 张参考图片短边不能小于 400 像素", index+1)
		}
		value, err := dashScopePublicMediaURL(image, true)
		if err != nil {
			return nil, fmt.Errorf("读取 HappyHorse 第 %d 张参考图片地址失败：%w", index+1, err)
		}
		media = append(media, dashScopeContractMedia{Type: "reference_image", URL: value})
	}
	return media, nil
}

func happyHorseDuration(value string) (int, error) {
	if strings.TrimSpace(value) == "" {
		return 5, nil
	}
	duration, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || duration < 3 || duration > 15 {
		return 0, errors.New("HappyHorse 视频时长仅支持 3-15 秒整数")
	}
	return duration, nil
}

func happyHorseResolution(value string) (string, error) {
	value = strings.ToUpper(strings.TrimSpace(value))
	if value == "" {
		return "1080P", nil
	}
	if value == "480" || value == "720" || value == "1080" {
		value += "P"
	}
	if value != "480P" && value != "720P" && value != "1080P" {
		return "", errors.New("HappyHorse 输出分辨率仅支持 480P、720P 或 1080P")
	}
	return value, nil
}

func happyHorseRatio(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "16:9", nil
	}
	switch value {
	case "16:9", "9:16", "1:1", "4:3", "3:4", "4:5", "5:4", "9:21", "21:9":
		return value, nil
	default:
		return "", errors.New("HappyHorse 画面比例不受支持")
	}
}
