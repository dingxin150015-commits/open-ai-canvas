package app

// Frozen pre-migration oracle; compiled only into tests, never the server.

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"
)

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
