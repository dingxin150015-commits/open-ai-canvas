package service

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func dashScopeVideoSeed(metadata map[string]interface{}) (*int64, error) {
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
			return nil, errors.New("DashScope 视频随机种子必须是整数")
		}
		value = int64(raw)
	case string:
		parsed, err := strconv.ParseInt(strings.TrimSpace(raw), 10, 64)
		if err != nil {
			return nil, errors.New("DashScope 视频随机种子必须是整数")
		}
		value = parsed
	default:
		return nil, errors.New("DashScope 视频随机种子必须是整数")
	}
	if value < 0 || value > 2147483647 {
		return nil, errors.New("DashScope 视频随机种子必须在 0-2147483647 之间")
	}
	return &value, nil
}

func dashScopeMetadataBool(metadata map[string]interface{}, key string, fallback bool) (bool, error) {
	if metadata == nil || metadata[key] == nil {
		return fallback, nil
	}
	switch value := metadata[key].(type) {
	case bool:
		return value, nil
	case string:
		normalized := strings.ToLower(strings.TrimSpace(value))
		if normalized == "true" || normalized == "1" {
			return true, nil
		}
		if normalized == "false" || normalized == "0" {
			return false, nil
		}
	}
	return false, fmt.Errorf("%s 必须是布尔值", key)
}

func dashScopePublicMediaURL(media providerMedia, allowDataURL bool) (string, error) {
	value := strings.TrimSpace(media.URL)
	if value == "" || strings.HasPrefix(value, "asset://") {
		value = strings.TrimSpace(media.DataURL)
	}
	if strings.HasPrefix(value, "data:") {
		if !allowDataURL {
			return "", errors.New("该模型的参考素材必须使用公网 HTTP(S) URL")
		}
		if !strings.Contains(value, ";base64,") {
			return "", errors.New("data URL 必须使用 base64 编码")
		}
		return value, nil
	}
	parsed, err := url.Parse(value)
	if err != nil || parsed.Scheme == "" || !isPublicMediaURL(value) {
		if allowDataURL {
			return "", errors.New("参考素材需要公网 HTTP(S) URL 或 base64 data URL")
		}
		return "", errors.New("参考素材需要公网 HTTP(S) URL")
	}
	if _, err := ValidateOutboundURL(value); err != nil {
		return "", err
	}
	return value, nil
}

func validatePromptRunes(prompt string, maxRunes int, modelLabel string, required bool) error {
	prompt = strings.TrimSpace(prompt)
	if required && prompt == "" {
		return fmt.Errorf("%s 必须提供提示词", modelLabel)
	}
	if utf8.RuneCountInString(prompt) > maxRunes {
		return fmt.Errorf("%s 提示词不能超过 %d 个字符", modelLabel, maxRunes)
	}
	return nil
}

func validateNegativePrompt(prompt string, maxRunes int, modelLabel string) error {
	if utf8.RuneCountInString(strings.TrimSpace(prompt)) > maxRunes {
		return fmt.Errorf("%s 负向提示词不能超过 %d 个字符", modelLabel, maxRunes)
	}
	return nil
}

func validateHappyHorsePrompt(prompt string, required bool) error {
	prompt = strings.TrimSpace(prompt)
	if required && prompt == "" {
		return errors.New("HappyHorse 必须提供提示词")
	}
	weight := 0
	for _, character := range prompt {
		if unicode.Is(unicode.Han, character) {
			weight += 2
		} else {
			weight++
		}
	}
	if weight > 5000 {
		return errors.New("HappyHorse 提示词不能超过 2500 个中文字符或 5000 个非中文字符")
	}
	return nil
}

func validateDashScopeVideoOperation(input canvasGenerationInput, modelLabel string, allowed ...string) error {
	operation := strings.ToLower(metadataString(input.Metadata, "videoEditOperation"))
	if operation == "" {
		return nil
	}
	for _, candidate := range allowed {
		if operation == candidate {
			return nil
		}
	}
	return fmt.Errorf("%s 不支持生成模式 %s", modelLabel, operation)
}

func validateProviderMediaSize(media providerMedia, maxBytes int64, label string) error {
	if maxBytes > 0 && media.Bytes > maxBytes {
		return fmt.Errorf("%s文件大小超过限制", label)
	}
	return nil
}

func validateProviderMediaMIME(media providerMedia, allowed []string, label string) error {
	value := strings.ToLower(strings.TrimSpace(strings.Split(media.MimeType, ";")[0]))
	if value == "" {
		return nil
	}
	for _, candidate := range allowed {
		if value == candidate {
			return nil
		}
	}
	return fmt.Errorf("%s格式 %s 不受支持", label, value)
}

func validateProviderMediaDuration(media providerMedia, minSeconds int, maxSeconds int, label string) error {
	if media.DurationMs <= 0 {
		return nil
	}
	minMillis := int64(minSeconds) * 1000
	maxMillis := int64(maxSeconds) * 1000
	if media.DurationMs < minMillis || media.DurationMs > maxMillis {
		return fmt.Errorf("%s时长必须在 %d-%d 秒之间", label, minSeconds, maxSeconds)
	}
	return nil
}
