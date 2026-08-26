package bailian

import (
	"log"
	"net/url"
	"strings"

	"infinite-canvas/backend/internal/provider"
)

// Discovery 阿里云百炼官方模型目录扩展。
type Discovery struct{}

func New() *Discovery { return &Discovery{} }

func (d *Discovery) GetProviderID() string { return "bailian" }

func (d *Discovery) Match(baseURL string, _ map[string]string) bool {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	dashscopeHost := host == "dashscope.aliyuncs.com" || strings.HasSuffix(host, ".dashscope.aliyuncs.com") || (strings.HasPrefix(host, "dashscope-") && strings.HasSuffix(host, ".aliyuncs.com"))
	return dashscopeHost || host == "maas.aliyuncs.com" || strings.HasSuffix(host, ".maas.aliyuncs.com")
}

// AdditionalModels 返回标准 /models 没有暴露的官方模型目录项。
// Manifest 只决定发现与支持状态，不代表允许真实调用。
func (d *Discovery) AdditionalModels(_ provider.DiscoveryConfig) []provider.Model {
	models, err := officialCatalogModels()
	if err != nil {
		log.Printf("[BailianCatalog] official manifest unavailable: %v", err)
		return nil
	}
	return models
}

func (d *Discovery) GetMetadata() provider.ProviderMetadata {
	return provider.ProviderMetadata{
		Name:             "bailian",
		Version:          officialCatalogVersion(),
		Description:      "Alibaba Cloud Bailian official versioned model catalog",
		Author:           "yingce",
		SupportedRegions: []string{"cn-beijing", "ap-southeast-1", "us-east-1"},
	}
}
