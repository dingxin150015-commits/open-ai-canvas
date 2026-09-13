package protocol

import "testing"

func TestModelContractModes(t *testing.T) {
	exact := &ModelContract{Mode: "exact", Models: []string{"model-a"}}
	if exact.Validate() != nil || exact.Check("models/model-a") != nil || exact.Check("model-b") == nil {
		t.Fatal("exact contract did not preserve model identity")
	}
	compatible := &ModelContract{Mode: "compatible"}
	if compatible.Validate() != nil || compatible.Check("new-model") != nil || compatible.Check("") == nil {
		t.Fatal("invalid compatible contract")
	}
	for _, invalid := range []*ModelContract{{Mode: "exact"}, {Mode: "compatible", Models: []string{"x"}}, {Mode: "unknown"}, {Mode: "exact", Models: []string{"x", "x"}}} {
		if invalid.Validate() == nil {
			t.Fatalf("accepted invalid contract %#v", invalid)
		}
	}
}

func TestModelMetadataRequiresNewHostContract(t *testing.T) {
	p := ManifestProvider{ModelContract: &ModelContract{Mode: "compatible"}}
	if validateProviderModelMetadata(p, "yingce.plugin/v3") == nil {
		t.Fatal("old host contract must reject model metadata")
	}
	if err := validateProviderModelMetadata(p, "yingce.plugin/v4"); err != nil {
		t.Fatal(err)
	}
}
