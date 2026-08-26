package database

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"infinite-canvas/backend/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestAssetIDColumnsUseSharedLimit(t *testing.T) {
	tests := []struct {
		value any
		field string
	}{
		{value: &model.Asset{}, field: "ID"},
		{value: &model.ProjectAssetLink{}, field: "AssetID"},
		{value: &model.ProjectAssetCandidate{}, field: "ResolvedAssetID"},
		{value: &model.AssetVersion{}, field: "AssetID"},
	}

	for _, test := range tests {
		parsed, err := schema.Parse(test.value, &sync.Map{}, schema.NamingStrategy{})
		if err != nil {
			t.Fatalf("parse schema: %v", err)
		}
		field := parsed.LookUpField(test.field)
		if field == nil {
			t.Fatalf("field %s not found in %s", test.field, parsed.Table)
		}
		if field.Size != model.AssetIDMaxLength {
			t.Fatalf("%s.%s size = %d, want %d", parsed.Table, field.DBName, field.Size, model.AssetIDMaxLength)
		}
	}
}

func TestPostgresAssetIDMigrationsCoverEveryAssetIDColumn(t *testing.T) {
	want := map[string]bool{
		"assets.id":                                  false,
		"project_asset_links.asset_id":               false,
		"project_asset_candidates.resolved_asset_id": false,
		"asset_versions.asset_id":                    false,
	}
	for _, migration := range assetIDColumnMigrations {
		key := migration.table + "." + migration.column
		if _, exists := want[key]; !exists {
			t.Fatalf("unexpected asset ID migration %s", key)
		}
		if !strings.Contains(migration.statement, fmt.Sprintf("varchar(%d)", model.AssetIDMaxLength)) {
			t.Fatalf("asset ID migration %s does not use limit %d", key, model.AssetIDMaxLength)
		}
		want[key] = true
	}
	for key, covered := range want {
		if !covered {
			t.Fatalf("missing asset ID migration %s", key)
		}
	}
}

func TestBackfillChannelModelCatalogState(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.ChannelModel{}); err != nil {
		t.Fatal(err)
	}
	configured := model.ChannelModel{ID: "configured", ChannelID: "channel", ModelKey: "configured-model", Capability: "text", Protocol: model.ChannelInterfaceChatCompletion, CapabilityConfigJSON: `{"version":1,"text":{"references":{"promptMaxChars":32000,"maxImages":0,"maxImageBytes":0,"maxVideos":0,"maxVideoBytes":0,"maxVideoDurationSeconds":0}}}`, BillingMode: "fixed_request", PriceVersion: 1}
	incomplete := model.ChannelModel{ID: "incomplete", ChannelID: "channel", ModelKey: "incomplete-model", Capability: "text", Protocol: model.ChannelInterfaceChatCompletion, BillingMode: "fixed_request", PriceVersion: 1}
	legacy := model.ChannelModel{ID: "legacy", ChannelID: "channel", ModelKey: "legacy-model", BillingMode: "fixed_request", PriceVersion: 1}
	if err := db.Create(&configured).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&incomplete).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&legacy).Error; err != nil {
		t.Fatal(err)
	}
	if err := backfillChannelModelCatalogState(db); err != nil {
		t.Fatal(err)
	}
	if err := db.First(&configured, "id = ?", configured.ID).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.First(&legacy, "id = ?", legacy.ID).Error; err != nil {
		t.Fatal(err)
	}
	if err := db.First(&incomplete, "id = ?", incomplete.ID).Error; err != nil {
		t.Fatal(err)
	}
	if configured.SupportStatus != model.ChannelModelSupportReady || configured.CatalogSource != "manual" {
		t.Fatalf("configured backfill = %#v", configured)
	}
	if legacy.SupportStatus != model.ChannelModelSupportPlanned || legacy.CatalogSource != "legacy" || legacy.ProviderModelKey != legacy.ModelKey {
		t.Fatalf("legacy backfill = %#v", legacy)
	}
	if incomplete.SupportStatus != model.ChannelModelSupportPlanned || incomplete.CatalogSource != "legacy" {
		t.Fatalf("incomplete configured backfill = %#v", incomplete)
	}
}
