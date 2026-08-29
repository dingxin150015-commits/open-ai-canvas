package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"infinite-canvas/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type failureEnvelope struct {
	Code          int    `json:"code"`
	Msg           string `json:"msg"`
	ErrorCode     string `json:"errorCode"`
	ErrorCategory string `json:"errorCategory"`
	Retryable     bool   `json:"retryable"`
	RequestID     string `json:"requestId"`
}

func TestFailServiceProjectsAppError(t *testing.T) {
	recorder, context := responseTestContext()
	err := service.NewAppError(http.StatusTooManyRequests, "请求过于频繁，请稍后重试")
	err.Code = 42901

	failService(context, err)

	response := decodeFailureEnvelope(t, recorder)
	if recorder.Code != http.StatusTooManyRequests || response.Code != 42901 || response.Msg != err.Message || response.ErrorCode != "request_throttled" || response.ErrorCategory != service.ErrorCategoryQuota || !response.Retryable || !strings.HasPrefix(response.RequestID, "req_") {
		t.Fatalf("response = status %d, body %#v", recorder.Code, response)
	}
}

func TestFailServiceProjectsModelErrorMetadata(t *testing.T) {
	recorder, context := responseTestContext()
	failService(context, service.ProviderRequestFailed("模型服务暂时不可用"))

	response := decodeFailureEnvelope(t, recorder)
	if recorder.Code != http.StatusBadGateway || response.Code != http.StatusBadGateway || response.ErrorCode != string(service.ErrCodeProviderRequestFailed) || response.ErrorCategory != service.ErrorCategoryProvider || !response.Retryable || response.Msg != "模型服务暂时不可用" {
		t.Fatalf("response = status %d, body %#v", recorder.Code, response)
	}
}

func TestFailHidesUnclassifiedClientError(t *testing.T) {
	recorder, context := responseTestContext()
	fail(context, http.StatusBadRequest, errors.New("invalid password=secret at C:\\private\\file"))

	response := decodeFailureEnvelope(t, recorder)
	if recorder.Code != http.StatusBadRequest || response.Msg != "请求格式或参数无效" || response.ErrorCode != "bad_request" || strings.Contains(recorder.Body.String(), "secret") || strings.Contains(recorder.Body.String(), "private") {
		t.Fatalf("unsafe client failure: %s", recorder.Body.String())
	}
}

func TestRequestCorrelationMiddlewareGeneratesRequestIDAndAcceptsSafeTraceID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RequestCorrelationMiddleware())
	router.GET("/test", func(c *gin.Context) { ok(c, gin.H{"requestId": RequestID(c), "traceId": TraceID(c)}) })

	valid := httptest.NewRecorder()
	validRequest := httptest.NewRequest(http.MethodGet, "/test", nil)
	validRequest.Header.Set("X-Request-ID", "client-request-123")
	validRequest.Header.Set("X-Canvas-Trace-ID", "client-trace-123")
	router.ServeHTTP(valid, validRequest)
	generated := valid.Header().Get("X-Request-ID")
	if !strings.HasPrefix(generated, "req_") || strings.Contains(valid.Body.String(), "client-request-123") || !strings.Contains(valid.Body.String(), generated) {
		t.Fatalf("request id was not regenerated safely: headers=%v body=%s", valid.Header(), valid.Body.String())
	}
	if valid.Header().Get("X-Canvas-Trace-ID") != "client-trace-123" || !strings.Contains(valid.Body.String(), "client-trace-123") {
		t.Fatalf("safe trace id was not preserved: headers=%v body=%s", valid.Header(), valid.Body.String())
	}

	invalid := httptest.NewRecorder()
	invalidRequest := httptest.NewRequest(http.MethodGet, "/test", nil)
	invalidRequest.Header.Set("X-Request-ID", "Bearer secret value")
	invalidRequest.Header.Set("X-Canvas-Trace-ID", "Bearer trace secret")
	router.ServeHTTP(invalid, invalidRequest)
	if !strings.HasPrefix(invalid.Header().Get("X-Request-ID"), "req_") || !strings.HasPrefix(invalid.Header().Get("X-Canvas-Trace-ID"), "trace_") || strings.Contains(invalid.Body.String(), "Bearer") {
		t.Fatalf("unsafe correlation id was accepted: headers=%v body=%s", invalid.Header(), invalid.Body.String())
	}
}

func TestFailServiceHidesUnclassifiedInternalError(t *testing.T) {
	recorder, context := responseTestContext()
	failService(context, errors.New("database password=secret"))

	response := decodeFailureEnvelope(t, recorder)
	if recorder.Code != http.StatusInternalServerError || response.Code != http.StatusInternalServerError {
		t.Fatalf("response = status %d, body %#v", recorder.Code, response)
	}
	if response.Msg != internalErrorMessage || strings.Contains(recorder.Body.String(), "password=secret") {
		t.Fatalf("internal error leaked in response: %s", recorder.Body.String())
	}
}

func TestFailInternalKeepsStatusWithoutLeakingCause(t *testing.T) {
	recorder, context := responseTestContext()
	failInternal(context, http.StatusServiceUnavailable, errors.New("redis://user:password@private-host"))

	response := decodeFailureEnvelope(t, recorder)
	if recorder.Code != http.StatusServiceUnavailable || response.Msg != "服务暂时不可用，请稍后重试" {
		t.Fatalf("response = status %d, body %#v", recorder.Code, response)
	}
	if strings.Contains(recorder.Body.String(), "private-host") {
		t.Fatalf("internal cause leaked in response: %s", recorder.Body.String())
	}
}

func TestStructuredHandlerLogContainsCorrelationWithoutCause(t *testing.T) {
	var output bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&output)
	t.Cleanup(func() { log.SetOutput(previous) })
	_, context := responseTestContext()
	context.Set(requestIDContextKey, "req_structured_12345678")

	failInternal(context, http.StatusInternalServerError, errors.New("database password=secret at C:\\private"))

	logged := output.String()
	for _, expected := range []string{"request_id=req_structured_12345678", "status=500", "category=internal", "error_code=internal_error"} {
		if !strings.Contains(logged, expected) {
			t.Fatalf("structured log missing %q: %s", expected, logged)
		}
	}
	if strings.Contains(logged, "secret") || strings.Contains(logged, "private") {
		t.Fatalf("structured log leaked cause: %s", logged)
	}
}

func responseTestContext() (*httptest.ResponseRecorder, *gin.Context) {
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodGet, "/api/test", nil)
	return recorder, context
}

func decodeFailureEnvelope(t *testing.T, recorder *httptest.ResponseRecorder) failureEnvelope {
	t.Helper()
	var response failureEnvelope
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return response
}
