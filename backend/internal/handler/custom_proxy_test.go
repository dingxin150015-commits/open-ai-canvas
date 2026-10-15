package handler

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"yingce/backend/internal/service"

	"github.com/gin-gonic/gin"
)

func TestCustomRelayForwardsOpenAIRequestWithoutBrowserHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const apiKey = "relay-secret-key"
	upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+apiKey {
			t.Errorf("Authorization = %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("User-Agent") != "Custom Relay Agent" {
			t.Errorf("User-Agent = %q", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("X-Gateway-Tenant") != "tenant-a" {
			t.Errorf("X-Gateway-Tenant = %q", r.Header.Get("X-Gateway-Tenant"))
		}
		for _, name := range []string{"Cookie", "Origin", "Referer", "X-Canvas-Upstream-URL", "X-Forwarded-For"} {
			if value := r.Header.Get(name); value != "" {
				t.Errorf("upstream received %s = %q", name, value)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Set-Cookie", "upstream=unsafe")
		_, _ = io.WriteString(w, `{"data":[]}`)
	}))
	defer upstream.Close()
	useCustomRelayTestClient(t, upstream)

	request := httptest.NewRequest(http.MethodGet, "/api/ai/custom", nil)
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("X-Canvas-Upstream-URL", "https://8.8.8.8/v1/models")
	request.Header.Set("X-Canvas-Upstream-Format", "openai")
	request.Header.Set(service.CustomRelayHeadersHeader, base64.StdEncoding.EncodeToString([]byte(`[{"name":"User-Agent","value":"Custom Relay Agent"},{"name":"X-Gateway-Tenant","value":"tenant-a"}]`)))
	request.Header.Set("Cookie", "browser=session")
	request.Header.Set("Origin", "https://canvas.example.com")
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request = request

	proxyCustomRelayRequest(context, defaultCustomRelayTestPolicy())
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if response.Header().Get("Set-Cookie") != "" {
		t.Fatal("upstream Set-Cookie should not be forwarded")
	}
	if strings.Contains(response.Body.String(), apiKey) {
		t.Fatal("response leaked API key")
	}
}

func TestCustomRelayReturnsVideoContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "video/mp4")
		_, _ = w.Write([]byte("video-content"))
	}))
	defer upstream.Close()
	useCustomRelayTestClient(t, upstream)

	request := httptest.NewRequest(http.MethodGet, "/api/ai/custom", nil)
	request.Header.Set("Authorization", "Bearer test-key")
	request.Header.Set("X-Canvas-Upstream-URL", "https://8.8.8.8/v1/videos/task-1/content")
	request.Header.Set("X-Canvas-Upstream-Format", "openai")
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request = request

	proxyCustomRelayRequest(context, defaultCustomRelayTestPolicy())
	if response.Code != http.StatusOK || response.Header().Get("Content-Type") != "video/mp4" || response.Body.String() != "video-content" {
		t.Fatalf("status = %d, content-type = %q, body = %q", response.Code, response.Header().Get("Content-Type"), response.Body.String())
	}
}

func TestCustomRelayConvertsGeminiAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const apiKey = "gemini-secret-key"
	upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-goog-api-key") != apiKey {
			t.Errorf("x-goog-api-key = %q", r.Header.Get("x-goog-api-key"))
		}
		if r.Header.Get("Authorization") != "" {
			t.Errorf("Authorization should not be forwarded, got %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"candidates":[]}`)
	}))
	defer upstream.Close()
	useCustomRelayTestClient(t, upstream)

	request := httptest.NewRequest(http.MethodPost, "/api/ai/custom", strings.NewReader(`{"contents":[]}`))
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Canvas-Upstream-URL", "https://8.8.8.8/v1beta/models/gemini-test:generateContent")
	request.Header.Set("X-Canvas-Upstream-Format", "gemini")
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request = request

	proxyCustomRelayRequest(context, defaultCustomRelayTestPolicy())
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestCustomRelayConvertsClaudeAuthentication(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const apiKey = "claude-secret-key"
	upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/messages" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Header.Get("x-api-key") != apiKey {
			t.Errorf("x-api-key = %q", r.Header.Get("x-api-key"))
		}
		if r.Header.Get("Authorization") != "" {
			t.Errorf("Authorization should not be forwarded, got %q", r.Header.Get("Authorization"))
		}
		if r.Header.Get("anthropic-version") != "2023-06-01" {
			t.Errorf("anthropic-version = %q", r.Header.Get("anthropic-version"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"content":[{"type":"text","text":"claude relay works"}]}`)
	}))
	defer upstream.Close()
	useCustomRelayTestClient(t, upstream)

	request := httptest.NewRequest(http.MethodPost, "/api/ai/custom", strings.NewReader(`{"model":"claude-test","messages":[]}`))
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Canvas-Upstream-URL", "https://8.8.8.8/v1/messages")
	request.Header.Set("X-Canvas-Upstream-Format", "claude")
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request = request

	proxyCustomRelayRequest(context, defaultCustomRelayTestPolicy())
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "claude relay works") {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
}

func TestCustomRelayStreamsBeforeUpstreamCompletes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	const apiKey = "stream-secret"
	release := make(chan struct{})
	upstream := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: first\n\n")
		w.(http.Flusher).Flush()
		<-release
		_, _ = io.WriteString(w, "data: second\n\n")
	}))
	defer upstream.Close()
	useCustomRelayTestClient(t, upstream)

	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context, _ := gin.CreateTestContext(w)
		context.Request = r
		proxyCustomRelayRequest(context, defaultCustomRelayTestPolicy())
	}))
	defer proxy.Close()
	request, _ := http.NewRequest(http.MethodPost, proxy.URL, strings.NewReader(`{"model":"test"}`))
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("X-Canvas-Upstream-URL", "https://8.8.8.8/v1/responses")
	request.Header.Set("X-Canvas-Upstream-Format", "openai")
	client := &http.Client{Timeout: 3 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		close(release)
		t.Fatal(err)
	}
	reader := bufio.NewReader(response.Body)
	line, err := reader.ReadString('\n')
	if err != nil {
		close(release)
		_ = response.Body.Close()
		t.Fatal(err)
	}
	if line != "data: first\n" {
		close(release)
		_ = response.Body.Close()
		t.Fatalf("first streamed line = %q", line)
	}
	close(release)
	_ = response.Body.Close()
}

func TestRelayStreamRedactorHandlesSplitSecret(t *testing.T) {
	redactor := newRelayStreamRedactor("split-secret")
	output := append(redactor.Push([]byte("before split-"), false), redactor.Push([]byte("secret after"), true)...)
	if bytes.Contains(output, []byte("split-secret")) || !bytes.Contains(output, []byte("[REDACTED]")) {
		t.Fatalf("redacted output = %q", output)
	}
}

func TestCustomRelayRejectsOversizedDeclaredBodyBeforeConnecting(t *testing.T) {
	gin.SetMode(gin.TestMode)
	connected := false
	previous := customRelayClient
	customRelayClient = func(time.Duration) *http.Client {
		connected = true
		return http.DefaultClient
	}
	t.Cleanup(func() { customRelayClient = previous })

	request := httptest.NewRequest(http.MethodPost, "/api/ai/custom", strings.NewReader(`{"model":"test"}`))
	request.ContentLength = (defaultCustomRelayTestPolicy().CustomRelayRequestMB << 20) + 1
	request.Header.Set("Authorization", "Bearer test-key")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Canvas-Upstream-URL", "https://8.8.8.8/v1/responses")
	request.Header.Set("X-Canvas-Upstream-Format", "openai")
	response := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(response)
	context.Request = request

	proxyCustomRelayRequest(context, defaultCustomRelayTestPolicy())
	if response.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	if connected {
		t.Fatal("oversized request should not create an upstream client")
	}
}

func useCustomRelayTestClient(t *testing.T, upstream *httptest.Server) {
	t.Helper()
	upstreamURL, err := url.Parse(upstream.URL)
	if err != nil {
		t.Fatal(err)
	}
	client := upstream.Client()
	transport := client.Transport
	client.Transport = relayTestRoundTripper(func(request *http.Request) (*http.Response, error) {
		if request.URL.Host != "8.8.8.8" {
			return nil, errors.New("unexpected relay test host")
		}
		forwarded := request.Clone(request.Context())
		forwarded.URL.Scheme = upstreamURL.Scheme
		forwarded.URL.Host = upstreamURL.Host
		forwarded.Host = upstreamURL.Host
		return transport.RoundTrip(forwarded)
	})
	previous := customRelayClient
	customRelayClient = func(time.Duration) *http.Client { return client }
	t.Cleanup(func() { customRelayClient = previous })
}

type relayTestRoundTripper func(*http.Request) (*http.Response, error)

func (roundTrip relayTestRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return roundTrip(request)
}

func defaultCustomRelayTestPolicy() service.RuntimeRequestPolicy {
	return service.RuntimeRequestPolicy{
		CustomRelayRequestMB: 32, CustomRelayResponseMB: 32, CustomRelayTimeoutMinutes: 10,
	}
}
