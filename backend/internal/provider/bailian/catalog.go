package bailian

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"

	"infinite-canvas/backend/internal/provider"
)

//go:embed catalog.generated.json
var catalogJSON []byte

type catalogManifest struct {
	Provider string                 `json:"provider"`
	Version  string                 `json:"version"`
	Models   []catalogManifestModel `json:"models"`
}

type catalogManifestModel struct {
	ID                     string                      `json:"id"`
	DisplayName            string                      `json:"displayName"`
	ProviderModelKey       string                      `json:"providerModelKey"`
	Capability             string                      `json:"capability"`
	Protocol               string                      `json:"protocol"`
	SupportedEndpointTypes []string                    `json:"supportedEndpointTypes"`
	APIPath                string                      `json:"apiPath"`
	SupportStatus          provider.ModelSupportStatus `json:"supportStatus"`
	SupportReason          string                      `json:"supportReason"`
	SupportedOperations    []string                    `json:"supportedOperations"`
	DocumentationPaths     []string                    `json:"documentationPaths"`
	CatalogSource          string                      `json:"catalogSource"`
	CatalogVersion         string                      `json:"catalogVersion"`
}

var (
	catalogOnce    sync.Once
	catalogItems   []provider.Model
	catalogVersion string
	catalogErr     error
)

func officialCatalogModels() ([]provider.Model, error) {
	catalogOnce.Do(func() {
		manifest := catalogManifest{}
		if err := json.Unmarshal(catalogJSON, &manifest); err != nil {
			catalogErr = fmt.Errorf("parse Bailian catalog manifest: %w", err)
			return
		}
		if strings.TrimSpace(manifest.Provider) != "bailian" || strings.TrimSpace(manifest.Version) == "" {
			catalogErr = fmt.Errorf("invalid Bailian catalog manifest header")
			return
		}
		seen := make(map[string]bool, len(manifest.Models))
		catalogItems = make([]provider.Model, 0, len(manifest.Models))
		for _, item := range manifest.Models {
			id := strings.TrimSpace(item.ID)
			if id == "" || seen[id] {
				catalogErr = fmt.Errorf("invalid or duplicate Bailian model ID %q", id)
				return
			}
			if !validSupportStatus(item.SupportStatus) {
				catalogErr = fmt.Errorf("invalid support status %q for Bailian model %s", item.SupportStatus, id)
				return
			}
			if strings.TrimSpace(item.Capability) == "" || strings.TrimSpace(item.Protocol) == "" || len(item.DocumentationPaths) == 0 {
				catalogErr = fmt.Errorf("incomplete Bailian model contract for %s", id)
				return
			}
			seen[id] = true
			catalogItems = append(catalogItems, provider.Model{
				ID:                     id,
				DisplayName:            strings.TrimSpace(item.DisplayName),
				ProviderModelKey:       strings.TrimSpace(item.ProviderModelKey),
				Provider:               "bailian",
				Capability:             []string{strings.TrimSpace(item.Capability)},
				Protocol:               strings.TrimSpace(item.Protocol),
				SupportedEndpointTypes: uniqueStrings(item.SupportedEndpointTypes),
				SupportedOperations:    uniqueStrings(item.SupportedOperations),
				APIPath:                strings.TrimSpace(item.APIPath),
				SupportStatus:          item.SupportStatus,
				SupportReason:          strings.TrimSpace(item.SupportReason),
				CatalogSource:          strings.TrimSpace(item.CatalogSource),
				CatalogVersion:         strings.TrimSpace(item.CatalogVersion),
				DocumentationPaths:     uniqueStrings(item.DocumentationPaths),
				Deprecated:             item.SupportStatus == provider.ModelSupportDeprecated,
				RequiresPlan:           item.SupportStatus == provider.ModelSupportPlanned,
			})
		}
		sort.Slice(catalogItems, func(left int, right int) bool { return catalogItems[left].ID < catalogItems[right].ID })
		catalogVersion = manifest.Version
	})
	if catalogErr != nil {
		return nil, catalogErr
	}
	result := make([]provider.Model, len(catalogItems))
	copy(result, catalogItems)
	return result, nil
}

func officialCatalogVersion() string {
	_, _ = officialCatalogModels()
	return catalogVersion
}

func validSupportStatus(status provider.ModelSupportStatus) bool {
	switch status {
	case provider.ModelSupportReady, provider.ModelSupportPlanned, provider.ModelSupportUnsupported, provider.ModelSupportDeprecated:
		return true
	default:
		return false
	}
}

func uniqueStrings(values []string) []string {
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
