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
	Reason        string `json:"reason"`
	ErrorCode     string `json:"errorCode"`
	ErrorCategory string `json:"errorCategory"`
	Retryable     bool   `json:"retryable"`
	RequestID     string `json:"requestId"`
}

func TestFailServiceRegistrationCooldown(t *testing.T) {
	recorder, context := responseTestContext()
	failService(context, &service.EmailCodeCooldownError{Seconds: 47})
	response := decodeFailureEnvelope(t, recorder)
	if recorder.Code != http.StatusTooManyRequests || response.Code != service.CodeRateLimited || response.Reason != string(service.ReasonRateLimited) || response.ErrorCategory != service.ErrorCategoryQuota || !response.Retryable || recorder.Header().Get("Retry-After") != "47" || !strings.Contains(response.Msg, "47") {
		t.Fatalf("cooldown response: status=%d header=%s body=%#v", recorder.Code, recorder.Header().Get("Retry-After"), response)
	}
}

func TestFailServiceProjectsAppError(t *testing.T) {
	recorder, context := responseTestContext()
	err := service.RateLimited("请求过于频繁，请稍后重试")

	failService(context, err)

	response := decodeFailureEnvelope(t, recorder)
	if recorder.Code != http.StatusTooManyRequests || response.Code != service.CodeRateLimited || response.Reason != string(service.ReasonRateLimited) || response.ErrorCode != "request_throttled" || response.ErrorCategory != service.ErrorCategoryQuota || !response.Retryable || response.Msg != err.Message || !strings.HasPrefix(response.RequestID, "req_") {
		t.Fatalf("response = status %d, body %#v", recorder.Code, response)
	}
}

func TestFailServiceQuotaExceeded(t *testing.T) {
	recorder, context := responseTestContext()
	failService(context, service.QuotaExceeded("账号素材数量已达到 100 个上限"))

	response := decodeFailureEnvelope(t, recorder)
	if recorder.Code != http.StatusForbidden || response.Code != service.CodeQuotaExceeded || response.Reason != string(service.ReasonQuotaExceeded) || response.ErrorCode != string(service.ReasonQuotaExceeded) || response.ErrorCategory != service.ErrorCategoryQuota || response.Retryable {
		t.Fatalf("quota response = status %d, body %#v", recorder.Code, response)
	}
}

func TestFailServiceProjectsModelErrorMetadata(t *testing.T) {
	recorder, context := responseTestContext()
	failService(context, service.ProviderRequestFailed("模型服务暂时不可用"))

	response := decodeFailureEnvelope(t, recorder)
	if recorder.Code != http.StatusBadGateway || response.Code != http.StatusBadGateway || response.Reason != string(service.ErrCodeProviderRequestFailed) || response.ErrorCode != string(service.ErrCodeProviderRequestFailed) || response.ErrorCategory != service.ErrorCategoryProvider || !response.Retryable || response.Msg != "模型服务暂时不可用" {
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

func TestParsePaginationQueryUsesDefaultsAndRejectsInvalidValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name      string
		query     string
		page      int
		pageSize  int
		wantError bool
	}{
		{name: "defaults", query: "", page: 1, pageSize: 40},
		{name: "explicit values", query: "page=3&pageSize=25", page: 3, pageSize: 25},
		{name: "invalid page", query: "page=abc", wantError: true},
		{name: "zero page", query: "page=0", wantError: true},
		{name: "negative page size", query: "pageSize=-1", wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/test?"+tt.query, nil)
			context, _ := gin.CreateTestContext(httptest.NewRecorder())
			context.Request = request
			page, pageSize, err := parsePaginationQuery(context, 40)
			if (err != nil) != tt.wantError {
				t.Fatalf("error = %v, wantError = %t", err, tt.wantError)
			}
			if tt.wantError {
				return
			}
			if page != tt.page || pageSize != tt.pageSize {
				t.Fatalf("pagination = (%d, %d), want (%d, %d)", page, pageSize, tt.page, tt.pageSize)
			}
		})
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
