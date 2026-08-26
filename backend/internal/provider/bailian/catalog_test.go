package bailian

import (
	"testing"

	"infinite-canvas/backend/internal/provider"
)

func TestOfficialCatalogManifest(t *testing.T) {
	items, err := officialCatalogModels()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 80 {
		t.Fatalf("catalog items = %d, want 80", len(items))
	}
	seen := map[string]provider.Model{}
	for _, item := range items {
		if _, exists := seen[item.ID]; exists {
			t.Fatalf("duplicate model ID %q", item.ID)
		}
		wantStatus := provider.ModelSupportPlanned
		readyModels := map[string]bool{
			"wan3.0-video": true,
			"wan2.7-t2v":   true, "wan2.7-t2v-2026-06-12": true, "wan2.7-t2v-2026-04-25": true,
			"wan2.7-i2v": true, "wan2.7-i2v-2026-04-25": true,
			"wan2.7-r2v": true, "wan2.7-r2v-2026-06-12": true,
			"happyhorse-1.1-t2v": true, "happyhorse-1.1-i2v": true, "happyhorse-1.1-r2v": true,
		}
		if readyModels[item.ID] {
			wantStatus = provider.ModelSupportReady
		}
		if item.SupportStatus != wantStatus {
			t.Fatalf("model %q status = %q, want %q", item.ID, item.SupportStatus, wantStatus)
		}
		if item.Protocol == "" || len(item.DocumentationPaths) == 0 {
			t.Fatalf("incomplete contract for %q", item.ID)
		}
		seen[item.ID] = item
	}
	for _, id := range []string{"wan3.0-video", "wan3.0-video-prime", "wan2.7-t2v", "happyhorse-1.1-t2v", "qwen-image-3.0-pro"} {
		if _, exists := seen[id]; !exists {
			t.Fatalf("official catalog missing %q", id)
		}
	}
	if len(seen["wan3.0-video"].SupportedOperations) != 5 {
		t.Fatalf("wan3.0-video operations = %#v", seen["wan3.0-video"].SupportedOperations)
	}
	if officialCatalogVersion() != "2026-08-25" {
		t.Fatalf("catalog version = %q", officialCatalogVersion())
	}
}

func TestDiscoveryMatchesBailianAndReturnsOfficialCatalog(t *testing.T) {
	discovery := New()
	if !discovery.Match("https://dashscope.aliyuncs.com/compatible-mode/v1", nil) {
		t.Fatal("DashScope URL did not match Bailian discovery")
	}
	if discovery.Match("https://api.openai.com/v1", nil) {
		t.Fatal("unrelated URL matched Bailian discovery")
	}
	if got := len(discovery.AdditionalModels(provider.DiscoveryConfig{})); got != 80 {
		t.Fatalf("AdditionalModels() count = %d, want 80", got)
	}
}
