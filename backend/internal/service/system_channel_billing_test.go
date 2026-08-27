package service

import (
	"encoding/json"
	"testing"

	"infinite-canvas/backend/internal/model"
	"infinite-canvas/backend/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSystemChannelTaskBillingUsesExactRequestedPriceTier(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file:"+newID()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if sqlDB, openErr := database.DB(); openErr == nil {
		sqlDB.SetMaxOpenConns(1)
	}
	if err := database.AutoMigrate(
		&model.SystemSetting{},
		&model.AdminAuditEvent{},
		&model.ModelChannel{},
		&model.ChannelModel{},
		&model.ChannelModelPriceTier{},
		&model.IDSequence{},
		&model.BillingOrder{},
	); err != nil {
		t.Fatal(err)
	}

	channel := model.ModelChannel{ID: "channel-1", Scope: model.ChannelScopeSystem, Enabled: true, Name: "Bailian"}
	channelModel := model.ChannelModel{
		ID: "model-1", ChannelID: channel.ID, ModelKey: "wan3.0-video", ProviderModelKey: "wan3.0-video",
		Capability: "video", Protocol: model.ChannelInterfaceDashScopeVideo, SupportStatus: model.ChannelModelSupportReady,
		Enabled: true, PriceConfigured: true,
	}
	if err := database.Create(&channel).Error; err != nil {
		t.Fatal(err)
	}
	if err := database.Create(&channelModel).Error; err != nil {
		t.Fatal(err)
	}
	tiers := []model.ChannelModelPriceTier{
		{
			ID: "tier-480-2", ChannelModelID: channelModel.ID,
			SelectorKey:  `{"operation":"text_to_video","videoSeconds":"2","vquality":"480p"}`,
			SelectorJSON: `{"operation":"text_to_video","videoSeconds":"2","vquality":"480p"}`,
			Resolution:   "480p", VideoSeconds: 2, BillingMode: "per_second", UnitPriceMicrocredits: 1_000_000,
			PriceConfigured: true, Enabled: true, PriceVersion: 1,
		},
		{
			ID: "tier-1080-5", ChannelModelID: channelModel.ID,
			SelectorKey:  `{"operation":"text_to_video","videoSeconds":"5","vquality":"1080p"}`,
			SelectorJSON: `{"operation":"text_to_video","videoSeconds":"5","vquality":"1080p"}`,
			Resolution:   "1080p", VideoSeconds: 5, BillingMode: "per_second", UnitPriceMicrocredits: 9_000_000,
			PriceConfigured: true, Enabled: true, PriceVersion: 1,
		},
	}
	if err := database.Create(&tiers).Error; err != nil {
		t.Fatal(err)
	}

	service := New(repository.New(database), t.TempDir())
	input := map[string]any{
		"mode": "video",
		"config": map[string]any{
			"channelId": channel.ID, "model": channelModel.ModelKey,
			"vquality": "480p", "videoSeconds": "2", "size": "16:9",
		},
		"metadata": map[string]any{"videoEditOperation": "text_to_video"},
	}
	order, err := service.taskBillingOrder("user-1", &model.Task{ID: "task-1", Type: "canvas_video", Operation: "text_to_video"}, input)
	if err != nil {
		t.Fatal(err)
	}
	if order == nil || order.PriceTierID != "tier-480-2" || order.Quantity != 2 || order.AmountMicrocredits != 2_000_000 {
		t.Fatalf("billing order = %#v", order)
	}
}

func TestSystemChannelQuoteUsesExactRequestedPriceTierWithoutPersistence(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file:"+newID()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if sqlDB, openErr := database.DB(); openErr == nil {
		sqlDB.SetMaxOpenConns(1)
	}
	if err := database.AutoMigrate(
		&model.SystemSetting{},
		&model.AdminAuditEvent{},
		&model.ModelChannel{},
		&model.ChannelModel{},
		&model.ChannelModelPriceTier{},
		&model.IDSequence{},
		&model.BillingOrder{},
	); err != nil {
		t.Fatal(err)
	}
	profile, err := json.Marshal(DefaultModelCapabilityConfigForModel(string(model.ChannelInterfaceDashScopeVideo), "wan3.0-video"))
	if err != nil {
		t.Fatal(err)
	}
	channel := model.ModelChannel{ID: "channel-quote", Scope: model.ChannelScopeSystem, Enabled: true, Name: "Bailian"}
	channelModel := model.ChannelModel{
		ID: "model-quote", ChannelID: channel.ID, ModelKey: "wan3.0-video", ProviderModelKey: "wan3.0-video",
		Capability: "video", Protocol: model.ChannelInterfaceDashScopeVideo, SupportStatus: model.ChannelModelSupportReady,
		CapabilityConfigJSON: string(profile), Enabled: true, PriceConfigured: true,
	}
	if err := database.Create(&channel).Error; err != nil {
		t.Fatal(err)
	}
	if err := database.Create(&channelModel).Error; err != nil {
		t.Fatal(err)
	}
	tiers := []model.ChannelModelPriceTier{
		{
			ID: "tier-quote-480-2", ChannelModelID: channelModel.ID,
			SelectorKey: `{"operation":"text_to_video","videoSeconds":"2","vquality":"480p"}`, SelectorJSON: `{"operation":"text_to_video","videoSeconds":"2","vquality":"480p"}`,
			Resolution: "480p", VideoSeconds: 2, BillingMode: "per_second", UnitPriceMicrocredits: 1_000_000,
			PriceConfigured: true, Enabled: true, PriceVersion: 1,
		},
		{
			ID: "tier-quote-1080-5", ChannelModelID: channelModel.ID,
			SelectorKey: `{"operation":"text_to_video","videoSeconds":"5","vquality":"1080p"}`, SelectorJSON: `{"operation":"text_to_video","videoSeconds":"5","vquality":"1080p"}`,
			Resolution: "1080p", VideoSeconds: 5, BillingMode: "per_second", UnitPriceMicrocredits: 9_000_000,
			PriceConfigured: true, Enabled: true, PriceVersion: 1,
		},
	}
	if err := database.Create(&tiers).Error; err != nil {
		t.Fatal(err)
	}

	service := New(repository.New(database), t.TempDir())
	quote, err := service.QuoteSystemChannelModel(channelModel.ID, ModelRequestIntent{
		Capability: "video", Operation: "text_to_video",
		Inputs: map[string]int{}, Options: map[string]any{"vquality": "480p", "videoSeconds": 2, "aspectRatio": "16:9"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if quote.ModelID != channelModel.ID || quote.BillingMode != "per_second" || quote.Quantity != 2 || quote.AmountMicrocredits != 2_000_000 || quote.Estimated {
		t.Fatalf("quote = %#v", quote)
	}
	var billingOrders int64
	if err := database.Model(&model.BillingOrder{}).Count(&billingOrders).Error; err != nil || billingOrders != 0 {
		t.Fatalf("quote must not persist billing orders, count=%d err=%v", billingOrders, err)
	}
}

func TestSystemChannelQuoteRejectsUnsupportedOrUnpricedSelection(t *testing.T) {
	database, err := gorm.Open(sqlite.Open("file:"+newID()+"?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := database.AutoMigrate(&model.SystemSetting{}, &model.AdminAuditEvent{}, &model.ModelChannel{}, &model.ChannelModel{}, &model.ChannelModelPriceTier{}, &model.IDSequence{}); err != nil {
		t.Fatal(err)
	}
	profile, err := json.Marshal(DefaultModelCapabilityConfigForModel(string(model.ChannelInterfaceDashScopeVideo), "wan3.0-video"))
	if err != nil {
		t.Fatal(err)
	}
	channel := model.ModelChannel{ID: "channel-reject", Scope: model.ChannelScopeSystem, Enabled: true, Name: "Bailian"}
	channelModel := model.ChannelModel{
		ID: "model-reject", ChannelID: channel.ID, ModelKey: "wan3.0-video", ProviderModelKey: "wan3.0-video",
		Capability: "video", Protocol: model.ChannelInterfaceDashScopeVideo, SupportStatus: model.ChannelModelSupportReady,
		CapabilityConfigJSON: string(profile), Enabled: true, PriceConfigured: true,
	}
	if err := database.Create(&channel).Error; err != nil {
		t.Fatal(err)
	}
	if err := database.Create(&channelModel).Error; err != nil {
		t.Fatal(err)
	}
	service := New(repository.New(database), t.TempDir())
	if _, err := service.QuoteSystemChannelModel(channelModel.ID, ModelRequestIntent{Capability: "image"}); err == nil {
		t.Fatal("capability mismatch must be rejected")
	}
	if _, err := service.QuoteSystemChannelModel(channelModel.ID, ModelRequestIntent{Capability: "video", Operation: "text_to_video", Options: map[string]any{"vquality": "480p", "videoSeconds": 2}}); err == nil {
		t.Fatal("missing exact price tier must be rejected")
	}
}
