package service

import (
	"fmt"
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
