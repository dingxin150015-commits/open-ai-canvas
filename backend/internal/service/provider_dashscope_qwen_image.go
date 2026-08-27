package service

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	qwenImageMinPixels = 512 * 512
	qwenImageMaxPixels = 2048 * 2048
	qwenImageMaxBytes  = 10 * 1024 * 1024
)

func isQwenImage30Model(value string) bool {
	modelName := strings.ToLower(strings.TrimSpace(strings.TrimPrefix(value, "models/")))
	return modelName == "qwen-image-3.0-pro" || modelName == "qwen-image-3.0"
}

func qwenImage30CapabilityConfig() *ImageCapabilityConfig {
	return &ImageCapabilityConfig{
		References: ImageReferenceConfig{
			PromptMaxChars: 32000,
			MaxImages:      3,
			MaxImageBytes:  qwenImageMaxBytes,
			MaskSupported:  false,
		},
		Size: ImageSizeConfig{
			Parameter: "size",
			Values: []string{
				"auto",
				"1:1", "2:3", "3:2", "3:4", "4:3", "9:16", "16:9",
				"1024x1024", "768x1152", "1024x1536", "1152x768", "1536x1024",
				"720x1280", "1080x1920", "1280x720", "1920x1080",
			},
			Default:     "auto",
			AllowCustom: true,
		},
		Quality:               ImageQualityConfig{Supported: false, Values: []string{}, Default: "auto"},
		TransparentBackground: VideoBooleanConfig{Supported: false, Default: false},
		ResponseFormat:        ParameterSupport{Supported: false},
		OutputFormat:          ParameterSupport{Supported: false},
		MaxOutputs:            6,
	}
}

func buildQwenImage30Request(input canvasGenerationInput) (map[string]interface{}, error) {
	if !isQwenImage30Model(input.Config.Model) {
		return nil, errors.New("当前模型不属于 Qwen Image 3.0 同步协议")
	}
	if input.Mask != nil {
		return nil, errors.New("Qwen Image 3.0 不支持蒙版参数")
	}
	if len(input.ReferenceImages) > 3 {
		return nil, errors.New("Qwen Image 3.0 最多支持 3 张参考图")
	}
	prompt := strings.TrimSpace(withSystemPrompt(input.Config, input.Prompt))
	if prompt == "" {
		return nil, errors.New("Qwen Image 3.0 必须提供提示词")
	}
	if quality := strings.ToLower(strings.TrimSpace(input.Config.Quality)); quality != "" && quality != "auto" {
		return nil, errors.New("Qwen Image 3.0 不支持图片质量参数")
	}
	if parseBool(input.Config.TransparentBackground, false) {
		return nil, errors.New("Qwen Image 3.0 不支持透明背景参数")
	}

	content := make([]map[string]interface{}, 0, len(input.ReferenceImages)+1)
	for index, reference := range input.ReferenceImages {
		raw, mimeType, err := mediaBytes(reference)
		if err != nil {
			return nil, fmt.Errorf("读取第 %d 张 Qwen 参考图失败：%w", index+1, err)
		}
		if len(raw) > qwenImageMaxBytes {
			return nil, fmt.Errorf("第 %d 张 Qwen 参考图超过 10 MB", index+1)
		}
		mimeType = normalizedMediaMimeType(mimeType, raw)
		if !qwenImageMimeAllowed(mimeType) {
			return nil, fmt.Errorf("第 %d 张 Qwen 参考图格式不受支持", index+1)
		}
		content = append(content, map[string]interface{}{"image": dataURL(mimeType, raw)})
	}
	content = append(content, map[string]interface{}{"text": prompt})

	count := 1
	if value := strings.TrimSpace(input.Config.Count); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 6 {
			return nil, errors.New("Qwen Image 3.0 输出数量必须为 1-6")
		}
		count = parsed
	}
	promptExtend, err := dashScopeMetadataBool(input.Metadata, "prompt_extend", true)
	if err != nil {
		return nil, err
	}
	promptExtendMode := strings.ToLower(strings.TrimSpace(metadataString(input.Metadata, "prompt_extend_mode")))
	if promptExtendMode == "" {
		promptExtendMode = "direct"
	}
	if promptExtendMode != "direct" && promptExtendMode != "agent" {
		return nil, errors.New("Qwen Image 3.0 提示词改写模式必须为 direct 或 agent")
	}
	if !promptExtend && input.Metadata != nil && input.Metadata["prompt_extend_mode"] != nil {
		return nil, errors.New("关闭提示词改写时不能设置改写模式")
	}
	if len(input.ReferenceImages) > 0 && promptExtendMode == "agent" {
		return nil, errors.New("Qwen Image 3.0 图片编辑不支持 agent 提示词改写")
	}
	enableThinking, err := dashScopeMetadataBool(input.Metadata, "enable_thinking", promptExtend)
	if err != nil {
		return nil, err
	}
	if !promptExtend && enableThinking {
		return nil, errors.New("Qwen Image 3.0 关闭提示词改写时不能开启思考模式")
	}
	watermark, err := dashScopeMetadataBool(input.Metadata, "watermark", false)
	if err != nil {
		return nil, err
	}
	negativePrompt := metadataString(input.Metadata, "negative_prompt")
	if utf8.RuneCountInString(negativePrompt) > 500 {
		return nil, errors.New("Qwen Image 3.0 负面提示词不能超过 500 个字符")
	}
	seed, err := qwenImageSeed(input.Metadata)
	if err != nil {
		return nil, err
	}
	size, err := qwenImage30Size(input.Config.Size)
	if err != nil {
		return nil, err
	}

	parameters := map[string]interface{}{
		"n":               count,
		"prompt_extend":   promptExtend,
		"watermark":       watermark,
		"enable_thinking": enableThinking,
	}
	if promptExtend {
		parameters["prompt_extend_mode"] = promptExtendMode
	}
	if negativePrompt != "" {
		parameters["negative_prompt"] = negativePrompt
	}
	if seed != nil {
		parameters["seed"] = *seed
	}
	if size != "" {
		parameters["size"] = size
	}

	return map[string]interface{}{
		"model": input.Config.Model,
		"input": map[string]interface{}{
			"messages": []map[string]interface{}{{
				"role":    "user",
				"content": content,
			}},
		},
		"parameters": parameters,
	}, nil
}

func qwenImage30Size(value string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(value, "×", "x")))
	if normalized == "" || normalized == "auto" {
		return "", nil
	}
	if mapped := map[string]string{
		"1:1": "1024x1024", "2:3": "1024x1536", "3:2": "1536x1024",
		"3:4": "768x1024", "4:3": "1024x768", "9:16": "720x1280", "16:9": "1280x720",
	}[normalized]; mapped != "" {
		normalized = mapped
	}
	if strings.Contains(normalized, ":") {
		parts := strings.Split(normalized, ":")
		if len(parts) != 2 {
			return "", errors.New("Qwen 图片比例格式无效")
		}
		widthRatio, widthErr := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		heightRatio, heightErr := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if widthErr != nil || heightErr != nil || widthRatio <= 0 || heightRatio <= 0 {
			return "", errors.New("Qwen 图片比例必须是正数")
		}
		ratio := widthRatio / heightRatio
		if ratio < 1.0/8.0 || ratio > 8 {
			return "", errors.New("Qwen 图片宽高比必须在 1:8 到 8:1 之间")
		}
		width := math.Sqrt(1024 * 1024 * ratio)
		height := math.Sqrt(1024 * 1024 / ratio)
		if width > 2048 {
			height *= 2048 / width
			width = 2048
		}
		if height > 2048 {
			width *= 2048 / height
			height = 2048
		}
		return validateQwenImage30Pixels(int(math.Round(width)), int(math.Round(height)))
	}
	normalized = strings.ReplaceAll(normalized, "*", "x")
	parts := strings.Split(normalized, "x")
	if len(parts) != 2 {
		return "", errors.New("Qwen 图片尺寸请使用宽x高或宽:高")
	}
	width, widthErr := strconv.Atoi(strings.TrimSpace(parts[0]))
	height, heightErr := strconv.Atoi(strings.TrimSpace(parts[1]))
	if widthErr != nil || heightErr != nil {
		return "", errors.New("Qwen 图片尺寸必须是整数")
	}
	return validateQwenImage30Pixels(width, height)
}

func validateQwenImage30Pixels(width int, height int) (string, error) {
	if width <= 0 || height <= 0 {
		return "", errors.New("Qwen 图片宽高必须是正整数")
	}
	pixels := int64(width) * int64(height)
	if pixels < qwenImageMinPixels || pixels > qwenImageMaxPixels {
		return "", errors.New("Qwen 图片总像素必须在 512x512 到 2048x2048 之间")
	}
	ratio := float64(width) / float64(height)
	if ratio < 1.0/8.0 || ratio > 8 {
		return "", errors.New("Qwen 图片宽高比必须在 1:8 到 8:1 之间")
	}
	return fmt.Sprintf("%d*%d", width, height), nil
}

func qwenImageMimeAllowed(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "image/jpeg", "image/png", "image/bmp", "image/x-ms-bmp", "image/tiff", "image/webp", "image/gif":
		return true
	default:
		return false
	}
}

func qwenImageSeed(metadata map[string]interface{}) (*int64, error) {
	if metadata == nil || metadata["seed"] == nil {
		return nil, nil
	}
	var value int64
	switch raw := metadata["seed"].(type) {
	case int:
		value = int64(raw)
	case int64:
		value = raw
	case float64:
		if math.Trunc(raw) != raw {
			return nil, errors.New("Qwen 图片随机种子必须是整数")
		}
		value = int64(raw)
	case string:
		parsed, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil {
			return nil, errors.New("Qwen 图片随机种子必须是整数")
		}
		value = parsed
	default:
		return nil, errors.New("Qwen 图片随机种子必须是整数")
	}
	if value < 0 || value > 2147483647 {
		return nil, errors.New("Qwen 图片随机种子必须在 0-2147483647 之间")
	}
	return &value, nil
}
