package service

import (
	"strings"
	"testing"
)

func TestSanitizeAPICallPayloadHidesPromptsSecretsAndSignedURLs(t *testing.T) {
	payload := []byte(`{
        "model":"qwen-image-3.0-pro",
        "prompt":"private cinematic prompt",
        "messages":[{"role":"user","content":"private conversation"}],
        "input":{"text":"private input"},
        "api_key":"secret-key",
        "result":"https://cdn.example.com/output.png?Signature=secret-signature&token=secret-token",
        "size":"1024*1024"
    }`)
	sanitized := SanitizeAPICallPayload(payload, "application/json")
	for _, secret := range []string{"private cinematic prompt", "private conversation", "private input", "secret-key", "secret-signature", "secret-token"} {
		if strings.Contains(sanitized, secret) {
			t.Fatalf("sanitized payload leaked %q: %s", secret, sanitized)
		}
	}
	for _, diagnostic := range []string{"qwen-image-3.0-pro", "1024*1024", "[文本内容已隐藏", "[REDACTED]"} {
		if !strings.Contains(sanitized, diagnostic) {
			t.Fatalf("sanitized payload missing %q: %s", diagnostic, sanitized)
		}
	}
}

func TestSafeTaskLogPayloadAllowsIdentifiersAndRejectsUnstructuredDiagnostics(t *testing.T) {
	if got := safeTaskLogPayload("provider_task_123:processing"); got != "provider_task_123:processing" {
		t.Fatalf("safe token = %q", got)
	}
	unsafe := safeTaskLogPayload("prompt content with password=secret https://private.example/path")
	if !strings.Contains(unsafe, "诊断内容已隐藏") || strings.Contains(unsafe, "secret") || strings.Contains(unsafe, "private.example") {
		t.Fatalf("unsafe payload = %q", unsafe)
	}
	jsonPayload := safeTaskLogPayload(`{"stage":"poll","prompt":"private prompt","providerRequestId":"task-1"}`)
	if strings.Contains(jsonPayload, "private prompt") || !strings.Contains(jsonPayload, "task-1") {
		t.Fatalf("JSON payload = %q", jsonPayload)
	}
}
