package service

import (
	"infinite-canvas/backend/internal/model"
	"testing"
)

func TestHasValidPriceRequiresConfiguredFlag(t *testing.T) {
	scalar := &model.ChannelModel{
		BillingMode:           "fixed_request",
		UnitPriceMicrocredits: 100,
		PriceConfigured:       false,
	}
	if HasValidPrice(scalar) {
		t.Fatal("draft scalar price must not be routable")
	}
	scalar.PriceConfigured = true
	if !HasValidPrice(scalar) {
		t.Fatal("configured scalar price should be routable")
	}

	tiered := &model.ChannelModel{PriceTiers: []model.ChannelModelPriceTier{{
		BillingMode:           "fixed_request",
		UnitPriceMicrocredits: 100,
		Enabled:               true,
		PriceConfigured:       false,
	}}}
	if HasValidPrice(tiered) {
		t.Fatal("draft tier price must not be routable")
	}
	tiered.PriceTiers[0].PriceConfigured = true
	if !HasValidPrice(tiered) {
		t.Fatal("configured tier price should be routable")
	}
}
