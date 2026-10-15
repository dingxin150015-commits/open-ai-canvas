// 渠道模型的上游目录拉取、导入预览与初始同步。

package app

import (
	"context"
	"log"
	"strings"

	"yingce/backend/internal/model"
	"yingce/backend/internal/repository"
)

func (s *Service) FetchAdminChannelModels(ctx context.Context, actor *model.User, channelID string) (*AdminChannelModelFetchResult, error) {
	if err := s.RequireAdmin(actor); err != nil {
		return nil, err
	}
	channel, err := s.adminSystemChannel(channelID)
	if err != nil {
		return nil, err
	}
	headers, err := ParseOutboundHeadersJSON(channel.HeadersJSON)
	if err != nil {
		return nil, err
	}
	// 使用完整 CatalogItem，避免模型标识列表丢掉官方能力、协议、操作与支持状态。
	catalog, err := s.FetchChannelModelCatalog(ctx, actor, ChannelModelsRequest{BaseURL: channel.BaseURL, AllowLocalChannel: channel.AllowLocalChannel, APIKey: channel.APIKey, APIFormat: channel.APIFormat, Headers: headers})
	if err != nil {
		return nil, err
	}
	// 只按当前未删除记录去重；普通手动删除的模型仍可重新拉取，已合并进模型家族的 SKU 除外。
	existing, err := s.repo.ChannelModels(channelID, true)
	if err != nil {
		return nil, err
	}
	existingByKey := make(map[string]model.ChannelModel, len(existing))
	for _, item := range existing {
		existingByKey[channelModelCatalogKey(item.ModelKey)] = item
	}
	retired := retiredChannelModelKeys(channel.RetiredModelsJSON)
	missing := make([]model.ChannelModel, 0, len(catalog))
	updates := make([]repository.ChannelModelCatalogUpdate, 0, len(catalog))
	result := adminChannelModelCatalogResult(catalog)
	for _, catalogItem := range catalog {
		name := strings.TrimPrefix(strings.TrimSpace(catalogItem.ID), "models/")
		if name == "" {
			continue
		}
		key := channelModelCatalogKey(name)
		if retired[key] {
			result.SkippedRetired++
			continue
		}
		desired := channelModelFromCatalog("", channelID, catalogItem)
		if current, ok := existingByKey[key]; ok {
			if !catalogMayBeEnriched(current) {
				result.SkippedConfigured++
			} else if update := catalogEnrichmentUpdate(current, desired); update != nil {
				for field := range update.Changes {
					result.UpdateFieldCounts[field]++
				}
				updates = append(updates, *update)
			} else {
				result.Unchanged++
			}
			continue
		}
		// 自动发现不能绕过支持状态或定价边界。
		modelID, idErr := s.repo.NextPrefixedID("MODEL")
		if idErr != nil {
			return nil, idErr
		}
		desired.ID = modelID
		missing = append(missing, desired)
	}
	result.Added, result.Updated, err = s.repo.SyncChannelModelCatalog(missing, updates)
	if err != nil {
		return nil, err
	}
	if result.Added > 0 || result.Updated > 0 {
		s.invalidateRouteCatalog()
	}
	if len(result.UpdateFieldCounts) > 0 {
		log.Printf("channel model catalog sync summary added=%d updated=%d update_fields=%v", result.Added, result.Updated, result.UpdateFieldCounts)
	}
	return result, nil
}

func adminChannelModelCatalogResult(catalog []ChannelModelCatalogItem) *AdminChannelModelFetchResult {
	result := &AdminChannelModelFetchResult{
		Models:              make([]string, 0, len(catalog)),
		CapabilityCounts:    map[string]int{},
		SupportStatusCounts: map[string]int{},
		UpdateFieldCounts:   map[string]int{},
	}
	for _, item := range catalog {
		name := strings.TrimPrefix(strings.TrimSpace(item.ID), "models/")
		if name == "" {
			continue
		}
		result.Models = append(result.Models, name)
		result.CapabilityCounts[normalizeCatalogModelType(item.ModelType)]++
		result.SupportStatusCounts[string(normalizeCatalogSupportStatus(string(item.SupportStatus)))]++
		if strings.Contains(item.CatalogSource, "upstream") {
			result.UpstreamCount++
		}
		if strings.Contains(item.CatalogSource, "official") {
			result.SupplementalCount++
			result.OfficialCatalogReady = true
		}
	}
	return result
}

// PreviewAdminChannelModels 只读取上游模型目录，不修改渠道模型配置。
func (s *Service) PreviewAdminChannelModels(ctx context.Context, actor *model.User, channelID string) ([]string, error) {
	if err := s.RequireAdmin(actor); err != nil {
		return nil, err
	}
	catalog, err := s.fetchAdminChannelModelCatalog(ctx, actor, channelID)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(catalog))
	for _, item := range catalog {
		names = append(names, strings.TrimPrefix(strings.TrimSpace(item.ID), "models/"))
	}
	return uniqueNonEmpty(names), nil
}

// ImportAdminChannelModels 只导入管理员明确选择、且仍存在于上游目录中的模型。
func (s *Service) ImportAdminChannelModels(ctx context.Context, actor *model.User, channelID string, selected []string) (*AdminChannelModelFetchResult, error) {
	if err := s.RequireAdmin(actor); err != nil {
		return nil, err
	}
	channel, err := s.adminSystemChannel(channelID)
	if err != nil {
		return nil, err
	}
	models, err := s.fetchAdminChannelModelCatalog(ctx, actor, channelID)
	if err != nil {
		return nil, err
	}
	if len(selected) == 0 {
		return nil, BadAuthRequest("请至少选择一个要导入的模型")
	}
	if len(selected) > 500 {
		return nil, BadAuthRequest("单次最多导入 500 个模型")
	}
	available := make(map[string]ChannelModelCatalogItem, len(models))
	for _, item := range models {
		available[channelModelCatalogKey(item.ID)] = item
	}
	chosenCatalog := make([]ChannelModelCatalogItem, 0, len(selected))
	seen := make(map[string]struct{}, len(selected))
	for _, rawName := range selected {
		name := strings.TrimPrefix(strings.TrimSpace(rawName), "models/")
		key := channelModelCatalogKey(name)
		if key == "" {
			continue
		}
		canonical, ok := available[key]
		if !ok {
			return nil, BadAuthRequest("所选模型不在上游模型目录中：" + name)
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		chosenCatalog = append(chosenCatalog, canonical)
	}
	if len(chosenCatalog) == 0 {
		return nil, BadAuthRequest("请至少选择一个有效的模型")
	}

	existing, err := s.repo.ChannelModels(channelID, true)
	if err != nil {
		return nil, err
	}
	existingByKey := make(map[string]model.ChannelModel, len(existing))
	for _, item := range existing {
		existingByKey[channelModelCatalogKey(item.ModelKey)] = item
	}
	retired := retiredChannelModelKeys(channel.RetiredModelsJSON)
	missing := make([]model.ChannelModel, 0, len(chosenCatalog))
	updates := make([]repository.ChannelModelCatalogUpdate, 0, len(chosenCatalog))
	result := adminChannelModelCatalogResult(chosenCatalog)
	for _, catalogItem := range chosenCatalog {
		name := strings.TrimPrefix(strings.TrimSpace(catalogItem.ID), "models/")
		key := channelModelCatalogKey(name)
		if retired[key] {
			result.SkippedRetired++
			continue
		}
		desired := channelModelFromCatalog("", channelID, catalogItem)
		if current, ok := existingByKey[key]; ok {
			if !catalogMayBeEnriched(current) {
				result.SkippedConfigured++
			} else if update := catalogEnrichmentUpdate(current, desired); update != nil {
				for field := range update.Changes {
					result.UpdateFieldCounts[field]++
				}
				updates = append(updates, *update)
			} else {
				result.Unchanged++
			}
			continue
		}
		modelID, idErr := s.repo.NextPrefixedID("MODEL")
		if idErr != nil {
			return nil, idErr
		}
		desired.ID = modelID
		missing = append(missing, desired)
	}
	result.Added, result.Updated, err = s.repo.SyncChannelModelCatalog(missing, updates)
	if err != nil {
		return nil, err
	}
	if result.Added > 0 || result.Updated > 0 {
		s.invalidateRouteCatalog()
	}
	if len(result.UpdateFieldCounts) > 0 {
		log.Printf("selected channel model catalog sync summary added=%d updated=%d update_fields=%v", result.Added, result.Updated, result.UpdateFieldCounts)
	}
	return result, nil
}

func (s *Service) fetchAdminChannelModelCatalog(ctx context.Context, actor *model.User, channelID string) ([]ChannelModelCatalogItem, error) {
	channel, err := s.adminSystemChannel(channelID)
	if err != nil {
		return nil, err
	}
	headers, err := ParseOutboundHeadersJSON(channel.HeadersJSON)
	if err != nil {
		return nil, err
	}
	models, err := s.FetchChannelModelCatalog(ctx, actor, ChannelModelsRequest{BaseURL: channel.BaseURL, AllowLocalChannel: channel.AllowLocalChannel, APIKey: channel.APIKey, APIFormat: channel.APIFormat, Headers: headers})
	if err != nil {
		return nil, err
	}
	return models, nil
}
