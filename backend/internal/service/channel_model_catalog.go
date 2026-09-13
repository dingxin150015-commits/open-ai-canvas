package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"infinite-canvas/backend/internal/model"
)

type ChannelModelsRequest struct {
	BaseURL           string           `json:"baseUrl"`
	AllowLocalChannel bool             `json:"allowLocalChannel"`
	APIKey            string           `json:"apiKey"`
	APIFormat         string           `json:"apiFormat"`
	Headers           []OutboundHeader `json:"headers"`
}

type channelModelsPayload struct {
	NextPageToken string             `json:"nextPageToken"`
	Data          []channelModelItem `json:"data"`
	Models        []channelModelItem `json:"models"`
	Error         *providerError     `json:"error"`
	Code          *int               `json:"code"`
	Msg           string             `json:"msg"`
}

type channelModelItem struct {
	GeminiDisplayName      string                        `json:"displayName"`
	GenerationMethods      []string                      `json:"supportedGenerationMethods"`
	ID                     string                        `json:"id"`
	Name                   string                        `json:"name"`
	DisplayName            string                        `json:"display_name"`
	ModelType              string                        `json:"model_type"`
	SupportedEndpointTypes []string                      `json:"supported_endpoint_types"`
	DefaultParameters      channelModelCatalogParameters `json:"default_parameters"`
	Options                channelModelCatalogOptions    `json:"options"`
	SupportsImages         *bool                         `json:"supports_images"`
	MinImages              *int                          `json:"min_images"`
	MaxImages              *int                          `json:"max_images"`
}

type channelModelCatalogParameters struct {
	AspectRatio     string `json:"aspect_ratio"`
	DurationSeconds string `json:"duration_seconds"`
	Resolution      string `json:"resolution"`
}

type channelModelCatalogOptions struct {
	AspectRatio     []ChannelModelCatalogOption `json:"aspect_ratio"`
	DurationSeconds []ChannelModelCatalogOption `json:"duration_seconds"`
	Resolution      []ChannelModelCatalogOption `json:"resolution"`
}

func (s *Service) FetchChannelModels(ctx context.Context, actor *model.User, input ChannelModelsRequest) ([]string, error) {
	items, err := s.FetchChannelModelCatalog(ctx, actor, input)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]bool, len(items))
	models := make([]string, 0, len(items))
	for _, item := range items {
		name := strings.TrimSpace(item.ID)
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		models = append(models, name)
	}
	sort.Strings(models)
	return models, nil
}

// ChannelModelCatalogItem 是前端自定义渠道拉取模型目录后的最小合同；
// 协议、能力和可选参数均来自上游公开元数据，不展开供应商内部兼容模型。
type ChannelModelCatalogItem struct {
	ID                     string                               `json:"id"`
	DisplayName            string                               `json:"displayName,omitempty"`
	ProviderModelKey       string                               `json:"providerModelKey,omitempty"`
	ModelType              string                               `json:"modelType,omitempty"`
	Protocol               string                               `json:"protocol,omitempty"`
	SupportedEndpointTypes []string                             `json:"supportedEndpointTypes,omitempty"`
	SupportedOperations    []string                             `json:"supportedOperations,omitempty"`
	SupportStatus          model.ChannelModelSupportStatus      `json:"supportStatus"`
	SupportReason          string                               `json:"supportReason,omitempty"`
	CatalogSource          string                               `json:"catalogSource,omitempty"`
	CatalogVersion         string                               `json:"catalogVersion,omitempty"`
	DocumentationPaths     []string                             `json:"documentationPaths,omitempty"`
	APIPath                string                               `json:"apiPath,omitempty"`
	DefaultParameters      ChannelModelCatalogDefaultParameters `json:"defaultParameters,omitempty"`
	Options                ChannelModelCatalogOptions           `json:"options,omitempty"`
	SupportsImages         *bool                                `json:"supportsImages,omitempty"`
	MinImages              *int                                 `json:"minImages,omitempty"`
	MaxImages              *int                                 `json:"maxImages,omitempty"`
}

type ChannelModelCatalogDefaultParameters struct {
	AspectRatio     string `json:"aspectRatio,omitempty"`
	DurationSeconds string `json:"durationSeconds,omitempty"`
	Resolution      string `json:"resolution,omitempty"`
}

type ChannelModelCatalogOptions struct {
	AspectRatio     []ChannelModelCatalogOption `json:"aspectRatio,omitempty"`
	DurationSeconds []ChannelModelCatalogOption `json:"durationSeconds,omitempty"`
	Resolution      []ChannelModelCatalogOption `json:"resolution,omitempty"`
}

type ChannelModelCatalogOption struct {
	Value string `json:"value"`
	Label string `json:"label,omitempty"`
}

func (s *Service) FetchChannelModelCatalog(ctx context.Context, actor *model.User, input ChannelModelsRequest) ([]ChannelModelCatalogItem, error) {
	if actor == nil || strings.TrimSpace(actor.ID) == "" {
		return nil, Unauthorized("请先登录")
	}
	baseURL := strings.TrimRight(strings.TrimSpace(input.BaseURL), "/")
	apiKey := strings.TrimSpace(input.APIKey)
	if baseURL == "" {
		return nil, BadAuthRequest("请填写 Base URL")
	}
	if apiKey == "" {
		return nil, BadAuthRequest("请填写 API Key")
	}
	apiFormat := strings.ToLower(strings.TrimSpace(input.APIFormat))
	if apiFormat == "" {
		apiFormat = "openai"
	}
	if apiFormat != "openai" && apiFormat != "gemini" {
		return nil, BadAuthRequest("接口协议不支持拉取模型")
	}
	headers, err := NormalizeOutboundHeaders(input.Headers)
	if err != nil {
		return nil, err
	}

	target := apiURL(baseURL, "/models")
	if apiFormat == "gemini" {
		if !strings.HasSuffix(strings.ToLower(baseURL), "/v1beta") {
			baseURL += "/v1beta"
		}
		target = baseURL + "/models"
	}
	if _, err := s.validateChannelOutboundURL(target, input.AllowLocalChannel, false); err != nil {
		return nil, err
	}
	requestContext := withProviderOutboundPolicy(ctx, providerConfig{BaseURL: baseURL, AllowLocalChannel: s.effectiveAllowLocalChannel(input.AllowLocalChannel)})
	request, err := http.NewRequestWithContext(requestContext, http.MethodGet, target, nil)
	if err != nil {
		return nil, BadAuthRequest("模型服务地址无效")
	}
	if apiFormat == "gemini" {
		request.Header.Set("x-goog-api-key", apiKey)
	} else {
		request.Header.Set("Authorization", "Bearer "+apiKey)
	}
	ApplyOutboundHeaders(request, headers)

	// 只代理固定的模型目录 GET；用户密钥仅用于本次请求，不写入数据库或日志。
	items, err := fetchChannelCatalogPages(request, apiFormat, doBinary)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]bool, len(items))
	catalog := make([]ChannelModelCatalogItem, 0, len(items))
	for _, item := range items {
		name := strings.TrimPrefix(strings.TrimSpace(firstNonEmpty(item.ID, item.Name)), "models/")
		if name == "" || seen[name] {
			continue
		}
		seen[name] = true
		modelType := normalizeCatalogModelType(item.ModelType)
		endpointTypes := normalizeCatalogEndpointTypes(item.SupportedEndpointTypes)
		catalog = append(catalog, ChannelModelCatalogItem{
			ID:                     name,
			DisplayName:            strings.TrimSpace(firstNonEmpty(item.GeminiDisplayName, item.DisplayName)),
			ProviderModelKey:       name,
			ModelType:              modelType,
			Protocol:               catalogProtocol(modelType, endpointTypes, apiFormat),
			SupportedEndpointTypes: endpointTypes,
			SupportStatus:          model.ChannelModelSupportUnknown,
			SupportReason:          catalogConfigurationReason(item.GenerationMethods),
			CatalogSource:          "upstream",
			CatalogVersion:         "upstream",
			DefaultParameters: ChannelModelCatalogDefaultParameters{
				AspectRatio:     strings.TrimSpace(item.DefaultParameters.AspectRatio),
				DurationSeconds: strings.TrimSpace(item.DefaultParameters.DurationSeconds),
				Resolution:      strings.TrimSpace(item.DefaultParameters.Resolution),
			},
			Options: ChannelModelCatalogOptions{
				AspectRatio:     normalizeCatalogOptions(item.Options.AspectRatio),
				DurationSeconds: normalizeCatalogOptions(item.Options.DurationSeconds),
				Resolution:      normalizeCatalogOptions(item.Options.Resolution),
			},
			SupportsImages: item.SupportsImages,
			MinImages:      item.MinImages,
			MaxImages:      item.MaxImages,
		})
	}
	sort.Slice(catalog, func(left int, right int) bool {
		return catalog[left].ID < catalog[right].ID
	})
	// 官方厂商 Manifest 是目录事实的一部分，不受实验插件开关控制；
	// 支持状态仍会阻止未完成 Adapter 的模型被定价、启用或路由。
	catalog = extendChannelModelCatalogWithRegistry(s.protocolRegistry(), baseURL, catalog)
	return catalog, nil
}

func catalogConfigurationReason(methods []string) string {
	if len(methods) == 0 {
		return "已发现，待管理员配置能力与请求协议"
	}
	return "上游方法：" + strings.Join(methods, "、") + "；请配置输出能力与协议"
}

// A listing is atomic: partial pages never reach catalog synchronization.
func fetchChannelCatalogPages(request *http.Request, format string, fetch func(*http.Request) ([]byte, string, error)) ([]channelModelItem, error) {
	ctx, cancel := context.WithTimeout(request.Context(), 2*time.Minute)
	defer cancel()
	result := []channelModelItem{}
	seen := map[string]bool{}
	cursor := ""
	for page := 0; page < 100; page++ {
		next := request.Clone(ctx)
		parsed := *request.URL
		next.URL = &parsed
		if cursor != "" {
			query := next.URL.Query()
			query.Set("pageToken", cursor)
			next.URL.RawQuery = query.Encode()
		}
		data, _, err := fetch(next)
		if err != nil {
			return nil, channelModelsUpstreamError(err)
		}
		var payload channelModelsPayload
		if err := json.Unmarshal(data, &payload); err != nil {
			return nil, WrapAppError(http.StatusBadGateway, "模型目录响应无效，未导入部分结果", err)
		}
		if payload.Error != nil || (payload.Code != nil && *payload.Code != 0) {
			return nil, NewAppError(http.StatusBadGateway, "模型目录拉取失败，未导入部分结果")
		}
		items := payload.Data
		if format == "gemini" {
			items = payload.Models
		}
		result = append(result, items...)
		if len(result) > 10000 {
			return nil, BadAuthRequest("模型目录超过数量上限，未导入部分结果")
		}
		cursor = strings.TrimSpace(payload.NextPageToken)
		if format != "gemini" || cursor == "" {
			return result, nil
		}
		if seen[cursor] || len(cursor) > 4096 {
			return nil, BadAuthRequest("模型目录分页游标异常，未导入部分结果")
		}
		seen[cursor] = true
	}
	return nil, BadAuthRequest("模型目录超过分页上限，未导入部分结果")
}

func catalogProtocol(modelType string, endpointTypes []string, apiFormat string) string {
	normalized := make(map[string]bool, len(endpointTypes))
	for _, endpointType := range endpointTypes {
		normalized[strings.ToLower(strings.TrimSpace(endpointType))] = true
	}
	switch {
	case normalized["openai-chat"] || normalized["chat-completion"] || normalized["chat"]:
		return string(model.ChannelInterfaceChatCompletion)
	case normalized["openai-response"] || normalized["responses"]:
		return string(model.ChannelInterfaceOpenAIResponse)
	case normalized["openai-image"] || normalized["image"]:
		return string(model.ChannelInterfaceOpenAIImage)
	case normalized["openai-video"] || normalized["video"]:
		return string(model.ChannelInterfaceNewAPIChannel2)
	case normalized["openai-audio"] || normalized["audio"]:
		return string(model.ChannelInterfaceOpenAIAudio)
	}
	if strings.EqualFold(strings.TrimSpace(apiFormat), "openai") {
		switch modelType {
		case "text":
			return string(model.ChannelInterfaceChatCompletion)
		case "image":
			return string(model.ChannelInterfaceOpenAIImage)
		case "video":
			return string(model.ChannelInterfaceNewAPIChannel2)
		case "audio":
			return string(model.ChannelInterfaceOpenAIAudio)
		}
	}
	return ""
}

func normalizeCatalogModelType(value string) string {
	normalized := strings.ToLower(strings.TrimSpace(value))
	switch normalized {
	case "text", "image", "video", "audio":
		return normalized
	default:
		return ""
	}
}

func normalizeCatalogOptions(options []ChannelModelCatalogOption) []ChannelModelCatalogOption {
	seen := make(map[string]bool, len(options))
	normalized := make([]ChannelModelCatalogOption, 0, len(options))
	for _, option := range options {
		value := strings.TrimSpace(option.Value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		normalized = append(normalized, ChannelModelCatalogOption{Value: value, Label: strings.TrimSpace(option.Label)})
	}
	return normalized
}

func normalizeCatalogEndpointTypes(values []string) []string {
	seen := make(map[string]bool, len(values))
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		item := strings.TrimSpace(value)
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		normalized = append(normalized, item)
	}
	return normalized
}

func channelModelsUpstreamError(err error) error {
	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr
	}
	var httpErr providerHTTPError
	if !errors.As(err, &httpErr) {
		return WrapAppError(http.StatusBadGateway, "连接模型服务失败，请检查渠道地址和网络", err)
	}
	switch httpErr.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return NewAppError(http.StatusBadGateway, "模型服务鉴权失败，请检查 API Key")
	case http.StatusNotFound:
		return NewAppError(http.StatusBadGateway, "模型服务未提供 /models 接口")
	case http.StatusTooManyRequests:
		return NewAppError(http.StatusBadGateway, "模型服务请求过于频繁或额度不足")
	default:
		return WrapAppError(http.StatusBadGateway, httpErr.Error(), err)
	}
}
