package protocol

import (
	"fmt"
	"strings"
)

// ModelContract describes structural protocol support, never account entitlement.
// The same contract is checked during configuration and before serialization.
type ModelContract struct {
	Mode   string   `json:"mode"`
	Models []string `json:"models,omitempty"`
}

type ModelCatalog struct {
	Hosts   []string       `json:"hosts"`
	Entries []CatalogModel `json:"entries"`
}

type CatalogModel struct {
	ID                  string   `json:"id"`
	DisplayName         string   `json:"displayName,omitempty"`
	SupportStatus       string   `json:"supportStatus,omitempty"`
	SupportReason       string   `json:"supportReason,omitempty"`
	SupportedOperations []string `json:"supportedOperations,omitempty"`
	DocumentationPaths  []string `json:"documentationPaths,omitempty"`
	CatalogSource       string   `json:"catalogSource,omitempty"`
	CatalogVersion      string   `json:"catalogVersion,omitempty"`
	APIPath             string   `json:"apiPath,omitempty"`
}

func (c *ModelContract) Validate() error {
	if c == nil {
		return nil
	}
	if c.Mode != "compatible" && c.Mode != "exact" {
		return fmt.Errorf("invalid model contract mode %q", c.Mode)
	}
	if (c.Mode == "exact") != (len(c.Models) > 0) {
		return fmt.Errorf("only exact model contracts require models")
	}
	seen := map[string]bool{}
	for _, key := range c.Models {
		if key == "" || key != strings.TrimSpace(key) || seen[key] {
			return fmt.Errorf("invalid or duplicate model contract key")
		}
		seen[key] = true
	}
	return nil
}

func (c *ModelContract) Check(model string) error {
	if c == nil {
		return nil
	} // Older protocol packages retain their existing configurable contract.
	key := strings.TrimPrefix(strings.TrimSpace(model), "models/")
	if key == "" {
		return fmt.Errorf("模型标识不能为空")
	}
	if c.Mode == "compatible" {
		return nil
	}
	for _, supported := range c.Models {
		if key == supported {
			return nil
		}
	}
	return fmt.Errorf("当前协议插件尚不支持模型 %s", key)
}

func validateProviderModelMetadata(provider ManifestProvider, version string) error {
	if provider.ModelContract == nil && provider.ModelCatalog == nil {
		return nil
	}
	if version != "yingce.plugin/v4" {
		return fmt.Errorf("model contracts and catalogs require yingce.plugin/v4")
	}
	if err := provider.ModelContract.Validate(); err != nil {
		return err
	}
	if provider.ModelCatalog == nil {
		return nil
	}
	if len(provider.Capabilities) != 1 || len(provider.ModelCatalog.Hosts) == 0 || len(provider.ModelCatalog.Entries) > 2000 {
		return fmt.Errorf("invalid model catalog dimensions")
	}
	for _, host := range provider.ModelCatalog.Hosts {
		if host == "" || strings.ContainsAny(host, "/:@?#*") || host != strings.ToLower(strings.TrimSpace(host)) {
			return fmt.Errorf("invalid catalog host")
		}
	}
	seen := map[string]bool{}
	for _, entry := range provider.ModelCatalog.Entries {
		if entry.ID == "" || seen[entry.ID] {
			return fmt.Errorf("invalid catalog model identity")
		}
		seen[entry.ID] = true
		if entry.SupportStatus == "ready" && provider.ModelContract != nil && provider.ModelContract.Check(entry.ID) != nil {
			return fmt.Errorf("catalog ready model is outside execution contract")
		}
	}
	return nil
}
