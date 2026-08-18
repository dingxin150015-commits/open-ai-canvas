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
	"time"
)

// runDashScopeImageTask 实现百炼（DashScope）图片生成协议
// API 文档：https://help.aliyun.com/zh/model-studio/developer-reference/text-to-image-api
func runDashScopeImageTask(ctx context.Context, input canvasGenerationInput) (map[string]interface{}, error) {
	if input.Mask != nil {
		return nil, errors.New("DashScope 图片协议不支持蒙版编辑，请移除蒙版后重试")
	}

	// 构建请求体
	body := map[string]interface{}{
		"model": input.Config.Model,
		"input": map[string]interface{}{
			"prompt": withSystemPrompt(input.Config, input.Prompt),
		},
		"parameters": map[string]interface{}{},
	}

	// 设置图片尺寸
	if size := normalizeDashScopeImageSize(input.Config.Size); size != "" {
		body["parameters"].(map[string]interface{})["size"] = size
	}

	// 处理参考图（DashScope 支持参考图）
	if len(input.ReferenceImages) > 0 {
		if len(input.ReferenceImages) > 1 {
			return nil, errors.New("DashScope 图片协议当前只支持 1 张参考图")
		}
		// 将参考图转为 base64
		raw, _, err := mediaBytes(input.ReferenceImages[0])
		if err != nil {
			return nil, fmt.Errorf("读取 DashScope 参考图失败：%w", err)
		}
		body["input"].(map[string]interface{})["ref_img"] = base64.StdEncoding.EncodeToString(raw)
	}

	// 提交异步任务
	taskID, err := submitDashScopeTask(ctx, input.Config, body)
	if err != nil {
		return nil, err
	}

	// 轮询任务状态
	for deadline := providerPollingDeadline(ctx); time.Now().Before(deadline); {
		result, err := pollDashScopeTask(ctx, input.Config, taskID)
		if err != nil {
			return nil, err
		}

		status := strings.ToLower(strings.TrimSpace(result.Output.TaskStatus))
		switch status {
		case "succeeded":
			// 下载图片 URL 并转为 data URL（避免 24h 过期）
			images, err := dashScopeImageDataURLs(ctx, input.Config, result)
			if err != nil {
				return nil, fmt.Errorf("DashScope 图片任务 %s 结果处理失败：%w", taskID, err)
			}
			return map[string]interface{}{"mode": "image", "images": images}, nil
		case "failed":
			errorMsg := "DashScope 图片生成失败"
			if result.Output.Message != "" {
				errorMsg = fmt.Sprintf("DashScope 图片生成失败：%s", result.Output.Message)
			}
			return nil, errors.New(errorMsg)
		case "pending", "running":
			// 继续轮询
		default:
			return nil, fmt.Errorf("DashScope 返回未知状态：%s", status)
		}

		if err := sleepContext(ctx, 3*time.Second); err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("DashScope 图片生成超时（任务 %s）", taskID)
}

// submitDashScopeTask 提交 DashScope 异步任务
func submitDashScopeTask(ctx context.Context, config providerConfig, body map[string]interface{}) (string, error) {
	if resumed := resumedProviderRequestID(ctx); resumed != "" {
		return resumed, nil
	}

	var response dashScopeResponse
	if err := postDashScopeJSON(withProviderRequestKind(ctx, "create"), config, "/api/v1/services/aigc/multimodal-generation/generation", body, &response); err != nil {
		return "", err
	}

	taskID := strings.TrimSpace(response.Output.TaskID)
	if taskID == "" {
		return "", errors.New("DashScope 接口没有返回任务 ID")
	}

	return taskID, nil
}

// pollDashScopeTask 轮询 DashScope 任务状态
func pollDashScopeTask(ctx context.Context, config providerConfig, taskID string) (dashScopeResponse, error) {
	var response dashScopeResponse
	if err := getDashScopeJSON(withProviderRequestKind(ctx, "poll"), config, "/api/v1/tasks/"+taskID, &response); err != nil {
		return response, err
	}
	return response, nil
}

// dashScopeImageDataURLs 下载 DashScope 图片 URL 并转为 data URL
func dashScopeImageDataURLs(ctx context.Context, config providerConfig, response dashScopeResponse) ([]string, error) {
	results := response.Output.Results
	if len(results) == 0 {
		return nil, errors.New("DashScope 接口没有返回图片")
	}

	images := make([]string, 0, len(results))
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
		images = append(images, dataURL(mimeType, data))
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

// DashScope API 响应结构
type dashScopeResponse struct {
	Output struct {
		TaskID     string `json:"task_id"`
		TaskStatus string `json:"task_status"` // PENDING, RUNNING, SUCCEEDED, FAILED
		Message    string `json:"message"`
		Results    []struct {
			URL string `json:"url"`
		} `json:"results"`
	} `json:"output"`
	RequestID string `json:"request_id"`
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
	// DashScope 轮询接口不需要 /v1 前缀，直接使用完整路径
	url := strings.TrimRight(config.BaseURL, "/") + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	ApplyOutboundHeaders(req, config.Headers)
	return doJSON(req, target)
}
