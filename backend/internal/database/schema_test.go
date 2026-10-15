package database

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"yingce/backend/internal/model"

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

func TestModelsIncludeCloudAgentCanvasMutation(t *testing.T) {
	for _, item := range Models() {
		if _, ok := item.(*model.CloudAgentCanvasMutation); ok {
			return
		}
	}
	t.Fatal("Models must include CloudAgentCanvasMutation")
}

func TestBackfillChannelModelCatalogState(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:catalog-backfill?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.ChannelModel{}); err != nil {
		t.Fatal(err)
	}
	configured := model.ChannelModel{ID: "configured", ChannelID: "channel", ModelKey: "configured-model", Capability: "text", Protocol: model.ChannelInterfaceChatCompletion, CapabilityConfigJSON: `{"version":1,"text":{"references":{"promptMaxChars":32000,"maxImages":0,"maxImageBytes":0,"maxVideos":0,"maxVideoBytes":0,"maxVideoDurationSeconds":0}}}`, BillingMode: "fixed_request", PriceVersion: 1}
	incomplete := model.ChannelModel{ID: "incomplete", ChannelID: "channel", ModelKey: "incomplete-model", Capability: "text", Protocol: model.ChannelInterfaceChatCompletion, BillingMode: "fixed_request", PriceVersion: 1}
	legacy := model.ChannelModel{ID: "legacy", ChannelID: "channel", ModelKey: "legacy-model", BillingMode: "fixed_request", PriceVersion: 1}
	for _, item := range []*model.ChannelModel{&configured, &incomplete, &legacy} {
		if err := db.Create(item).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := backfillChannelModelCatalogState(db); err != nil {
		t.Fatal(err)
	}
	for _, item := range []*model.ChannelModel{&configured, &incomplete, &legacy} {
		if err := db.First(item, "id = ?", item.ID).Error; err != nil {
			t.Fatal(err)
		}
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

func TestMigrationPlanAcceptsProviderDurationVersionNineLineage(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:duration-v9-lineage?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&schemaMigration{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&schemaMigration{Version: 9, Name: "provider_duration_metering", Checksum: providerDurationMeteringChecksum}).Error; err != nil {
		t.Fatal(err)
	}
	plan, err := migrationsForDatabase(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan) != int(CurrentSchemaVersion) {
		t.Fatalf("plan length = %d, want %d", len(plan), CurrentSchemaVersion)
	}
	if plan[8].version != 9 || plan[8].name != "provider_duration_metering" {
		t.Fatalf("migration 9 = %#v", plan[8])
	}
	if plan[9].version != 10 || plan[9].name != "channel_presentation" {
		t.Fatalf("migration 10 = %#v", plan[9])
	}
	if plan[30].version != 31 || plan[30].name != "builtin_tools" {
		t.Fatalf("legacy migration 31 = %#v", plan[30])
	}
	if plan[40].version != 41 || plan[40].name != "builtin_skill_tombstones" {
		t.Fatalf("legacy migration 41 = %#v", plan[40])
	}
	if plan[41].version != 42 || plan[41].name != "comfy_bridge_runtime" || plan[41].checksum != "sha256:comfy-bridge-runtime-v42-20260927" {
		t.Fatalf("legacy migration 42 = %#v", plan[41])
	}
	for index, item := range plan {
		if item.version != int64(index+1) {
			t.Fatalf("legacy migration gap at %d: %#v", index+1, item)
		}
	}
	last := plan[len(plan)-1]
	if last.version != CurrentSchemaVersion || last.name != "upload_reservations" {
		t.Fatalf("last migration = %#v", last)
	}
}

func TestMigrationPlanKeepsOfficialVersionNineLineage(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:official-v9-lineage?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&schemaMigration{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&schemaMigration{Version: 9, Name: "channel_presentation", Checksum: "sha256:channel-presentation-v9-20260908"}).Error; err != nil {
		t.Fatal(err)
	}
	plan, err := migrationsForDatabase(db)
	if err != nil {
		t.Fatal(err)
	}
	if plan[8].version != 9 || plan[8].name != "channel_presentation" {
		t.Fatalf("migration 9 = %#v", plan[8])
	}
	if plan[29].version != 30 || plan[29].name != "builtin_tools" {
		t.Fatalf("official migration 30 = %#v", plan[29])
	}
	if plan[40].version != 41 || plan[40].name != "resource_thumbnail" || plan[40].checksum != resourceThumbnailChecksum {
		t.Fatalf("official migration 41 = %#v", plan[40])
	}
	if plan[43].version != 44 || plan[43].name != "provider_duration_metering" {
		t.Fatalf("official migration 44 = %#v", plan[43])
	}
	last := plan[len(plan)-1]
	if last.version != CurrentSchemaVersion || last.name != "comfy_bridge_runtime" {
		t.Fatalf("last migration = %#v", last)
	}
}

func TestMigrationPlanPreservesV158LocalVersionFortyOneLineage(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:v158-local-lineage?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&schemaMigration{}); err != nil {
		t.Fatal(err)
	}
	for _, record := range []schemaMigration{
		{Version: 9, Name: "channel_presentation", Checksum: "sha256:channel-presentation-v9-20260908"},
		{Version: 41, Name: "provider_duration_metering", Checksum: providerDurationMeteringV41Checksum},
		{Version: 42, Name: "comfy_bridge_runtime", Checksum: "sha256:comfy-bridge-runtime-v42-20260927"},
	} {
		if err := db.Create(&record).Error; err != nil {
			t.Fatal(err)
		}
	}
	plan, err := migrationsForDatabase(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan) != int(CurrentSchemaVersion) {
		t.Fatalf("plan length = %d, want %d", len(plan), CurrentSchemaVersion)
	}
	for index, item := range plan {
		if item.version != int64(index+1) {
			t.Fatalf("migration gap at %d: %#v", index+1, item)
		}
	}
	for index, name := range []string{"provider_duration_metering", "comfy_bridge_runtime", "resource_thumbnail", "cloud_agent_pi_sessions", "topup_sale_strategies"} {
		if plan[40+index].name != name {
			t.Fatalf("migration %d = %#v, want %s", 41+index, plan[40+index], name)
		}
	}
	if last := plan[len(plan)-1]; last.version != CurrentSchemaVersion || last.name != "upload_reservations" {
		t.Fatalf("last legacy migration = %#v", last)
	}
}

func TestMigrationPlanPreservesV160LocalVersionFortyFiveLineage(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:v160-local45-lineage?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&schemaMigration{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&schemaMigration{Version: 45, Name: "comfy_bridge_runtime", Checksum: "sha256:comfy-bridge-runtime-v45-20261003"}).Error; err != nil {
		t.Fatal(err)
	}
	plan, err := migrationsForDatabase(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan) != int(CurrentSchemaVersion) || plan[44].name != "comfy_bridge_runtime" || plan[50].name != "upload_reservations" {
		t.Fatalf("local45 plan = %#v", plan)
	}
	for index, item := range plan {
		if item.version != int64(index+1) {
			t.Fatalf("local45 migration gap at %d: %#v", index+1, item)
		}
	}
}

func TestMigrationPlanKeepsOfficialVersionFortyFiveLineage(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file:official45-lineage?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&schemaMigration{}); err != nil {
		t.Fatal(err)
	}
	if err := db.Create(&schemaMigration{Version: 45, Name: "upload_reservations", Checksum: "sha256:upload-reservations-v45-20261005"}).Error; err != nil {
		t.Fatal(err)
	}
	plan, err := migrationsForDatabase(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan) != int(CurrentSchemaVersion) || plan[43].name != "provider_duration_metering" || plan[44].name != "upload_reservations" || plan[50].name != "comfy_bridge_runtime" {
		t.Fatalf("official45 plan = %#v", plan)
	}
}
