package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// runDashScopeVideoTask 实现百炼（DashScope）视频生成协议（异步模式）
//
// API 文档：
// - Wan 2.7 文生视频：https://platform.qianwenai.com/docs/api-reference/video-generation/wan27-text-to-video/create-task
// - Wan 2.7 图生视频：https://platform.qianwenai.com/docs/api-reference/video-generation/wan27-image-to-video/create-task
// - Wan 2.7 参考生视频：https://platform.qianwenai.com/docs/api-reference/video-generation/wan27-reference-to-video/create-task
// - HappyHorse 系列：https://platform.qianwenai.com/docs/api-reference/video-generation/happyhorse-text-to-video/create-task
// - Wan 3.0：https://platform.qianwenai.com/docs/api-reference/video-generation/wan30-video/create-task
//
// 支持的模型和能力：
// - wan2.7-t2v: 文生视频，2-15秒，720P/1080P
// - wan2.7-i2v: 图生视频（首帧），1张输入图像
// - wan2.7-r2v: 参考生视频，最多5个参考素材（图片或视频）
// - happyhorse-1.1-t2v/i2v/r2v: HappyHorse 系列，物理真实
// - wan3.0-video: 全能模型
//
// API 流程：
//  1. POST /services/aigc/video-generation/video-synthesis (Header: X-DashScope-Async: enable)
//     → 返回 task_id
//  2. 轮询 GET /tasks/{task_id}，每15秒一次
//     → task_status: PENDING → RUNNING → SUCCEEDED / FAILED
//  3. 任务成功后从 video_url 下载视频（24小时有效期）
//  4. 转换为 data URL 返回，由前端资源同步链写入平台 OSS
//
// 参数限制由配置系统管理，在 validateVideoTask 中统一验证
func runDashScopeVideoTask(ctx context.Context, input canvasGenerationInput) (map[string]interface{}, error) {
	log.Printf("[DashScope Video] ========== 开始处理视频生成任务 ==========")
	log.Printf("[DashScope Video] Model: %s", input.Config.Model)

	// 1. 构建请求体（自动识别 t2v/i2v/r2v 模式）
	requestBody, err := buildDashScopeVideoRequest(input)
	if err != nil {
		log.Printf("[DashScope Video] ❌ 构建请求体失败: %v", err)
		return nil, err
	}

	// 2. 提交异步任务
	taskID, err := submitDashScopeVideoTask(ctx, input.Config, requestBody)
	if err != nil {
		log.Printf("[DashScope Video] ❌ 提交任务失败: %v", err)
		return nil, fmt.Errorf("提交视频生成任务失败：%w", err)
	}
	log.Printf("[DashScope Video] ✅ 任务提交成功，TaskID: %s", taskID)

	// 3. 轮询任务状态
	// 超时上限跟随系统任务策略（ctx deadline，默认 30 分钟），与其他视频协议保持一致，
	// 不再自行计算短上限，否则会在 DashScope 侧仍为 RUNNING 时提前放弃。
	result, err := pollDashScopeVideoTask(ctx, input.Config, taskID)
	if err != nil {
		log.Printf("[DashScope Video] ❌ 轮询任务失败: %v", err)
		return nil, err
	}
	log.Printf("[DashScope Video] ✅ 任务完成，已收到临时视频地址")

	// 5. 下载视频并交给统一资源化链
	videoDataURL, mimeType, err := downloadDashScopeVideo(ctx, input.Config, result.Output.VideoURL)
	if err != nil {
		log.Printf("[DashScope Video] ❌ 下载视频失败: %v", err)
		return nil, fmt.Errorf("下载视频失败：%w", err)
	}
	log.Printf("[DashScope Video] ✅ 视频下载成功，MimeType: %s", mimeType)

	// 6. 返回结果（格式与现有视频协议一致）
	log.Printf("[DashScope Video] ========== 任务完成 ==========")
	return map[string]interface{}{
		"mode": "video",
		"video": map[string]interface{}{
			"dataUrl":  videoDataURL,
			"mimeType": mimeType,
		},
	}, nil
}

// buildDashScopeVideoRequest 构建视频生成请求体
// 根据模型ID和输入自动识别 t2v/i2v/r2v 模式
func buildDashScopeVideoRequest(input canvasGenerationInput) (map[string]interface{}, error) {
	model := input.Config.Model
	if isWan30VideoModel(model) || strings.HasPrefix(strings.ToLower(strings.TrimPrefix(strings.TrimSpace(model), "models/")), "wan3.0-") {
		return buildWan30VideoRequest(input)
	}
	if isWan27ReadyVideoModel(model) || strings.HasPrefix(strings.ToLower(strings.TrimPrefix(strings.TrimSpace(model), "models/")), "wan2.7-") {
		return buildWan27VideoRequest(input)
	}
	if isHappyHorse11ReadyVideoModel(model) || strings.HasPrefix(strings.ToLower(strings.TrimPrefix(strings.TrimSpace(model), "models/")), "happyhorse-") {
		return buildHappyHorseVideoRequest(input)
	}
	log.Printf("[DashScope Video] --- 开始构建请求体 ---")
	log.Printf("[DashScope Video] Input Model: %s", model)

	// 参数转换前后对比
	// 视频参数统一从 input.Config 读取，与其他视频协议保持一致；
	// metadata 只承载业务上下文（会话、节点等），不含生成参数。
	durationNorm := normalizeDashScopeVideoDuration(input.Config.VideoSeconds)
	log.Printf("[DashScope Video] Duration: Config.VideoSeconds=%q -> %d", input.Config.VideoSeconds, durationNorm)

	resolutionNorm := normalizeDashScopeVideoResolution(input.Config.VQuality)
	log.Printf("[DashScope Video] Resolution: Config.VQuality=%q -> %s", input.Config.VQuality, resolutionNorm)

	ratioNorm := normalizeDashScopeVideoRatio(input.Config.Size)
	log.Printf("[DashScope Video] Ratio: Config.Size=%q -> %s", input.Config.Size, ratioNorm)

	watermark := parseBool(input.Config.VideoWatermark, false)
	log.Printf("[DashScope Video] Watermark: Config.VideoWatermark=%q -> %v", input.Config.VideoWatermark, watermark)

	// 构建基础请求体
	body := map[string]interface{}{
		"model": model,
		"input": map[string]interface{}{
			"prompt": input.Prompt,
		},
		"parameters": map[string]interface{}{
			"duration":      durationNorm,
			"resolution":    resolutionNorm,
			"ratio":         ratioNorm,
			"watermark":     watermark,
			"prompt_extend": true, // 默认开启提示词扩写
		},
	}

	// 构建 media 数组（根据模型和输入自动判断）
	media, err := buildDashScopeVideoMedia(model, input)
	if err != nil {
		return nil, err
	}
	log.Printf("[DashScope Video] Media count: %d", len(media))

	// 参考素材数量下限由模型能力配置（VideoReferenceConfig.MinImages）经 validateVideoTask 统一校验，
	// 与其他视频协议保持一致，此处不再按模型名硬编码判断。

	// 如果有 media，添加到请求体
	// ✅ 修复问题3：添加类型断言防御，避免panic
	if len(media) > 0 {
		if inputMap, ok := body["input"].(map[string]interface{}); ok {
			inputMap["media"] = media
			log.Printf("[DashScope Video] ✅ Media 已添加到请求体")
		} else {
			log.Printf("[DashScope Video] ⚠️ 类型断言失败，无法添加 media")
		}
	}

	// 处理负向提示词（如果有）
	if negativePrompt := metadataString(input.Metadata, "negative_prompt"); negativePrompt != "" {
		if inputMap, ok := body["input"].(map[string]interface{}); ok {
			inputMap["negative_prompt"] = negativePrompt
			log.Printf("[DashScope Video] ✅ 已添加负向提示词")
		}
	}

	// 处理声音克隆（仅 Wan 2.7 r2v 支持，HappyHorse 不支持）
	// 优先从 Metadata["reference_voice"] 读取（如果前端实现了专用上传）
	// 降级到 ReferenceAudios[0]（复用现有音频素材输入）
	isWan27 := strings.Contains(model, "wan2.7") || strings.Contains(model, "wan27")
	if strings.Contains(model, "r2v") && isWan27 {
		referenceVoice := metadataString(input.Metadata, "reference_voice")
		if referenceVoice == "" && len(input.ReferenceAudios) > 0 {
			// 从第一个音频素材提取 URL
			if url, err := mediaReferenceURL(input.ReferenceAudios[0]); err == nil {
				referenceVoice = url
			}
		}
		if referenceVoice != "" {
			if inputMap, ok := body["input"].(map[string]interface{}); ok {
				inputMap["reference_voice"] = referenceVoice
				log.Printf("[DashScope Video] ✅ 已添加参考音色")
			}
		}
	}

	return body, nil
}

// buildDashScopeVideoMedia 构建 media 数组
// 自动判断素材类型（first_frame、reference_image、reference_video）
// 素材地址统一走 mediaReferenceURL（与其他 JSON 视频协议一致）。
// 参考素材已在 hydrateGenerationMedia 阶段换成对象存储签名 URL，这里不再自行编码 data URL；
// 取不到地址时直接返回错误，不能静默跳过，否则会退化成 media 为空的 "Field required: input.media"。
func buildDashScopeVideoMedia(model string, input canvasGenerationInput) ([]map[string]interface{}, error) {
	var media []map[string]interface{}

	// 判断模式和模型系列
	isI2V := strings.Contains(model, "i2v")
	isR2V := strings.Contains(model, "r2v")
	isWan27 := strings.Contains(model, "wan2.7") || strings.Contains(model, "wan27")
	isHappyHorse := strings.Contains(model, "happyhorse")

	// 处理参考图片
	if isI2V {
		// i2v 模式
		imageCount := len(input.ReferenceImages)
		if imageCount > 0 {
			// 首帧（必须）
			url, err := mediaReferenceURL(input.ReferenceImages[0])
			if err != nil {
				return nil, fmt.Errorf("读取首帧图地址失败：%w", err)
			}
			media = append(media, map[string]interface{}{"type": "first_frame", "url": url})
			log.Printf("[DashScope Video] media[%d] type=first_frame", len(media)-1)

			// 末帧（仅 Wan 2.7 支持，HappyHorse 不支持）
			if isWan27 && imageCount >= 2 {
				url, err := mediaReferenceURL(input.ReferenceImages[1])
				if err != nil {
					return nil, fmt.Errorf("读取末帧图地址失败：%w", err)
				}
				media = append(media, map[string]interface{}{"type": "last_frame", "url": url})
				log.Printf("[DashScope Video] media[%d] type=last_frame", len(media)-1)
			}
		}
	} else {
		// t2v / r2v 模式：处理参考图片
		maxImages := 5 // Wan 2.7 r2v 最多5张
		if isHappyHorse && isR2V {
			maxImages = 9 // HappyHorse r2v 支持1-9张
		}

		for i, refImg := range input.ReferenceImages {
			if i >= maxImages {
				break
			}
			url, err := mediaReferenceURL(refImg)
			if err != nil {
				return nil, fmt.Errorf("读取第 %d 张参考图地址失败：%w", i+1, err)
			}

			mediaType := "reference_image"
			media = append(media, map[string]interface{}{"type": mediaType, "url": url})
			log.Printf("[DashScope Video] media[%d] type=%s", len(media)-1, mediaType)
		}
	}

	// 处理参考视频（仅 Wan 2.7 r2v 支持，HappyHorse 不支持）
	if isR2V && isWan27 {
		for i, refVid := range input.ReferenceVideos {
			if len(media) >= 5 {
				break
			}
			url, err := mediaReferenceURL(refVid)
			if err != nil {
				return nil, fmt.Errorf("读取第 %d 个参考视频地址失败：%w", i+1, err)
			}
			media = append(media, map[string]interface{}{"type": "reference_video", "url": url})
			log.Printf("[DashScope Video] media[%d] type=reference_video", len(media)-1)
		}
	}

	// 处理驱动音频（仅 Wan 2.7 i2v 支持，HappyHorse 不支持）
	if isI2V && isWan27 {
		for _, refAudio := range input.ReferenceAudios {
			url, err := mediaReferenceURL(refAudio)
			if err != nil {
				return nil, fmt.Errorf("读取驱动音频地址失败：%w", err)
			}
			media = append(media, map[string]interface{}{"type": "driving_audio", "url": url})
			log.Printf("[DashScope Video] media[%d] type=driving_audio", len(media)-1)
			break // 只支持 1 个驱动音频
		}
	}

	return media, nil
}

// submitDashScopeVideoTask 提交异步视频生成任务
func submitDashScopeVideoTask(ctx context.Context, config providerConfig, requestBody map[string]interface{}) (string, error) {
	startTime := time.Now()
	path := "/api/v1/services/aigc/video-generation/video-synthesis"

	log.Printf("[DashScope Video] --- 开始提交任务 ---")

	// 检查 context deadline
	if deadline, ok := ctx.Deadline(); ok {
		log.Printf("[DashScope Video] Context deadline: %v (剩余 %v)", deadline, time.Until(deadline))
	} else {
		log.Printf("[DashScope Video] Context 无 deadline")
	}

	data, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("[DashScope Video] ❌ 序列化请求体失败: %v", err)
		return "", fmt.Errorf("序列化请求体失败：%w", err)
	}
	log.Printf("[DashScope Video] ✅ 请求体序列化成功，大小: %d bytes", len(data))

	url, err := dashScopeNativeEndpoint(config.BaseURL, path)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		log.Printf("[DashScope Video] ❌ 创建HTTP请求失败: %v", err)
		return "", err
	}

	// 注意：真实 API Key 必须原样写入 Header，脱敏只用于日志输出
	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DashScope-Async", "enable") // 启用异步模式
	ApplyOutboundHeaders(req, config.Headers)

	log.Printf("[DashScope Video] HTTP Headers:")
	log.Printf("[DashScope Video]   Content-Type: %s", req.Header.Get("Content-Type"))
	log.Printf("[DashScope Video]   X-DashScope-Async: %s", req.Header.Get("X-DashScope-Async"))

	log.Printf("[DashScope Video] 开始调用 doJSON...")
	var response dashScopeVideoSubmitResponse
	if err := doJSON(req, &response); err != nil {
		elapsed := time.Since(startTime)
		log.Printf("[DashScope Video] ❌ doJSON 失败 (耗时 %v): %v", elapsed, err)
		return "", err
	}
	elapsed := time.Since(startTime)
	log.Printf("[DashScope Video] ✅ doJSON 成功 (耗时 %v)", elapsed)

	log.Printf("[DashScope Video] 响应已接收 request_id=%s task_status=%s has_task_id=%v", response.RequestID, response.Output.TaskStatus, response.Output.TaskID != "")

	// 检查是否有错误
	if response.Code != "" && response.Code != "Success" {
		log.Printf("[DashScope Video] ❌ API返回错误码: %s - %s", response.Code, response.Message)
		// 翻译 DashScope 的参数错误为用户可读提示
		if strings.EqualFold(strings.TrimSpace(response.Code), "InvalidParameter") && strings.Contains(strings.ToLower(response.Message), "input.media") {
			return "", fmt.Errorf("模型 %s 需要参考素材（图片/视频），请先添加后再生成", config.Model)
		}
		return "", fmt.Errorf("DashScope API 错误：%s - %s", response.Code, response.Message)
	}

	taskID := response.Output.TaskID
	if taskID == "" {
		log.Printf("[DashScope Video] ❌ API未返回任务ID")
		return "", errors.New("API 未返回任务 ID")
	}

	log.Printf("[DashScope Video] ✅ 任务提交成功，TaskID: %s", taskID)
	return taskID, nil
}

// pollDashScopeVideoTask 轮询视频任务状态
// 轮询间隔 15 秒（官方 query-result.md 建议值）；
// 超时上限跟随系统任务策略（ctx deadline，默认 30 分钟），与其他视频协议统一。
func pollDashScopeVideoTask(ctx context.Context, config providerConfig, taskID string) (*dashScopeVideoTaskResult, error) {
	const pollInterval = 15 * time.Second // 官方建议15秒
	const initialDelay = 5 * time.Second  // 首次延迟5秒

	deadline := providerPollingDeadline(ctx)
	log.Printf("[DashScope Video] --- 开始轮询任务状态 ---")
	log.Printf("[DashScope Video] TaskID: %s", taskID)
	log.Printf("[DashScope Video] 轮询截止时间: %v（剩余 %v）", deadline, time.Until(deadline))

	// 首次延迟，避免立即轮询
	log.Printf("[DashScope Video] 等待 %v 后开始首次轮询...", initialDelay)
	if err := sleepContext(ctx, initialDelay); err != nil {
		return nil, err
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	// 在循环外创建一次超时 channel，避免每轮重复创建 timer
	timeoutCh := time.After(time.Until(deadline))

	pollStart := time.Now()
	round := 0

	for {
		select {
		case <-ctx.Done():
			log.Printf("[DashScope Video] ❌ 轮询被取消（已轮询 %d 轮，耗时 %v）：%v", round, time.Since(pollStart), ctx.Err())
			return nil, ctx.Err()

		case <-timeoutCh:
			// 超时：DashScope 侧任务仍可能在跑，task_id 24 小时内可查，提示用户稍后查询
			elapsed := time.Since(pollStart)
			log.Printf("[DashScope Video] ❌ 轮询超时（已轮询 %d 轮，耗时 %v）TaskID=%s", round, elapsed, taskID)
			return nil, fmt.Errorf("视频生成超时（%v），任务ID: %s，可稍后查询", elapsed.Round(time.Second), taskID)

		case <-ticker.C:
			round++
			log.Printf("[DashScope Video] 轮询第 %d 轮（已耗时 %v）...", round, time.Since(pollStart))
			result, err := fetchDashScopeVideoTaskStatus(ctx, config, taskID)
			if err != nil {
				log.Printf("[DashScope Video] ❌ 第 %d 轮查询失败：%v", round, err)
				return nil, err
			}
			log.Printf("[DashScope Video] 第 %d 轮 task_status=%q has_video_url=%v message=%q", round, result.Output.TaskStatus, result.Output.VideoURL != "", result.Output.Message)

			done, outcomeErr := dashScopeVideoTaskOutcome(config.Model, result)
			if outcomeErr != nil {
				return nil, outcomeErr
			}
			if done {
				log.Printf("[DashScope Video] ✅ 轮询成功（共 %d 轮，耗时 %v）", round, time.Since(pollStart))
				return result, nil
			}
		}
	}
}

func dashScopeVideoTaskOutcome(modelName string, result *dashScopeVideoTaskResult) (bool, error) {
	if result == nil {
		return false, errors.New("DashScope 视频任务状态响应为空")
	}
	switch result.Output.TaskStatus {
	case "SUCCEEDED":
		if strings.TrimSpace(result.Output.VideoURL) == "" {
			return false, errors.New("任务成功但未返回视频 URL")
		}
		return true, nil
	case "FAILED":
		message := strings.TrimSpace(result.Output.Message)
		if message == "" {
			message = "未知错误"
		}
		if strings.Contains(message, "Field required") && strings.Contains(message, "input.media") {
			return false, fmt.Errorf("模型 %s 需要参考素材（图片/视频），请先添加后再生成", modelName)
		}
		return false, fmt.Errorf("视频生成失败：%s", message)
	case "CANCELED":
		return false, errors.New("DashScope 视频任务已取消")
	case "UNKNOWN":
		return false, errors.New("DashScope 视频任务不存在或任务 ID 已过期")
	case "PENDING", "RUNNING":
		return false, nil
	default:
		return false, fmt.Errorf("未知任务状态：%s", result.Output.TaskStatus)
	}
}

// fetchDashScopeVideoTaskStatus 查询单次任务状态
func fetchDashScopeVideoTaskStatus(ctx context.Context, config providerConfig, taskID string) (*dashScopeVideoTaskResult, error) {
	path := fmt.Sprintf("/api/v1/tasks/%s", taskID)

	url, err := dashScopeNativeEndpoint(config.BaseURL, path)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.APIKey)
	ApplyOutboundHeaders(req, config.Headers)

	var result dashScopeVideoTaskResult
	if err := doJSON(req, &result); err != nil {
		return nil, err
	}

	// 检查API错误
	if result.Code != "" && result.Code != "Success" {
		return nil, fmt.Errorf("DashScope API 错误：%s - %s", result.Code, result.Message)
	}

	return &result, nil
}

// downloadDashScopeVideo 下载短期结果并转换为统一 data URL；持久化由前端资源同步链负责。
func downloadDashScopeVideo(ctx context.Context, _ providerConfig, videoURL string) (string, string, error) {
	// 创建下载超时上下文（5分钟）
	downloadCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	if _, err := ValidateOutboundURL(videoURL); err != nil {
		return "", "", fmt.Errorf("视频结果地址不安全：%w", err)
	}

	req, err := http.NewRequestWithContext(downloadCtx, http.MethodGet, videoURL, nil)
	if err != nil {
		return "", "", err
	}
	// 结果 URL 是短期签名地址，通常属于独立 OSS 主机。API Key 和渠道
	// 自定义请求头只允许发往模型端点，禁止跨主机转发到结果下载地址。

	// ✅ 修复问题5：使用自定义HTTP Client，设置超时
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("下载视频请求失败：%w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("下载视频失败，HTTP 状态码：%d", resp.StatusCode)
	}

	// ✅ 修复问题4：检查 Content-Length，避免下载过大文件浪费带宽
	const maxVideoSize = 500 * 1024 * 1024 // 500MB
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		if size, err := strconv.ParseInt(contentLength, 10, 64); err == nil {
			if size > maxVideoSize {
				return "", "", fmt.Errorf("视频文件过大（%d MB），超过500MB限制", size/(1024*1024))
			}
		}
	}

	// 当前实现：下载到内存（带500MB软限制）
	limitedReader := io.LimitReader(resp.Body, maxVideoSize)

	var buf bytes.Buffer
	written, err := io.Copy(&buf, limitedReader)
	if err != nil {
		return "", "", fmt.Errorf("读取视频数据失败：%w", err)
	}

	// 检查是否达到限制
	if written >= maxVideoSize {
		// 尝试再读一个字节，确认是否真的超限
		extra := make([]byte, 1)
		n, _ := resp.Body.Read(extra)
		if n > 0 {
			return "", "", fmt.Errorf("视频文件超过500MB限制")
		}
	}

	// 获取MIME类型
	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "video/mp4" // DashScope 默认返回 MP4
	}

	// 转换为 data URL
	videoData := buf.Bytes()
	return dataURL(mimeType, videoData), mimeType, nil
}

// normalizeDashScopeVideoDuration 规范化视频时长（wan2.7 支持 2-15 秒整数）
// 入参取自 providerConfig.VideoSeconds（字符串），空值或非法值回退到 5 秒。
func normalizeDashScopeVideoDuration(value string) int {
	duration, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || duration <= 0 {
		return 5 // 与其他视频协议一致的默认时长
	}
	if duration < 2 {
		return 2
	}
	if duration > 15 {
		return 15
	}
	return duration
}

// normalizeDashScopeVideoResolution 规范化分辨率（wan2.7 仅支持 720P / 1080P）
// 入参取自 providerConfig.VQuality，如 "720p"、"1080p"、"auto"。
func normalizeDashScopeVideoResolution(value string) string {
	if strings.Contains(strings.ToUpper(strings.TrimSpace(value)), "1080") {
		return "1080P"
	}
	return "720P"
}

// normalizeDashScopeVideoRatio 规范化画面比例（wan2.7 支持 16:9/9:16/1:1/4:3/3:4）
// 入参取自 providerConfig.Size，既可能是 "16:9" 这类比例，也可能是 "1280x720" 这类像素尺寸。
func normalizeDashScopeVideoRatio(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	allowed := map[string]bool{"16:9": true, "9:16": true, "1:1": true, "4:3": true, "3:4": true}
	if allowed[normalized] {
		return normalized
	}
	// 像素尺寸（如 1280x720）按最接近的受支持比例折算，避免直接丢弃用户选择。
	parts := strings.Split(strings.ReplaceAll(normalized, "*", "x"), "x")
	if len(parts) != 2 {
		return "16:9"
	}
	width, widthErr := strconv.Atoi(strings.TrimSpace(parts[0]))
	height, heightErr := strconv.Atoi(strings.TrimSpace(parts[1]))
	if widthErr != nil || heightErr != nil || width <= 0 || height <= 0 {
		return "16:9"
	}
	ratio := float64(width) / float64(height)
	candidates := []struct {
		name  string
		value float64
	}{
		{name: "16:9", value: 16.0 / 9},
		{name: "9:16", value: 9.0 / 16},
		{name: "1:1", value: 1},
		{name: "4:3", value: 4.0 / 3},
		{name: "3:4", value: 3.0 / 4},
	}
	bestName := "16:9"
	bestDifference := math.MaxFloat64
	for _, candidate := range candidates {
		if difference := math.Abs(ratio - candidate.value); difference < bestDifference {
			bestName = candidate.name
			bestDifference = difference
		}
	}
	return bestName
}

// DashScope 视频任务提交响应
type dashScopeVideoSubmitResponse struct {
	RequestID string `json:"request_id"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Output    struct {
		TaskID     string `json:"task_id"`
		TaskStatus string `json:"task_status"` // PENDING
	} `json:"output"`
}

// DashScope 视频任务查询响应
type dashScopeVideoTaskResult struct {
	RequestID string `json:"request_id"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	Output    struct {
		TaskID     string `json:"task_id"`
		TaskStatus string `json:"task_status"` // PENDING, RUNNING, SUCCEEDED, FAILED
		Message    string `json:"message,omitempty"`
		VideoURL   string `json:"video_url,omitempty"` // 成功时返回
		Usage      struct {
			Duration   int    `json:"duration"`    // 视频时长（秒）
			VideoCount int    `json:"video_count"` // 固定为1
			Ratio      string `json:"ratio"`       // 如 "16:9"
			SR         int    `json:"SR"`          // 720 或 1080
		} `json:"usage,omitempty"`
	} `json:"output"`
}
