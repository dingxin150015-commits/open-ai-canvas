package service

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"infinite-canvas/backend/internal/model"
	"infinite-canvas/backend/internal/protocol"
	"strings"
)

// Evaluate against a single live registry snapshot; discovery never grants execution.
func evaluateChannelModel(registry *protocol.Registry, item *model.ChannelModel) model.ChannelModelEligibility {
	result := model.ChannelModelEligibility{CanEditConfiguration: true}
	fail := func(code, reason string) model.ChannelModelEligibility {
		result.ReasonCode, result.Reason = code, reason
		return result
	}
	if item == nil {
		return fail("configuration_required", "请先配置模型")
	}
	if _, err := effectiveChannelModelCapability(item); err != nil {
		return fail("configuration_required", "请完善模型能力参数")
	}
	result.ConfigurationValid = true
	if item.SupportStatus == model.ChannelModelSupportUnsupported || item.SupportStatus == model.ChannelModelSupportDeprecated {
		return fail("model_unavailable", "模型已明确标记为不支持或已废弃，不能执行")
	}
	adapter, ok := registry.Resolve(string(item.Protocol))
	if !ok || !adapter.Metadata().Enabled || adapter.Metadata().UnavailableReason != "" {
		return fail("protocol_unavailable", "请求协议未安装、未启用或不可用，请选择有效协议")
	}
	metadata := adapter.Metadata()
	result.ProtocolVersion = metadata.Version
	if !containsPluginSurface(metadata.Scopes, protocol.SurfaceAdminSystemChannel) || !containsPluginCapability(metadata.Categories, protocol.Capability(item.Capability)) {
		return fail("protocol_incompatible", "协议不支持当前系统渠道能力")
	}
	key := firstNonEmpty(item.ProviderModelKey, item.ModelKey)
	if strings.TrimSpace(key) == "" {
		return fail("configuration_required", "请填写模型标识")
	}
	if err := validateBailianLegacyProtocol(key, string(item.Protocol)); err != nil {
		return fail("protocol_incompatible", err.Error())
	}
	if err := metadata.ModelContract.Check(key); err != nil {
		return fail("model_unsupported", err.Error())
	}
	// A documented missing adapter cannot be replaced by an implicit generic fallback.
	if item.SupportStatus == model.ChannelModelSupportPlanned && strings.Contains(item.CatalogSource, "official") && metadata.ModelContract == nil {
		return fail("model_contract_required", "该模型尚缺明确执行合同，请安装支持此模型的协议插件")
	}
	for _, tier := range item.PriceTiers {
		if !tier.Enabled {
			continue
		}
		if err := metadata.ModelContract.Check(firstNonEmpty(tier.ProviderModelKey, key)); err != nil {
			return fail("sku_unsupported", err.Error())
		}
	}
	result.ExecutionSupported = true
	if !HasValidPrice(item) {
		return fail("price_required", "请明确配置有效价格；默认零值不代表免费")
	}
	result.CanEnable = true
	result.CanExecute = item.Enabled
	if !item.Enabled {
		return fail("model_disabled", "模型尚未启用")
	}
	return result
}

// Inspect configuration without ID allocation, database writes or provider calls.
func (s *Service) InspectAdminChannelModel(actor *model.User, channelID string, req ChannelModelRequest) (*model.ChannelModelEligibility, error) {
	if err := s.RequireAdmin(actor); err != nil {
		return nil, err
	}
	channel, err := s.adminSystemChannel(channelID)
	if err != nil {
		return nil, err
	}
	req.Draft = false
	key, upstreamKey, capability, protocolID, err := s.normalizeChannelModelContract(channel, req)
	if err != nil {
		return nil, err
	}
	profile, err := NormalizeModelCapabilityConfigForModel(capability, string(protocolID), upstreamKey, req.CapabilityConfig)
	if err != nil {
		return nil, err
	}
	encoded, err := json.Marshal(profile)
	if err != nil {
		return nil, err
	}
	tiers, err := s.normalizeChannelModelPriceTiers(req, capability, protocolID, upstreamKey)
	if err != nil {
		return nil, err
	}
	candidate := model.ChannelModel{ChannelID: channelID, ModelKey: key, ProviderModelKey: upstreamKey, Capability: capability, Protocol: protocolID, CapabilityConfigJSON: string(encoded), PriceTiers: tiers, Enabled: req.Enabled != nil && *req.Enabled}
	existing, err := s.repo.ChannelModelByKeyIncludingDisabled(channelID, key)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		candidate.SupportStatus = existing.SupportStatus
		candidate.CatalogSource = existing.CatalogSource
	}
	if err := validateChannelModelTierCapabilities(tiers, candidate.CapabilityConfigJSON, capability); err != nil {
		return nil, err
	}
	result := s.channelModelEligibility(&candidate)
	return &result, nil
}

func (s *Service) channelModelEligibility(item *model.ChannelModel) model.ChannelModelEligibility {
	return evaluateChannelModel(s.protocolRegistry(), item)
}
func (s *Service) channelModelSupported(item *model.ChannelModel) bool {
	return s.channelModelEligibility(item).ExecutionSupported
}
func (s *Service) requireChannelModelSupported(item *model.ChannelModel) error {
	result := s.channelModelEligibility(item)
	if !result.ExecutionSupported {
		return BadAuthRequest(result.Reason)
	}
	return nil
}
