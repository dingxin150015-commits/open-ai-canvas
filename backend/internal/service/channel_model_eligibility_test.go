package service

import (
	"encoding/json"
	"infinite-canvas/backend/internal/model"
	"infinite-canvas/backend/internal/protocol"
	"testing"
)

func eligibilityFixture(t *testing.T, id, protocolID, capability string) model.ChannelModel {
	t.Helper()
	profile, err := json.Marshal(DefaultModelCapabilityConfigForModel(protocolID, id))
	if err != nil {
		t.Fatal(err)
	}
	return model.ChannelModel{ModelKey: id, ProviderModelKey: id, Capability: capability, Protocol: model.ChannelInterfaceType(protocolID), SupportStatus: model.ChannelModelSupportPlanned, CatalogSource: "upstream", CapabilityConfigJSON: string(profile), PriceConfigured: true, BillingMode: "fixed_request", Enabled: true}
}

func completeModelFixtureProfile(t *testing.T, item *model.ChannelModel) {
	t.Helper()
	encoded, err := json.Marshal(DefaultModelCapabilityConfigForModel(string(item.Protocol), firstNonEmpty(item.ProviderModelKey, item.ModelKey)))
	if err != nil {
		t.Fatal(err)
	}
	item.CapabilityConfigJSON = string(encoded)
}

func TestChannelEligibilityUsesConfigurationNotDiscoveryStatus(t *testing.T) {
	registry := loadOfficialFallbackRegistry()
	candidate := eligibilityFixture(t, "new-compatible-model", "chat-completion", "text")
	imported := evaluateChannelModel(registry, &candidate)
	candidate.SupportStatus = model.ChannelModelSupportReady
	manual := evaluateChannelModel(registry, &candidate)
	if imported != manual || !imported.CanExecute {
		t.Fatalf("import/manual divergence: %#v %#v", imported, manual)
	}
	candidate.Enabled = false
	stopped := evaluateChannelModel(registry, &candidate)
	if !stopped.CanEnable || stopped.CanExecute {
		t.Fatal("stopped model eligibility invalid")
	}
	candidate.PriceConfigured = false
	if evaluateChannelModel(registry, &candidate).CanEnable {
		t.Fatal("implicit zero cannot enable")
	}
}

func TestChannelEligibilityRejectsMissingPluginsAndWrongExactModels(t *testing.T) {
	empty, _ := protocol.NewRegistry()
	candidate := eligibilityFixture(t, "qwen-image-3.0-pro", "dashscope-image", "image")
	if evaluateChannelModel(empty, &candidate).ExecutionSupported {
		t.Fatal("missing plugin allowed")
	}
	registry := loadOfficialFallbackRegistry()
	if !evaluateChannelModel(registry, &candidate).ExecutionSupported {
		t.Fatal("supported Qwen contract unavailable")
	}
	candidate.ProviderModelKey = "wan2.7-image-pro"
	if evaluateChannelModel(registry, &candidate).ExecutionSupported {
		t.Fatal("unimplemented model allowed")
	}
	candidate.ProviderModelKey = "qwen-image-3.0-pro"
	candidate.SupportStatus = model.ChannelModelSupportDeprecated
	if evaluateChannelModel(registry, &candidate).ExecutionSupported {
		t.Fatal("deprecated model allowed")
	}
}

func TestCatalogContributionUsesTheLivePluginRegistry(t *testing.T) {
	empty, _ := protocol.NewRegistry()
	if got := extendChannelModelCatalogWithRegistry(empty, "https://dashscope.aliyuncs.com", nil); len(got) != 0 {
		t.Fatal("disabled plugin still contributes")
	}
	if got := extendChannelModelCatalogWithRegistry(loadOfficialFallbackRegistry(), "https://dashscope.aliyuncs.com", nil); len(got) != 80 {
		t.Fatalf("expected 80 catalog facts, got %d", len(got))
	}
	if got := extendChannelModelCatalogWithRegistry(loadOfficialFallbackRegistry(), "https://dashscope.aliyuncs.com.evil.invalid", nil); len(got) != 0 {
		t.Fatal("untrusted host matched")
	}
}

func TestImportedModelCanBeConfiguredInPlace(t *testing.T) {
	svc, db := newChannelModelTestService(t)
	actor := &model.User{ID: "admin", Role: model.UserRoleAdmin}
	channel := model.ModelChannel{ID: "config-channel", UserID: actor.ID, Scope: model.ChannelScopeSystem, Enabled: true, Name: "fixture", BaseURL: "https://example.com/v1", APIKey: "fixture", APIFormat: "openai", ModelsJSON: `[]`}
	if err := db.Create(&channel).Error; err != nil {
		t.Fatal(err)
	}
	existing := model.ChannelModel{ID: "imported-model", ChannelID: channel.ID, ModelKey: "new-model", SupportStatus: model.ChannelModelSupportPlanned, CatalogSource: "upstream", Enabled: false}
	if err := db.Create(&existing).Error; err != nil {
		t.Fatal(err)
	}
	enabled := true
	saved, err := svc.SaveAdminChannelModel(actor, channel.ID, existing.ID, ChannelModelRequest{ModelKey: existing.ModelKey, Capability: "text", Protocol: "chat-completion", CapabilityConfig: DefaultModelCapabilityConfigForModel("chat-completion", existing.ModelKey), PriceTiers: []ChannelModelPriceTierRequest{{BillingMode: "fixed_request", PriceConfigured: true}}, Enabled: &enabled})
	if err != nil {
		t.Fatal(err)
	}
	if saved.ID != existing.ID || !saved.Eligibility.CanExecute {
		t.Fatalf("in-place configuration failed %#v", saved)
	}
	if _, err := svc.SystemChannelModel(channel.ID, existing.ModelKey); err != nil {
		t.Fatal("configured imported model blocked by repository", err)
	}
	public, err := svc.publicSystemChannelCatalog(nil)
	if err != nil || len(public) != 1 || len(public[0].Models) != 1 {
		t.Fatalf("configured import disappeared from public catalog: %#v %v", public, err)
	}
}

func TestIncompleteDraftCannotEnableAndPreservesIdentity(t *testing.T) {
	svc, db := newChannelModelTestService(t)
	actor := &model.User{ID: "admin", Role: model.UserRoleAdmin}
	channel := model.ModelChannel{ID: "draft-channel", UserID: actor.ID, Scope: model.ChannelScopeSystem, Enabled: true, Name: "fixture", BaseURL: "https://example.com", APIKey: "fixture", ModelsJSON: `[]`}
	if err := db.Create(&channel).Error; err != nil {
		t.Fatal(err)
	}
	enabled := true
	saved, err := svc.SaveAdminChannelModel(actor, channel.ID, "", ChannelModelRequest{Draft: true, ModelKey: "unclassified", Enabled: &enabled, BillingMode: "fixed_request"})
	if err != nil {
		t.Fatal(err)
	}
	if saved.Enabled || saved.Eligibility.CanExecute || !saved.Eligibility.CanEditConfiguration {
		t.Fatal("draft gained execution")
	}
}

func TestConfigurationInspectionNeverWritesOrAllocatesIDs(t *testing.T) {
	svc, db := newChannelModelTestService(t)
	actor := &model.User{ID: "admin", Role: model.UserRoleAdmin}
	channel := model.ModelChannel{ID: "inspection-channel", UserID: actor.ID, Scope: model.ChannelScopeSystem, Enabled: true, Name: "fixture", BaseURL: "https://must-not-be-contacted.invalid", APIKey: "fixture", ModelsJSON: `[]`}
	if err := db.Create(&channel).Error; err != nil {
		t.Fatal(err)
	}
	var before, after int64
	if err := db.Model(&model.IDSequence{}).Count(&before).Error; err != nil {
		t.Fatal(err)
	}
	enabled := true
	result, err := svc.InspectAdminChannelModel(actor, channel.ID, ChannelModelRequest{ModelKey: "new-model", Capability: "text", Protocol: "chat-completion", CapabilityConfig: DefaultModelCapabilityConfigForModel("chat-completion", "new-model"), PriceTiers: []ChannelModelPriceTierRequest{{BillingMode: "fixed_request", PriceConfigured: true}}, Enabled: &enabled})
	if err != nil || !result.CanEnable {
		t.Fatalf("invalid inspection %#v %v", result, err)
	}
	db.Model(&model.IDSequence{}).Count(&after)
	if before != after {
		t.Fatal("inspection allocated IDs")
	}
	var count int64
	db.Model(&model.ChannelModel{}).Count(&count)
	if count != 0 {
		t.Fatal("inspection saved a model")
	}
	if _, err := svc.InspectAdminChannelModel(&model.User{ID: "ordinary", Role: model.UserRoleUser}, channel.ID, ChannelModelRequest{}); err == nil {
		t.Fatal("ordinary user allowed to inspect admin configuration")
	}
}

func TestChannelEligibilityRechecksPluginLifecycleWithoutChangingSavedModel(t *testing.T) {
	adapter, ok := loadOfficialFallbackRegistry().Resolve("chat-completion")
	if !ok {
		t.Fatal("missing fixture protocol")
	}
	registry, err := protocol.NewRegistry(adapter)
	if err != nil {
		t.Fatal(err)
	}
	candidate := eligibilityFixture(t, "new-id", "chat-completion", "text")
	before, _ := json.Marshal(candidate)
	if !evaluateChannelModel(registry, &candidate).CanExecute {
		t.Fatal("configured protocol rejected")
	}
	registry.Unregister("chat-completion")
	if evaluateChannelModel(registry, &candidate).CanExecute {
		t.Fatal("removed protocol remained executable")
	}
	if err := registry.Register(adapter); err != nil {
		t.Fatal(err)
	}
	if !evaluateChannelModel(registry, &candidate).CanExecute {
		t.Fatal("restored protocol remained unavailable")
	}
	after, _ := json.Marshal(candidate)
	if string(before) != string(after) {
		t.Fatal("runtime eligibility rewrote persisted configuration")
	}
}
