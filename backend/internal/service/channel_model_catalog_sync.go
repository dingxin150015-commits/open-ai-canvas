package service

import (
	"encoding/json"
	"strings"

	"infinite-canvas/backend/internal/model"
	"infinite-canvas/backend/internal/repository"
)

func marshalCatalogStrings(values []string) string {
	values = uniqueCatalogStrings(values)
	encoded, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(encoded)
}

func unmarshalCatalogStrings(raw string) []string {
	var values []string
	if json.Unmarshal([]byte(strings.TrimSpace(raw)), &values) != nil {
		return []string{}
	}
	return uniqueCatalogStrings(values)
}

func channelModelFromCatalog(id string, channelID string, item ChannelModelCatalogItem) model.ChannelModel {
	modelKey := strings.TrimPrefix(strings.TrimSpace(item.ID), "models/")
	displayName := strings.TrimSpace(item.DisplayName)
	if displayName == "" {
		displayName = modelKey
	}
	result := model.ChannelModel{
		ID:                      id,
		ChannelID:               channelID,
		ModelKey:                modelKey,
		ProviderModelKey:        strings.TrimPrefix(strings.TrimSpace(firstNonEmpty(item.ProviderModelKey, modelKey)), "models/"),
		DisplayName:             displayName,
		Capability:              normalizeCatalogModelType(item.ModelType),
		Protocol:                model.ChannelInterfaceType(strings.TrimSpace(item.Protocol)),
		SupportStatus:           normalizeCatalogSupportStatus(string(item.SupportStatus)),
		SupportReason:           strings.TrimSpace(item.SupportReason),
		CatalogSource:           strings.TrimSpace(item.CatalogSource),
		CatalogVersion:          strings.TrimSpace(item.CatalogVersion),
		SupportedOperationsJSON: marshalCatalogStrings(item.SupportedOperations),
		DocumentationPathsJSON:  marshalCatalogStrings(item.DocumentationPaths),
		BillingMode:             "fixed_request",
		Enabled:                 false,
		PriceConfigured:         false,
		PriceVersion:            1,
	}
	if result.SupportStatus == model.ChannelModelSupportReady && (result.Capability == "text" || result.Capability == "image" || result.Capability == "video") {
		capabilityConfig := DefaultModelCapabilityConfigForModel(string(result.Protocol), result.ProviderModelKey)
		if normalized, err := NormalizeModelCapabilityConfig(result.Capability, string(result.Protocol), capabilityConfig); err == nil && normalized != nil {
			if encoded, err := json.Marshal(normalized); err == nil {
				result.CapabilityConfigJSON = string(encoded)
				result.CapabilityVersion = 1
			}
		}
	}
	return result
}

func catalogMayBeEnriched(current model.ChannelModel) bool {
	if current.Enabled || current.PriceConfigured || strings.TrimSpace(current.CapabilityConfigJSON) != "" {
		return false
	}
	for _, tier := range current.PriceTiers {
		if tier.PriceConfigured {
			return false
		}
	}
	return true
}

func catalogEnrichmentUpdate(current model.ChannelModel, desired model.ChannelModel) *repository.ChannelModelCatalogUpdate {
	if !catalogMayBeEnriched(current) {
		return nil
	}
	changes := map[string]any{}
	if strings.TrimSpace(current.ProviderModelKey) == "" {
		changes["provider_model_key"] = desired.ProviderModelKey
	}
	if (strings.TrimSpace(current.DisplayName) == "" || current.DisplayName == current.ModelKey) && current.DisplayName != desired.DisplayName {
		changes["display_name"] = desired.DisplayName
	}
	if strings.TrimSpace(current.Capability) == "" && strings.TrimSpace(desired.Capability) != "" {
		changes["capability"] = desired.Capability
	}
	if strings.TrimSpace(string(current.Protocol)) == "" && strings.TrimSpace(string(desired.Protocol)) != "" {
		changes["protocol"] = desired.Protocol
	}
	if strings.TrimSpace(current.CapabilityConfigJSON) == "" && strings.TrimSpace(desired.CapabilityConfigJSON) != "" {
		changes["capability_config_json"] = desired.CapabilityConfigJSON
		changes["capability_version"] = desired.CapabilityVersion
	}
	if current.SupportStatus != desired.SupportStatus {
		changes["support_status"] = desired.SupportStatus
	}
	if current.SupportReason != desired.SupportReason {
		changes["support_reason"] = desired.SupportReason
	}
	if current.CatalogSource != desired.CatalogSource {
		changes["catalog_source"] = desired.CatalogSource
	}
	if current.CatalogVersion != desired.CatalogVersion {
		changes["catalog_version"] = desired.CatalogVersion
	}
	if current.SupportedOperationsJSON != desired.SupportedOperationsJSON {
		changes["supported_operations_json"] = desired.SupportedOperationsJSON
	}
	if current.DocumentationPathsJSON != desired.DocumentationPathsJSON {
		changes["documentation_paths_json"] = desired.DocumentationPathsJSON
	}
	if len(changes) == 0 {
		return nil
	}
	return &repository.ChannelModelCatalogUpdate{ID: current.ID, ExpectedUpdatedAt: current.UpdatedAt, Changes: changes}
}

func populateChannelModelCatalogMetadata(item *model.ChannelModel) {
	item.SupportedOperations = unmarshalCatalogStrings(item.SupportedOperationsJSON)
	item.DocumentationPaths = unmarshalCatalogStrings(item.DocumentationPathsJSON)
}

func requireReadyChannelModel(item *model.ChannelModel) error {
	if item == nil || item.SupportStatus == model.ChannelModelSupportReady {
		return nil
	}
	return BadAuthRequest("该模型当前为" + channelModelSupportStatusLabel(item.SupportStatus) + "状态，仅可查看；完成执行器适配和测试后才能定价、启用或测试")
}

func channelModelSupportStatusLabel(status model.ChannelModelSupportStatus) string {
	switch status {
	case model.ChannelModelSupportUnsupported:
		return "不支持"
	case model.ChannelModelSupportDeprecated:
		return "已废弃"
	case model.ChannelModelSupportReady:
		return "可用"
	default:
		return "计划支持"
	}
}
