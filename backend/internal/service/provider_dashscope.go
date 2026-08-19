package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// runDashScopeImageTask 实现百炼（DashScope）图片生成协议（同步模式）
// API 文档：https://help.aliyun.com/zh/model-studio/developer-reference/text-to-image-api
func runDashScopeImageTask(ctx context.Context, input canvasGenerationInput) (map[string]interface{}, error) {
	if input.Mask != nil {
		return nil, errors.New("DashScope 图片协议不支持蒙版编辑，请移除蒙版后重试")
	}

	// 构建多模态 content 数组（DashScope 图片生成使用多模态消息格式）
	content := []map[string]interface{}{}

	// 1. 添加文本内容
	promptText := strings.TrimSpace(input.Prompt)
	if systemPrompt := strings.TrimSpace(input.Config.SystemPrompt); systemPrompt != "" {
		// System prompt 和 user prompt 合并
		promptText = systemPrompt + "\n\n" + promptText
	}
	if promptText != "" {
		content = append(content, map[string]interface{}{
			"text": promptText,
		})
	}

	// 2. 处理参考图（添加到 content 数组）
	if len(input.ReferenceImages) > 0 {
		if len(input.ReferenceImages) > 1 {
			return nil, errors.New("DashScope 图片协议当前只支持 1 张参考图")
		}
		raw, _, err := mediaBytes(input.ReferenceImages[0])
		if err != nil {
			return nil, fmt.Errorf("读取 DashScope 参考图失败：%w", err)
		}
		// 参考图作为独立 content 对象（使用 data URL 格式）
		content = append(content, map[string]interface{}{
			"image": "data:image/png;base64," + base64.StdEncoding.EncodeToString(raw),
		})
	}

	// 3. 构建请求体（多模态消息格式）
	body := map[string]interface{}{
		"model": input.Config.Model,
		"input": map[string]interface{}{
			"messages": []map[string]interface{}{
				{
					"role":    "user",
					"content": content, // content 必须是对象数组
				},
			},
		},
		"parameters": map[string]interface{}{},
	}

	// 4. 设置图片尺寸
	if size := normalizeDashScopeImageSize(input.Config.Size); size != "" {
		body["parameters"].(map[string]interface{})["size"] = size
	}

	// 5. 同步调用（不使用异步模式）
	var response dashScopeResponse
	if err := postDashScopeJSONSync(ctx, input.Config, "/api/v1/services/aigc/multimodal-generation/generation", body, &response); err != nil {
		return nil, err
	}

	// 6. 检查响应格式并提取图片 URL
	// 同步模式返回 choices 格式
	if len(response.Output.Choices) > 0 && len(response.Output.Choices[0].Message.Content) > 0 {
		images, err := dashScopeImageDataURLs(ctx, input.Config, response)
		if err != nil {
			return nil, fmt.Errorf("DashScope 图片结果处理失败：%w", err)
		}
		return map[string]interface{}{"mode": "image", "images": images}, nil
	}

	// 回退：检查异步模式的响应格式（保留兼容性）
	status := strings.ToLower(strings.TrimSpace(response.Output.TaskStatus))
	if status == "failed" || (status == "" && len(response.Output.Results) == 0) {
		errorMsg := "DashScope 图片生成失败"
		if response.Output.Message != "" {
			errorMsg = fmt.Sprintf("DashScope 图片生成失败：%s", response.Output.Message)
		}
		return nil, errors.New(errorMsg)
	}

	// 异步模式的图片处理（如果上面的 choices 格式失败）
	images, err := dashScopeImageDataURLs(ctx, input.Config, response)
	if err != nil {
		return nil, fmt.Errorf("DashScope 图片结果处理失败：%w", err)
	}

	return map[string]interface{}{"mode": "image", "images": images}, nil
}


// dashScopeImageDataURLs 下载 DashScope 图片 URL 并转为 data URL
func dashScopeImageDataURLs(ctx context.Context, config providerConfig, response dashScopeResponse) ([]map[string]string, error) {
	images := make([]map[string]string, 0)

	// 优先处理同步模式的 choices 格式
	if len(response.Output.Choices) > 0 {
		for _, choice := range response.Output.Choices {
			for _, item := range choice.Message.Content {
				imageURL := strings.TrimSpace(item.Image)
				if imageURL == "" {
					continue
				}

				// 下载图片并转为 data URL（DashScope 返回的 URL 有效期仅 24 小时）
				data, mimeType, err := getExternalBinary(withProviderRequestKind(ctx, "download"), imageURL)
				if err != nil {
					return nil, fmt.Errorf("DashScope 图片下载失败：%w", err)
				}
				mimeType = normalizedMediaMimeType(mimeType, data)
				images = append(images, map[string]string{"dataUrl": dataURL(mimeType, data)})
			}
		}

		if len(images) > 0 {
			return images, nil
		}
	}

	// 回退：处理异步模式的 results 格式
	results := response.Output.Results
	if len(results) == 0 {
		return nil, errors.New("DashScope 接口没有返回图片")
	}

	for _, result := range results {
		imageURL := strings.TrimSpace(result.URL)
		if imageURL == "" {
			continue
		}

		// 下载图片并转为 data URL（DashScope 返回的 URL 有效期仅 24 小时）
		data, mimeType, err := getExternalBinary(withProviderRequestKind(ctx, "download"), imageURL)
		if err != nil {
			return nil, fmt.Errorf("DashScope 图片下载失败：%w", err)
		}
		mimeType = normalizedMediaMimeType(mimeType, data)
		images = append(images, map[string]string{"dataUrl": dataURL(mimeType, data)})
	}

	if len(images) == 0 {
		return nil, errors.New("DashScope 接口没有返回可用图片")
	}

	return images, nil
}

// normalizeDashScopeImageSize 规范化 DashScope 图片尺寸
// 支持的尺寸：720*1280, 1280*720, 1024*1024 等
func normalizeDashScopeImageSize(value string) string {
	size := strings.TrimSpace(value)
	// 将 1024x1024 转为 1024*1024
	size = strings.ReplaceAll(size, "x", "*")

	// 常见尺寸映射
	switch size {
	case "1024*1024", "1:1":
		return "1024*1024"
	case "720*1280", "9:16":
		return "720*1280"
	case "1280*720", "16:9":
		return "1280*720"
	default:
		return "1024*1024" // 默认正方形
	}
}

// DashScope API 响应结构（同步模式使用 choices 格式）
type dashScopeResponse struct {
	Output struct {
		Choices []struct {
			FinishReason string `json:"finish_reason"`
			Message      struct {
				Role    string `json:"role"`
				Content []struct {
					Image string `json:"image"`
				} `json:"content"`
			} `json:"message"`
		} `json:"choices"`

		// 异步模式字段（保留以防需要回退到异步）
		TaskID     string `json:"task_id,omitempty"`
		TaskStatus string `json:"task_status,omitempty"`
		Message    string `json:"message,omitempty"`
		Results    []struct {
			URL string `json:"url"`
		} `json:"results,omitempty"`
	} `json:"output"`
	RequestID string `json:"request_id"`
}

// postDashScopeJSONSync 向 DashScope 发送同步 POST 请求（不带 X-DashScope-Async header）
func postDashScopeJSONSync(ctx context.Context, config providerConfig, path string, body interface{}, target interface{}) error {
	data, _ := json.Marshal(body)
	url := strings.TrimRight(config.BaseURL, "/") + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	// 同步模式：不设置 X-DashScope-Async header
	ApplyOutboundHeaders(req, config.Headers)
	return doJSON(req, target)
}

// postDashScopeJSON 向 DashScope 发送 POST 请求（带 X-DashScope-Async header）
func postDashScopeJSON(ctx context.Context, config providerConfig, path string, body interface{}, target interface{}) error {
	data, _ := json.Marshal(body)
	// DashScope API 路径不需要 /v1 前缀
	url := strings.TrimRight(config.BaseURL, "/") + path
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-Async", "enable")
	ApplyOutboundHeaders(req, config.Headers)
	return doJSON(req, target)
}

// getDashScopeJSON 向 DashScope 发送 GET 请求
func getDashScopeJSON(ctx context.Context, config providerConfig, path string, target interface{}) error {
	url := strings.TrimRight(config.BaseURL, "/") + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	ApplyOutboundHeaders(req, config.Headers)
	return doJSON(req, target)
}
