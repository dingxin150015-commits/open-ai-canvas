package service

import (
	"fmt"
	"infinite-canvas/backend/internal/model"
)

// Old Wanx/Wan plugins are different contracts, not aliases. Other gateways
// may expose these model keys legitimately through their own protocols.
func validateBailianLegacyProtocol(key, id string) error {
	if id != "dashscope-wanx-image" && id != "dashscope-wan-video" {
		return nil
	}
	if isQwenImage30Model(key) || isWan27ReadyVideoModel(key) || isHappyHorse11ReadyVideoModel(key) || isWan30VideoModel(key) {
		return fmt.Errorf("该模型不能使用旧 Wanx/Wan 协议，请选择对应的 DashScope 专属插件")
	}
	return nil
}

func (s *Service) bailianModelProtocolAvailable(cm model.ChannelModel) bool {
	if validateBailianLegacyProtocol(firstNonEmpty(cm.ProviderModelKey, cm.ModelKey), string(cm.Protocol)) != nil {
		return false
	}
	if cm.Protocol == model.ChannelInterfaceDashScopeImage || cm.Protocol == model.ChannelInterfaceDashScopeVideo {
		return s.protocolIsSelectable(string(cm.Protocol))
	}
	return true
}
