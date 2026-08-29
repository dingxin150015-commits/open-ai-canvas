package service

import (
	"encoding/json"
	"testing"

	"infinite-canvas/backend/internal/model"
	"infinite-canvas/backend/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPublicSystemChannelCatalogExcludesDisabledAndMarksUnpricedUnavailable(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file:"+newID()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(&model.ModelChannel{}, &model.ChannelModel{}, &model.ChannelModelPriceTier{}); err != nil {
		t.Fatal(err)
	}

	channel := model.ModelChannel{
		ID:        "channel-1",
		Scope:     model.ChannelScopeSystem,
		Enabled:   true,
		Name:      "Bailian",
		APIFormat: "openai",
	}
	if err := database.Create(&channel).Error; err != nil {
		t.Fatal(err)
	}

	profile, err := json.Marshal(DefaultModelCapabilityConfigForModel(string(model.ChannelInterfaceDashScopeVideo), "wan3.0-video"))
	if err != nil {
		t.Fatal(err)
	}
	models := []model.ChannelModel{
		{
			ID:                   "wan30-disabled",
			ChannelID:            channel.ID,
			ModelKey:             "wan3.0-video",
			ProviderModelKey:     "wan3.0-video",
			DisplayName:          "Wan 3.0",
			Capability:           "video",
			Protocol:             model.ChannelInterfaceDashScopeVideo,
			SupportStatus:        model.ChannelModelSupportReady,
			CapabilityConfigJSON: string(profile),
			Enabled:              false,
			PriceConfigured:      false,
		},
		{
			ID:                   "unpriced-enabled",
			ChannelID:            channel.ID,
			ModelKey:             "unpriced-video",
			ProviderModelKey:     "wan3.0-video",
			DisplayName:          "Unpriced Video",
			Capability:           "video",
			Protocol:             model.ChannelInterfaceDashScopeVideo,
			SupportStatus:        model.ChannelModelSupportReady,
			CapabilityConfigJSON: string(profile),
			Enabled:              true,
			PriceConfigured:      false,
		},
		{
			ID:                    "planned-priced",
			ChannelID:             channel.ID,
			ModelKey:              "planned-video",
			ProviderModelKey:      "planned-video",
			DisplayName:           "Planned Video",
			Capability:            "video",
			Protocol:              model.ChannelInterfaceDashScopeVideo,
			SupportStatus:         model.ChannelModelSupportPlanned,
			CapabilityConfigJSON:  string(profile),
			Enabled:               true,
			PriceConfigured:       true,
			BillingMode:           "fixed_request",
			UnitPriceMicrocredits: 0,
		},
	}
	if err := database.Create(&models).Error; err != nil {
		t.Fatal(err)
	}

	service := &Service{repo: repository.New(database)}
	catalog, err := service.publicSystemChannelCatalog(nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(catalog) != 1 || len(catalog[0].Models) != 1 {
		t.Fatalf("catalog = %#v, want one channel with one enabled model", catalog)
	}
	got := catalog[0].Models[0]
	if got.ModelKey != "unpriced-video" {
		t.Fatalf("model = %q, disabled and non-ready models must be omitted", got.ModelKey)
	}
	if got.Available {
		t.Fatal("enabled but unpriced model must remain unavailable to the creation UI")
	}
}
