package app

import (
	"net/url"
	"sort"
	"strings"

	"yingce/backend/internal/model"
	"yingce/backend/internal/protocol"
	"yingce/backend/internal/provider"
)

// extendChannelModelCatalog 在已通过统一安全出站链路获取的标准目录上补充厂商条目。
// 插件不持有密钥也不自行发 HTTP，避免绕过 SSRF、请求头和响应大小边界。
func extendChannelModelCatalog(baseURL string, apiFormat string, headers []OutboundHeader, catalog []ChannelModelCatalogItem) []ChannelModelCatalogItem {
	return extendChannelModelCatalogWithRegistry(loadOfficialFallbackRegistry(), baseURL, catalog)
}

func extendChannelModelCatalogWithRegistry(registry *protocol.Registry, baseURL string, catalog []ChannelModelCatalogItem) []ChannelModelCatalogItem {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return catalog
	}
	host := strings.ToLower(parsed.Hostname())
	extra := []provider.Model{}
	for _, metadata := range registry.List(protocol.SurfaceAdminSystemChannel, "", false) {
		if metadata.ModelCatalog == nil || len(metadata.Categories) != 1 {
			continue
		}
		matched := false
		for _, allowed := range metadata.ModelCatalog.Hosts {
			if host == allowed || strings.HasSuffix(host, "."+allowed) {
				matched = true
			}
		}
		if !matched {
			continue
		}
		for _, entry := range metadata.ModelCatalog.Entries {
			extra = append(extra, provider.Model{ID: entry.ID, ProviderModelKey: entry.ID, DisplayName: entry.DisplayName, Capability: []string{string(metadata.Categories[0])}, Protocol: metadata.ID, SupportStatus: provider.ModelSupportStatus(entry.SupportStatus), SupportReason: entry.SupportReason, SupportedOperations: entry.SupportedOperations, DocumentationPaths: entry.DocumentationPaths, CatalogSource: entry.CatalogSource, CatalogVersion: entry.CatalogVersion, APIPath: entry.APIPath})
		}
	}

	indexByID := make(map[string]int, len(catalog)+len(extra))
	for index := range catalog {
		indexByID[catalog[index].ID] = index
	}
	for _, model := range extra {
		item := providerCatalogItem(model)
		if item.ID == "" {
			continue
		}
		if index, exists := indexByID[item.ID]; exists {
			catalog[index] = enrichCatalogItem(catalog[index], item)
			continue
		}
		indexByID[item.ID] = len(catalog)
		catalog = append(catalog, item)
	}
	sort.Slice(catalog, func(left int, right int) bool {
		return catalog[left].ID < catalog[right].ID
	})
	return catalog
}

func providerCatalogItem(item provider.Model) ChannelModelCatalogItem {
	modelType := ""
	for _, capability := range item.Capability {
		if normalized := normalizeCatalogModelType(capability); normalized != "" {
			modelType = normalized
			break
		}
	}
	return ChannelModelCatalogItem{
		ID:                     strings.TrimPrefix(strings.TrimSpace(item.ID), "models/"),
		DisplayName:            strings.TrimSpace(item.DisplayName),
		ProviderModelKey:       strings.TrimPrefix(strings.TrimSpace(firstNonEmpty(item.ProviderModelKey, item.ID)), "models/"),
		ModelType:              modelType,
		Protocol:               strings.TrimSpace(item.Protocol),
		SupportedEndpointTypes: normalizeCatalogEndpointTypes(item.SupportedEndpointTypes),
		SupportedOperations:    uniqueCatalogStrings(item.SupportedOperations),
		SupportStatus:          normalizeCatalogSupportStatus(string(item.SupportStatus)),
		SupportReason:          strings.TrimSpace(item.SupportReason),
		CatalogSource:          strings.TrimSpace(item.CatalogSource),
		CatalogVersion:         strings.TrimSpace(item.CatalogVersion),
		DocumentationPaths:     uniqueCatalogStrings(item.DocumentationPaths),
		APIPath:                strings.TrimSpace(item.APIPath),
	}
}

func enrichCatalogItem(current ChannelModelCatalogItem, fallback ChannelModelCatalogItem) ChannelModelCatalogItem {
	if current.DisplayName == "" {
		current.DisplayName = fallback.DisplayName
	}
	if current.ProviderModelKey == "" {
		current.ProviderModelKey = fallback.ProviderModelKey
	}
	if current.ModelType == "" {
		current.ModelType = fallback.ModelType
	}
	if len(current.SupportedEndpointTypes) == 0 {
		current.SupportedEndpointTypes = fallback.SupportedEndpointTypes
	}
	if current.Protocol == "" {
		current.Protocol = fallback.Protocol
	}
	current.SupportedOperations = uniqueCatalogStrings(append(current.SupportedOperations, fallback.SupportedOperations...))
	current.DocumentationPaths = uniqueCatalogStrings(append(current.DocumentationPaths, fallback.DocumentationPaths...))
	if fallback.SupportStatus != "" {
		current.SupportStatus = fallback.SupportStatus
	}
	if fallback.SupportReason != "" {
		current.SupportReason = fallback.SupportReason
	}
	if fallback.CatalogSource != "" {
		if current.CatalogSource == "" || current.CatalogSource == fallback.CatalogSource {
			current.CatalogSource = fallback.CatalogSource
		} else {
			current.CatalogSource = current.CatalogSource + "+" + fallback.CatalogSource
		}
	}
	if fallback.CatalogVersion != "" {
		current.CatalogVersion = fallback.CatalogVersion
	}
	if fallback.APIPath != "" {
		current.APIPath = fallback.APIPath
	}
	return current
}

func normalizeCatalogSupportStatus(value string) model.ChannelModelSupportStatus {
	switch model.ChannelModelSupportStatus(strings.ToLower(strings.TrimSpace(value))) {
	case model.ChannelModelSupportReady:
		return model.ChannelModelSupportReady
	case model.ChannelModelSupportUnsupported:
		return model.ChannelModelSupportUnsupported
	case model.ChannelModelSupportDeprecated:
		return model.ChannelModelSupportDeprecated
	case model.ChannelModelSupportPlanned:
		return model.ChannelModelSupportPlanned
	default:
		return model.ChannelModelSupportUnknown
	}
}

func uniqueCatalogStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func outboundHeadersMap(headers []OutboundHeader) map[string]string {
	result := make(map[string]string, len(headers))
	for _, header := range headers {
		result[header.Name] = header.Value
	}
	return result
}

func modelCatalogRegion(baseURL string) string {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return ""
	}
	host := strings.ToLower(parsed.Hostname())
	if strings.Contains(host, "dashscope-us") {
		return "us-east-1"
	}
	if strings.Contains(host, "ap-southeast-1") {
		return "ap-southeast-1"
	}
	return "cn-beijing"
}
