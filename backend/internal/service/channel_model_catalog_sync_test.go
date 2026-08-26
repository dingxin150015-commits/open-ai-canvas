package service

import (
	"fmt"
	"testing"
	"time"

	"infinite-canvas/backend/internal/model"
)

func TestBailianOfficialCatalogExtendsSyntheticUpstreamFixture(t *testing.T) {
	upstream := make([]ChannelModelCatalogItem, 240)
	for index := range upstream {
		upstream[index] = ChannelModelCatalogItem{
			ID:            fmt.Sprintf("text-model-%03d", index),
			ModelType:     "text",
			SupportStatus: model.ChannelModelSupportPlanned,
			CatalogSource: "upstream",
		}
	}
	merged := extendChannelModelCatalog("https://dashscope.aliyuncs.com/compatible-mode/v1", "openai", nil, upstream)
	if len(merged) != 320 {
		t.Fatalf("merged catalog count = %d, want 320", len(merged))
	}
	videoCount := 0
	wan3Found := false
	for _, item := range merged {
		if item.ModelType == "video" {
			videoCount++
		}
		if item.ID == "wan3.0-video" {
			wan3Found = item.Protocol == string(model.ChannelInterfaceDashScopeVideo) && item.SupportStatus == model.ChannelModelSupportReady
		}
	}
	if videoCount == 0 || !wan3Found {
		t.Fatalf("videoCount=%d wan3Found=%v", videoCount, wan3Found)
	}
}

func TestCatalogEnrichmentPreservesConfiguredModels(t *testing.T) {
	desired := channelModelFromCatalog("MODEL-new", "channel-1", ChannelModelCatalogItem{
		ID: "wan3.0-video", DisplayName: "Wan 3.0", ProviderModelKey: "wan3.0-video", ModelType: "video",
		Protocol: string(model.ChannelInterfaceDashScopeVideo), SupportStatus: model.ChannelModelSupportPlanned,
		CatalogSource: "bailian-official", CatalogVersion: "2026-08-25", SupportedOperations: []string{"text_to_video"}, DocumentationPaths: []string{"video.md"},
	})
	untouched := model.ChannelModel{ID: "MODEL-old", ModelKey: desired.ModelKey, DisplayName: desired.ModelKey, UpdatedAt: time.Now()}
	update := catalogEnrichmentUpdate(untouched, desired)
	if update == nil || update.Changes["capability"] != "video" || update.Changes["support_status"] != model.ChannelModelSupportPlanned {
		t.Fatalf("catalog enrichment = %#v", update)
	}
	configured := untouched
	configured.Enabled = true
	configured.DisplayName = "管理员名称"
	configured.CapabilityConfigJSON = `{"video":{}}`
	if update := catalogEnrichmentUpdate(configured, desired); update != nil {
		t.Fatalf("configured model must not be enriched: %#v", update)
	}
	unchanged := desired
	unchanged.ID = "MODEL-existing"
	unchanged.UpdatedAt = time.Now()
	if update := catalogEnrichmentUpdate(unchanged, desired); update != nil {
		t.Fatalf("identical second pull must be idempotent: %#v", update)
	}
}

func TestUpstreamCatalogEnrichmentIsIdempotent(t *testing.T) {
	desired := channelModelFromCatalog("MODEL-upstream", "channel-1", ChannelModelCatalogItem{
		ID: "text-model", ProviderModelKey: "text-model", SupportStatus: model.ChannelModelSupportPlanned,
		SupportReason: "上游目录未提供经过项目验证的执行器合同", CatalogSource: "upstream", CatalogVersion: "upstream",
	})
	current := desired
	current.UpdatedAt = time.Now()
	if update := catalogEnrichmentUpdate(current, desired); update != nil {
		t.Fatalf("identical upstream catalog item must not be enriched: %#v", update)
	}
}

func TestPlannedChannelModelIsReadOnly(t *testing.T) {
	err := requireReadyChannelModel(&model.ChannelModel{SupportStatus: model.ChannelModelSupportPlanned})
	if err == nil {
		t.Fatal("planned model should be rejected")
	}
	if err := requireReadyChannelModel(&model.ChannelModel{SupportStatus: model.ChannelModelSupportReady}); err != nil {
		t.Fatalf("ready model rejected: %v", err)
	}
}

func TestReadyWan30CatalogItemCarriesCapabilityConfig(t *testing.T) {
	item := channelModelFromCatalog("MODEL-wan30", "channel-1", ChannelModelCatalogItem{
		ID: "wan3.0-video", ProviderModelKey: "wan3.0-video", ModelType: "video", Protocol: string(model.ChannelInterfaceDashScopeVideo),
		SupportStatus: model.ChannelModelSupportReady, CatalogSource: "bailian-official-docs", CatalogVersion: "2026-08-25",
	})
	if item.CapabilityVersion != 1 || item.CapabilityConfigJSON == "" {
		t.Fatalf("ready catalog item = %#v", item)
	}
	config, err := DecodeModelCapabilityConfig(item.CapabilityConfigJSON)
	if err != nil || config == nil || config.Video == nil || config.Video.Duration.Max != 30 || !config.Video.Duration.SmartSupported {
		t.Fatalf("capability config = %#v, error = %v", config, err)
	}
	current := model.ChannelModel{ID: "MODEL-existing", ChannelID: "channel-1", ModelKey: "wan3.0-video", UpdatedAt: time.Now()}
	update := catalogEnrichmentUpdate(current, item)
	if update == nil || update.Changes["capability_config_json"] == nil || update.Changes["support_status"] != model.ChannelModelSupportReady {
		t.Fatalf("catalog update = %#v", update)
	}
}

func TestReadyDashScopeVideoCatalogItemsCarryModelSpecificCapabilities(t *testing.T) {
	models := []string{
		"wan3.0-video",
		"wan2.7-t2v", "wan2.7-t2v-2026-06-12", "wan2.7-t2v-2026-04-25",
		"wan2.7-i2v", "wan2.7-i2v-2026-04-25",
		"wan2.7-r2v", "wan2.7-r2v-2026-06-12",
		"happyhorse-1.1-t2v", "happyhorse-1.1-i2v", "happyhorse-1.1-r2v",
	}
	for _, modelKey := range models {
		t.Run(modelKey, func(t *testing.T) {
			item := channelModelFromCatalog("MODEL-"+modelKey, "channel-1", ChannelModelCatalogItem{
				ID: modelKey, ProviderModelKey: modelKey, ModelType: "video", Protocol: string(model.ChannelInterfaceDashScopeVideo), SupportStatus: model.ChannelModelSupportReady,
			})
			config, err := DecodeModelCapabilityConfig(item.CapabilityConfigJSON)
			if err != nil || item.CapabilityVersion != 1 || config == nil || config.Video == nil {
				t.Fatalf("item=%#v config=%#v error=%v", item, config, err)
			}
			if _, err := NormalizeModelCapabilityConfig("video", string(model.ChannelInterfaceDashScopeVideo), config); err != nil {
				t.Fatalf("NormalizeModelCapabilityConfig() error = %v", err)
			}
		})
	}
	planned := channelModelFromCatalog("MODEL-prime", "channel-1", ChannelModelCatalogItem{
		ID: "wan3.0-video-prime", ProviderModelKey: "wan3.0-video-prime", ModelType: "video", Protocol: string(model.ChannelInterfaceDashScopeVideo), SupportStatus: model.ChannelModelSupportPlanned,
	})
	if planned.CapabilityConfigJSON != "" || planned.CapabilityVersion != 0 {
		t.Fatalf("planned model received executable capability: %#v", planned)
	}
}
