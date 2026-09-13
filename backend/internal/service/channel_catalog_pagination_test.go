package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGeminiCatalogPaginationKeepsCredentialsOffURL(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://example.test/v1beta/models", nil)
	req.Header.Set("x-goog-api-key", "fixture")
	calls := 0
	items, err := fetchChannelCatalogPages(req, "gemini", func(r *http.Request) ([]byte, string, error) {
		calls++
		if r.URL.Host != "example.test" || r.Header.Get("x-goog-api-key") != "fixture" || r.URL.Query().Get("key") != "" {
			t.Fatal("changed outbound authority")
		}
		if calls == 1 {
			return []byte(`{"models":[{"name":"models/a","displayName":"A","supportedGenerationMethods":["generateContent"]}],"nextPageToken":"cursor/+?="}`), "", nil
		}
		if r.URL.Query().Get("pageToken") != "cursor/+?=" {
			t.Fatal("cursor escaped incorrectly")
		}
		return []byte(`{"models":[{"name":"models/b"}]}`), "", nil
	})
	if err != nil || calls != 2 || len(items) != 2 || items[0].GeminiDisplayName != "A" || len(items[0].GenerationMethods) != 1 {
		t.Fatalf("invalid paged result %v %v", items, err)
	}
}

func TestGeminiOrdinaryTextTaskUsesInstalledCreateContract(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path != "/v1beta/models/gemini-fixture:generateContent" || r.Header.Get("x-goog-api-key") != "fixture" || r.URL.Query().Get("key") != "" {
			t.Errorf("unexpected request authority/path")
		}
		var body struct {
			Contents []struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"contents"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || len(body.Contents) != 1 || len(body.Contents[0].Parts) != 1 || body.Contents[0].Parts[0].Text != "fixture prompt" {
			t.Errorf("invalid Gemini message mapping")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"candidates":[{"content":{"parts":[{"text":"fixture reply"}]}}],"usageMetadata":{"promptTokenCount":1,"candidatesTokenCount":1}}`))
	}))
	defer server.Close()
	svc, _ := newChannelModelTestService(t)
	svc.runtimeCapabilities = RuntimeCapabilities{desktopLocalChannels: true}
	stream := false
	input := canvasGenerationInput{Mode: "text", Prompt: "fixture prompt", Config: providerConfig{BaseURL: server.URL, APIKey: "fixture", Model: "gemini-fixture", InterfaceType: "gemini-generate-content", APIFormat: "gemini", AllowLocalChannel: true}, TextOptions: canvasTextOptions{Stream: &stream}}
	raw, _ := json.Marshal(input)
	result, err := svc.processCanvasGenerationTask(context.Background(), "fixture-user", "", "canvas_text", "", string(raw))
	if err != nil || calls != 1 || result["text"] != "fixture reply" {
		t.Fatalf("ordinary Gemini task failed %#v %v (calls=%d)", result, err, calls)
	}
	for _, kind := range []string{"image", "video", "audio"} {
		media := providerMedia{DataURL: "data:image/png;base64,AQ=="}
		rejected := input
		switch kind {
		case "image":
			rejected.ReferenceImages = []providerMedia{media}
		case "video":
			rejected.ReferenceVideos = []providerMedia{media}
		case "audio":
			rejected.ReferenceAudios = []providerMedia{media}
		}
		ctx := withProtocolRegistry(withProviderOutboundPolicy(context.Background(), rejected.Config), svc.protocolRegistry())
		if _, err := runTextTask(ctx, rejected); err == nil {
			t.Fatalf("ordinary text silently dropped %s", kind)
		}
		if calls != 1 {
			t.Fatal("rejected media request contacted upstream")
		}
	}
}

func TestCatalogPaginationFailsAtomically(t *testing.T) {
	for _, failure := range []string{"network", "repeat", "json", "upstream"} {
		t.Run(failure, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "https://example.test/models", nil)
			calls := 0
			items, err := fetchChannelCatalogPages(req, "gemini", func(*http.Request) ([]byte, string, error) {
				calls++
				if calls == 1 || failure == "repeat" {
					return []byte(`{"models":[{"name":"a"}],"nextPageToken":"same"}`), "", nil
				}
				if failure == "network" {
					return nil, "", errors.New("fixture network failure")
				}
				if failure == "upstream" {
					return []byte(`{"error":{"message":"fixture"}}`), "", nil
				}
				return []byte(`invalid`), "", nil
			})
			if err == nil || items != nil || calls != 2 {
				t.Fatal("partial catalog must not be imported")
			}
		})
	}
}

func TestCatalogPageDisplaySizeDoesNotLimitUpstreamCollection(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://example.test/models", nil)
	calls := 0
	items, err := fetchChannelCatalogPages(req, "gemini", func(*http.Request) ([]byte, string, error) {
		calls++
		payload := channelModelsPayload{Models: make([]channelModelItem, 50)}
		if calls < 3 {
			payload.NextPageToken = string(rune('a' + calls))
		}
		data, _ := json.Marshal(payload)
		return data, "", nil
	})
	if err != nil || len(items) != 150 {
		t.Fatal("all upstream pages must be collected independently of UI page size")
	}
}
