package service

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Product capability profiles remain separate from wire serialization. All
// executable request and response contracts live in the declarative plugins.
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

const wan30VideoModel = "wan3.0-video"

// Resolve the actual routed model before allowing a media-only request. Normal
// capability admission still validates all reference counts and permissions.
func bailianTaskAllowsEmptyPrompt(input map[string]any) bool {
	if input["mode"] != "video" {
		return false
	}
	config, _ := input["config"].(map[string]any)
	name, _ := config["model"].(string)
	providerName, _ := config["providerModelKey"].(string)
	resolved := firstNonEmpty(providerName, name)
	images, _ := input["referenceImages"].([]any)
	videos, _ := input["referenceVideos"].([]any)
	if happyHorse11VideoKind(resolved) == "i2v" {
		return len(images) == 1 && len(videos) == 0
	}
	return wan27VideoKind(resolved) == "i2v" && len(images)+len(videos) > 0
}

func bailianOptionalPromptModel(name string) bool {
	return wan27VideoKind(name) == "i2v" || happyHorse11VideoKind(name) == "i2v"
}

func isWan30VideoModel(value string) bool {
	return strings.EqualFold(strings.TrimPrefix(strings.TrimSpace(value), "models/"), wan30VideoModel)
}
