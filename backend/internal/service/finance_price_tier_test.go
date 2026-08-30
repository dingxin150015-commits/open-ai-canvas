package service

import (
	"testing"

	"infinite-canvas/backend/internal/model"
	"infinite-canvas/backend/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestTaskBillingOrderMatchesSystemImagePriceTierFromRequestedSpec(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:"+newID()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.SystemSetting{}, &model.ChannelModel{}, &model.ChannelModelPriceTier{}); err != nil {
		t.Fatal(err)
	}
	channelModel := model.ChannelModel{
		ID: "channel-model-1", ChannelID: "channel-1", ModelKey: "gpt-image-2", ProviderModelKey: "gpt-image-2",
		Capability: "image", Protocol: model.ChannelInterfaceOpenAIImage, BillingMode: "fixed_request",
		PriceConfigured: true, Enabled: true, PriceVersion: 1,
	}
	if err := db.Create(&channelModel).Error; err != nil {
		t.Fatal(err)
	}
	defaultTier := model.ChannelModelPriceTier{
		ID: "tier-default", ChannelModelID: channelModel.ID, SelectorKey: `{}`, SelectorJSON: `{}`,
		Resolution: "*", ProviderModelKey: channelModel.ProviderModelKey, BillingMode: "fixed_request",
		UnitPriceMicrocredits: 1_000_000, PriceConfigured: true, Enabled: true, PriceVersion: 1,
	}
	tier := model.ChannelModelPriceTier{
		ID: "tier-2k", ChannelModelID: channelModel.ID, SelectorKey: `{"quality":"2k"}`, SelectorJSON: `{"quality":"2k"}`,
		Resolution: "*", ProviderModelKey: channelModel.ProviderModelKey, BillingMode: "fixed_request",
		UnitPriceMicrocredits: 4_000_000, PriceConfigured: true, Enabled: true, PriceVersion: 1,
	}
	if err := db.Create(&[]model.ChannelModelPriceTier{defaultTier, tier}).Error; err != nil {
		t.Fatal(err)
	}

	svc := New(repository.New(db), t.TempDir())
	order, err := svc.taskBillingOrder("user-1", &model.Task{ID: "task-1", Type: "canvas_image", Operation: "image"}, map[string]any{
		"mode": "image",
		"config": map[string]any{
			"channelId":   "channel-1",
			"model":       "gpt-image-2",
			"quality":     "2K",
			"size":        "*",
			"priceTierId": defaultTier.ID,
		},
	})
	if err != nil {
		t.Fatalf("taskBillingOrder() error = %v", err)
	}
	if order == nil || order.PriceTierID != tier.ID || order.AmountMicrocredits != tier.UnitPriceMicrocredits {
		t.Fatalf("taskBillingOrder() = %#v, want tier %s", order, tier.ID)
	}
}

func TestChannelModelPriceTierForBillingRejectsStaleExplicitTier(t *testing.T) {
	defaultTier := model.ChannelModelPriceTier{ID: "tier-default", SelectorJSON: `{}`, BillingMode: "fixed_request", PriceConfigured: true, Enabled: true}
	exactTier := model.ChannelModelPriceTier{ID: "tier-2k", SelectorJSON: `{"quality":"2k"}`, BillingMode: "fixed_request", PriceConfigured: true, Enabled: true}
	channelModel := model.ChannelModel{Capability: "image", Protocol: model.ChannelInterfaceOpenAIImage, PriceTiers: []model.ChannelModelPriceTier{defaultTier, exactTier}}
	intent := ModelRequestIntent{Capability: "image", Options: map[string]any{"quality": "2K"}}

	if tier := channelModelPriceTierForBilling(channelModel, defaultTier.ID, "image", intent); tier != nil {
		t.Fatalf("stale explicit tier = %#v, want nil", tier)
	}
	if tier := channelModelPriceTierForBilling(channelModel, exactTier.ID, "image", intent); tier == nil || tier.ID != exactTier.ID {
		t.Fatalf("matching explicit tier = %#v, want %s", tier, exactTier.ID)
	}
}
