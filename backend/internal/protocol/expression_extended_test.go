package protocol

import (
	"testing"
)

func TestStrictExpressionValues(t *testing.T) {
	for _, value := range []any{false, "false", "0"} {
		got, err := evaluateManifestValue(map[string]any{"$strictBool": map[string]any{"value": value, "default": true}}, nil)
		if err != nil || got != false {
			t.Fatalf("false lost: %v %v", got, err)
		}
	}
	for _, value := range []any{1.5, "many", true} {
		if _, err := evaluateManifestValue(map[string]any{"$strictInt": map[string]any{"value": value, "default": 5}}, nil); err == nil {
			t.Fatalf("invalid integer accepted: %v", value)
		}
	}
	got, err := evaluateManifestValue(map[string]any{"$strictInt": map[string]any{"value": "-1", "min": 2, "max": 30, "allow": []any{-1}}}, nil)
	if err != nil || got != int64(-1) {
		t.Fatal("smart duration lost")
	}
}

func TestAtomicPluginRegistration(t *testing.T) {
	info := Metadata{ID: "occupied", Version: "1"}
	registry, _ := NewRegistry(UnavailableAdapter{Info: info})
	newInfo := info
	newInfo.ID = "must-not-leak"
	if err := registry.RegisterAll([]Adapter{UnavailableAdapter{Info: newInfo}, UnavailableAdapter{Info: info}}); err == nil {
		t.Fatal("collision accepted")
	}
	if _, ok := registry.Get(newInfo.ID); ok {
		t.Fatal("partial plugin became executable")
	}
}

func TestPixelSizeRejectsNonFiniteAndOverflow(t *testing.T) {
	for _, value := range []string{"NaN:1", "+Inf:1", "1:0", "999999999999x999999999999", "0x1024"} {
		if _, err := resolveManifestPixelSize(map[string]any{"value": value, "minPixels": 262144, "maxPixels": 4194304, "minRatio": 0.125, "maxRatio": 8, "targetPixels": 1048576, "ratioMaxSide": 2048, "separator": "*"}); err == nil {
			t.Fatalf("bad size accepted: %s", value)
		}
	}
}

func TestFrameSelectionRejectsMissingDuplicateAndExtraMedia(t *testing.T) {
	items := []any{map[string]any{"id": "first"}, map[string]any{"id": "last"}}
	for _, ids := range [][]any{{"missing", "last"}, {"first", "first"}, {"first"}} {
		if _, err := evaluateManifestValue(map[string]any{"$selectByIds": map[string]any{"from": items, "ids": ids}}, nil); err == nil {
			t.Fatalf("invalid frame selection accepted: %v", ids)
		}
	}
}
