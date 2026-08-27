package handler

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"infinite-canvas/backend/internal/service"

	"github.com/gin-gonic/gin"
)

const internalErrorMessage = "系统处理失败，请稍后重试"
const requestIDContextKey = "canvas_request_id"

var requestIDFallbackCounter atomic.Uint64

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := normalizeRequestID(c.GetHeader("X-Request-ID"))
		if requestID == "" {
			requestID = newRequestID()
		}
		c.Set(requestIDContextKey, requestID)
		c.Request.Header.Set("X-Request-ID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

func RequestID(c *gin.Context) string {
	if c != nil {
		if value, exists := c.Get(requestIDContextKey); exists {
			if requestID := normalizeRequestID(fmt.Sprint(value)); requestID != "" {
				return requestID
			}
		}
	}
	requestID := newRequestID()
	if c != nil {
		c.Set(requestIDContextKey, requestID)
		c.Header("X-Request-ID", requestID)
	}
	return requestID
}

func HandleRecovery(c *gin.Context, recovered any) {
	failInternal(c, http.StatusInternalServerError, fmt.Errorf("recovered panic type %T", recovered))
}

func ok(c *gin.Context, data any) {
	_ = RequestID(c)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": data, "msg": "ok"})
}

// fail 只接受调用方已经确认可公开的错误；service 返回值统一交给 failService 投影。
func fail(c *gin.Context, status int, err error) {
	if isStructuredServiceError(err) {
		failService(c, err)
		return
	}
	if status >= http.StatusInternalServerError {
		failInternal(c, status, err)
		return
	}
	publicCode, category, retryable := handlerErrorMetadata(status)
	logHandlerFailure(c, status, category, publicCode, retryable, err)
	writeFailure(c, status, status, safeClientErrorMessage(status), publicCode, category, retryable)
}

func failService(c *gin.Context, err error) {
	var modelErr *service.ModelError
	if errors.As(err, &modelErr) && modelErr.AppError != nil {
		writeAppFailure(c, modelErr.AppError, string(modelErr.ErrorCode), err)
		return
	}
	var appErr *service.AppError
	if errors.As(err, &appErr) && validErrorStatus(appErr.Status) {
		writeAppFailure(c, appErr, appErr.PublicCode, err)
		return
	}
	failInternal(c, http.StatusInternalServerError, err)
}

func writeAppFailure(c *gin.Context, appErr *service.AppError, publicCode string, diagnostic error) {
	status := appErr.Status
	if !validErrorStatus(status) {
		status = http.StatusInternalServerError
	}
	code := appErr.Code
	if code == 0 {
		code = status
	}
	if strings.TrimSpace(publicCode) == "" {
		publicCode = appErr.PublicCode
	}
	category := strings.TrimSpace(appErr.Category)
	defaultCode, defaultCategory, defaultRetryable := handlerErrorMetadata(status)
	if strings.TrimSpace(publicCode) == "" {
		publicCode = defaultCode
	}
	if category == "" {
		category = defaultCategory
	}
	retryable := appErr.Retryable || defaultRetryable
	message := strings.TrimSpace(appErr.Message)
	if message == "" {
		message = safeInternalErrorMessage(status)
	}
	logHandlerFailure(c, status, category, publicCode, retryable, diagnostic)
	writeFailure(c, status, code, message, publicCode, category, retryable)
}

// failInternal 保留真实 HTTP 状态，但绝不把未分类错误原文写入响应。
func failInternal(c *gin.Context, status int, err error) {
	if !validErrorStatus(status) {
		status = http.StatusInternalServerError
	}
	publicCode, category, retryable := handlerErrorMetadata(status)
	logHandlerFailure(c, status, category, publicCode, retryable, err)
	writeFailure(c, status, status, safeInternalErrorMessage(status), publicCode, category, retryable)
}

func writeFailure(c *gin.Context, status int, code int, message string, publicCode string, category string, retryable bool) {
	c.JSON(status, gin.H{
		"code": code, "data": nil, "msg": message,
		"errorCode": publicCode, "errorCategory": category, "retryable": retryable, "requestId": RequestID(c),
	})
}

func validErrorStatus(status int) bool {
	return status >= http.StatusBadRequest && status <= 599
}

func safeInternalErrorMessage(status int) string {
	switch status {
	case http.StatusBadGateway:
		return "上游服务暂时不可用，请稍后重试"
	case http.StatusServiceUnavailable:
		return "服务暂时不可用，请稍后重试"
	case http.StatusGatewayTimeout:
		return "上游服务响应超时，请稍后重试"
	default:
		return internalErrorMessage
	}
}

func logHandlerFailure(c *gin.Context, status int, category string, publicCode string, retryable bool, err error) {
	method := ""
	route := ""
	if c != nil && c.Request != nil {
		method = c.Request.Method
		route = c.FullPath()
		if route == "" {
			route = "<unmatched>"
		}
	}
	// 只记录关联键、分类和 Go 类型；错误原文、URL、请求体和凭据均不得进入进程日志。
	log.Printf("handler request failed: request_id=%s method=%s route=%s status=%d category=%s error_code=%s retryable=%t error_type=%T", RequestID(c), method, route, status, category, publicCode, retryable, err)
}

func isStructuredServiceError(err error) bool {
	if err == nil {
		return false
	}
	var modelErr *service.ModelError
	if errors.As(err, &modelErr) {
		return true
	}
	var appErr *service.AppError
	return errors.As(err, &appErr)
}

func handlerErrorMetadata(status int) (string, string, bool) {
	switch status {
	case http.StatusBadRequest:
		return "bad_request", service.ErrorCategoryValidation, false
	case http.StatusUnauthorized:
		return "authentication_required", service.ErrorCategoryAuthentication, false
	case http.StatusForbidden:
		return "access_denied", service.ErrorCategoryAuthorization, false
	case http.StatusNotFound:
		return "not_found", service.ErrorCategoryNotFound, false
	case http.StatusConflict:
		return "state_conflict", service.ErrorCategoryConflict, false
	case http.StatusRequestTimeout, http.StatusGatewayTimeout:
		return "request_timeout", service.ErrorCategoryTimeout, true
	case http.StatusTooManyRequests:
		return "request_throttled", service.ErrorCategoryQuota, true
	case http.StatusBadGateway:
		return "upstream_unavailable", service.ErrorCategoryProvider, true
	case http.StatusServiceUnavailable:
		return "service_unavailable", service.ErrorCategoryInternal, true
	default:
		return "internal_error", service.ErrorCategoryInternal, status >= 500
	}
}

func safeClientErrorMessage(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "请求格式或参数无效"
	case http.StatusUnauthorized:
		return "请先登录"
	case http.StatusForbidden:
		return "没有权限执行此操作"
	case http.StatusNotFound:
		return "请求的资源不存在"
	case http.StatusConflict:
		return "数据状态已变化，请刷新后重试"
	case http.StatusTooManyRequests:
		return "请求过于频繁，请稍后重试"
	default:
		return http.StatusText(status)
	}
}

func normalizeRequestID(value string) string {
	value = strings.TrimSpace(value)
	if len(value) < 8 || len(value) > 128 {
		return ""
	}
	for _, character := range value {
		if (character >= 'a' && character <= 'z') || (character >= 'A' && character <= 'Z') || (character >= '0' && character <= '9') || strings.ContainsRune("._:-", character) {
			continue
		}
		return ""
	}
	return value
}

func newRequestID() string {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err == nil {
		return "req_" + hex.EncodeToString(buffer)
	}
	return fmt.Sprintf("req_fallback_%x_%x", time.Now().UnixNano(), requestIDFallbackCounter.Add(1))
}
