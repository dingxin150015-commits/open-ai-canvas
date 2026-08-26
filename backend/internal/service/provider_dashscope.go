package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// runDashScopeImageTask 实现百炼（DashScope）图片生成协议（同步模式）
//
// API 文档：
// - 千问图像编辑：https://platform.qianwenai.com/docs/api-reference/image-generation/qwen-image-editing
// - 万相图像编辑：https://platform.qianwenai.com/docs/developer-guides/image-generation/wan-image-editing
// - 模型选择指南：https://platform.qianwenai.com/docs/developer-guides/getting-started/image-models
//
// 支持的模型和能力（通过管理后台"最大参考图"配置）：
// - qwen-image-3.0-pro / qwen-image-3.0: 1-3 张输入图像
// - qwen-image-2.0-pro / qwen-image-2.0: 1-3 张输入图像
// - wan2.7-image-pro / wan2.7-image: 0-9 张输入图像（0张=文生图模式）
// - wan2.6-image: 1-4 张输入图像
//
// 参数限制由配置系统管理，在 validateImageTask 中统一验证。
// 此函数只负责构建符合 DashScope 协议的请求格式。
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

	// 2. 处理所有参考图（添加到 content 数组）
	// 注意：参考图数量已在 validateImageTask 中验证，这里直接处理
	for i, refImg := range input.ReferenceImages {
		raw, mimeType, err := mediaBytes(refImg)
		if err != nil {
			return nil, fmt.Errorf("读取第 %d 张参考图失败：%w", i+1, err)
		}
		// 使用实际的 MIME 类型（而不是硬编码为 image/png）
		mimeType = normalizedMediaMimeType(mimeType, raw)
		// 每张参考图作为独立 content 对象（使用 data URL 格式）
		content = append(content, map[string]interface{}{
			"image": dataURL(mimeType, raw),
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
// 支持的尺寸：720*1280, 1280*720, 1024*1024, 2048*2048, 4096*4096 等
func normalizeDashScopeImageSize(value string) string {
	size := strings.TrimSpace(value)
	if size == "" {
		return "" // 返回空，让 API 使用服务端默认值
	}

	// 规范化格式：x 或 X -> *
	size = strings.ReplaceAll(size, "x", "*")
	size = strings.ReplaceAll(size, "X", "*")

	// 常见比例别名映射
	aliasMap := map[string]string{
		"1:1":  "1024*1024",
		"9:16": "720*1280",
		"16:9": "1280*720",
		"2:3":  "1024*1536",
		"3:2":  "1536*1024",
		"3:4":  "768*1024",
		"4:3":  "1024*768",
	}

	if mapped, ok := aliasMap[size]; ok {
		return mapped
	}

	// 已经是具体尺寸格式（如 "2048*2048"），直接返回
	// 让配置系统和 API 验证是否合法
	return size
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
	endpoint, err := dashScopeNativeEndpoint(config.BaseURL, path)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
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
	endpoint, err := dashScopeNativeEndpoint(config.BaseURL, path)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
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
	endpoint, err := dashScopeNativeEndpoint(config.BaseURL, path)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	ApplyOutboundHeaders(req, config.Headers)
	return doJSON(req, target)
}

// dashScopeNativeEndpoint keeps the configured account/region host while moving
// from an OpenAI-compatible base path to the native DashScope API surface.
// Token Plan and workspace-specific MaaS keys are bound to their configured host,
// so replacing the host with the public DashScope domain would break entitlement.
func dashScopeNativeEndpoint(rawBaseURL string, path string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawBaseURL))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", errors.New("DashScope Base URL 无效")
	}
	if parsed.RawQuery != "" || parsed.Fragment != "" {
		return "", errors.New("DashScope Base URL 不能包含查询参数或片段")
	}

	basePath := strings.TrimRight(parsed.Path, "/")
	for _, suffix := range []string{"/compatible-mode/v1", "/api/v1"} {
		if strings.HasSuffix(strings.ToLower(basePath), suffix) {
			basePath = strings.TrimRight(basePath[:len(basePath)-len(suffix)], "/")
			break
		}
	}
	parsed.Path = basePath
	parsed.RawPath = ""
	return strings.TrimRight(parsed.String(), "/") + "/" + strings.TrimLeft(path, "/"), nil
}
