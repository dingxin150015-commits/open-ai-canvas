package service

import (
	"testing"

	"infinite-canvas/backend/internal/model"
)

func TestHasValidPriceRequiresConfiguredFlagAndAllowsExplicitFree(t *testing.T) {
	scalar := &model.ChannelModel{BillingMode: "fixed_request", UnitPriceMicrocredits: 0}
	if HasValidPrice(scalar) {
		t.Fatal("draft scalar price must not be routable")
	}
	scalar.PriceConfigured = true
	if !HasValidPrice(scalar) {
		t.Fatal("explicitly configured free scalar price should be routable")
	}
	scalar.UnitPriceMicrocredits = -1
	if HasValidPrice(scalar) {
		t.Fatal("negative scalar price must not be routable")
	}

	tiered := &model.ChannelModel{PriceTiers: []model.ChannelModelPriceTier{{
		BillingMode: "fixed_request", UnitPriceMicrocredits: 0, Enabled: true,
	}}}
	if HasValidPrice(tiered) {
		t.Fatal("draft tier price must not be routable")
	}
	tiered.PriceTiers[0].PriceConfigured = true
	if !HasValidPrice(tiered) {
		t.Fatal("explicitly configured free tier price should be routable")
	}
}

func TestValidateChannelModelPriceSupportsArkVideoTokensOnly(t *testing.T) {
	const outputPrice = int64(16_000_000)
	if !ValidateChannelModelPrice("token", "video", model.ChannelInterfaceVolcengineArkVideo, 0, 0, outputPrice, 0) {
		t.Fatal("Volcengine Ark video Token price should be valid")
	}
	for _, protocol := range []model.ChannelInterfaceType{model.ChannelInterfaceVolcengineJiMengVideo, model.ChannelInterfaceNewAPIVideo} {
		if ValidateChannelModelPrice("token", "video", protocol, 0, 0, outputPrice, 0) {
			t.Fatalf("video protocol %q should not support Token pricing", protocol)
		}
	}
}

func TestHasValidPriceUsesChannelProtocolForTokenTiers(t *testing.T) {
	tier := model.ChannelModelPriceTier{Enabled: true, PriceConfigured: true, BillingMode: "token", OutputTokenPriceMicrocredits: 16_000_000}
	ark := &model.ChannelModel{Capability: "video", Protocol: model.ChannelInterfaceVolcengineArkVideo, PriceTiers: []model.ChannelModelPriceTier{tier}}
	if !HasValidPrice(ark) {
		t.Fatal("Volcengine Ark video Token tier should be valid")
	}

	jimeng := &model.ChannelModel{Capability: "video", Protocol: model.ChannelInterfaceVolcengineJiMengVideo, PriceTiers: []model.ChannelModelPriceTier{tier}}
	if HasValidPrice(jimeng) {
		t.Fatal("JiMeng video Token tier should be invalid")
	}
}

func TestHasValidPriceAllowsExplicitFreeTextTokens(t *testing.T) {
	text := &model.ChannelModel{Capability: "text", Protocol: model.ChannelInterfaceChatCompletion, PriceTiers: []model.ChannelModelPriceTier{{Enabled: true, PriceConfigured: true, BillingMode: "token"}}}
	if !HasValidPrice(text) {
		t.Fatal("explicitly configured free text Token tier should be valid")
	}
	image := &model.ChannelModel{Capability: "image", Protocol: model.ChannelInterfaceOpenAIImage, PriceTiers: []model.ChannelModelPriceTier{{Enabled: true, PriceConfigured: true, BillingMode: "token"}}}
	if HasValidPrice(image) {
		t.Fatal("image Token tier should remain unsupported")
	}
}
