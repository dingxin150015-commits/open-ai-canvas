package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"infinite-canvas/backend/internal/service"

	"github.com/gin-gonic/gin"
)

const internalErrorMessage = "系统处理失败，请稍后重试"

func HandleRecovery(c *gin.Context, recovered any) {
	failInternal(c, http.StatusInternalServerError, fmt.Errorf("recovered panic type %T", recovered))
}

func ok(c *gin.Context, data any) {
	_ = RequestID(c)
	c.JSON(http.StatusOK, gin.H{"code": service.CodeOK, "data": data, "msg": "ok"})
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
	reason := service.ReasonForStatus(status)
	logHandlerFailure(c, status, reason, false, err)
	writeFailure(c, status, status, reason, safeClientErrorMessage(status), false)
}

func failService(c *gin.Context, err error) {
	var cooldown *service.EmailCodeCooldownError
	if errors.As(err, &cooldown) {
		c.Header("Retry-After", strconv.Itoa(cooldown.Seconds))
		writeFailure(c, http.StatusTooManyRequests, service.CodeRateLimited, service.ReasonRateLimited, cooldown.Error(), true)
		return
	}
	var modelErr *service.ModelError
	if errors.As(err, &modelErr) && modelErr.AppError != nil && validErrorStatus(modelErr.Status) {
		writeAppError(c, modelErr.AppError)
		return
	}
	var appErr *service.AppError
	if errors.As(err, &appErr) && validErrorStatus(appErr.Status) {
		writeAppError(c, appErr)
		return
	}
	failInternal(c, http.StatusInternalServerError, err)
}

func writeAppError(c *gin.Context, appErr *service.AppError) {
	code := appErr.Code
	if code == 0 {
		code = appErr.Status
	}
	reason := appErr.Reason
	if reason == "" {
		reason = service.ReasonForStatus(appErr.Status)
	}
	message := strings.TrimSpace(appErr.Message)
	if message == "" {
		message = safeInternalErrorMessage(appErr.Status)
	}
	if appErr.Status >= http.StatusInternalServerError {
		diagnosticErr := appErr.Cause
		if diagnosticErr == nil {
			diagnosticErr = appErr
		}
		logHandlerFailure(c, appErr.Status, reason, appErr.Retryable, diagnosticErr)
	}
	writeFailure(c, appErr.Status, code, reason, message, appErr.Retryable)
}

// failInternal 保留真实 HTTP 状态，但绝不把未分类错误原文写入响应。
func failInternal(c *gin.Context, status int, err error) {
	if !validErrorStatus(status) {
		status = http.StatusInternalServerError
	}
	reason := service.ReasonForStatus(status)
	logHandlerFailure(c, status, reason, true, err)
	writeFailure(c, status, status, reason, safeInternalErrorMessage(status), true)
}

func writeFailure(c *gin.Context, status int, code int, reason service.ErrorReason, message string, retryable bool) {
	body := gin.H{"code": code, "data": nil, "msg": message}
	if reason != "" {
		body["reason"] = reason
	}
	publicCode, category, defaultRetryable := handlerErrorMetadata(status)
	if reason != "" {
		publicCode, category, defaultRetryable = handlerMetadataForReason(reason, status)
	}
	body["errorCode"] = publicCode
	body["errorCategory"] = category
	body["retryable"] = retryable || defaultRetryable
	body["requestId"] = RequestID(c)
	c.JSON(status, body)
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

func logHandlerFailure(c *gin.Context, status int, reason service.ErrorReason, retryable bool, err error) {
	if err == nil {
		return
	}
	method := ""
	route := ""
	if c != nil && c.Request != nil {
		method = c.Request.Method
		route = c.FullPath()
		if route == "" {
			route = "<unmatched>"
		}
	}
	publicCode, category, defaultRetryable := handlerMetadataForReason(reason, status)
	// 只记录关联键、分类和 Go 类型；错误原文、URL、请求体和凭据不得进入进程日志。
	log.Printf("handler request failed: request_id=%s method=%s route=%s status=%d category=%s error_code=%s retryable=%t error_type=%T", RequestID(c), method, route, status, category, publicCode, retryable || defaultRetryable, err)
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

func handlerMetadataForReason(reason service.ErrorReason, status int) (string, string, bool) {
	switch string(reason) {
	case string(service.ErrCodeModelCapabilityNotSupported):
		return string(reason), service.ErrorCategoryCapability, false
	case string(service.ErrCodeModelPriceNotConfigured):
		return string(reason), service.ErrorCategoryPricing, false
	case string(service.ErrCodeModelRouteUnavailable):
		return string(reason), service.ErrorCategoryRouting, true
	case string(service.ErrCodeProviderRequestFailed):
		return string(reason), service.ErrorCategoryProvider, true
	case string(service.ErrCodeModelCatalogMismatch):
		return string(reason), service.ErrorCategoryConflict, false
	case string(service.ErrCodeInvalidModelSelection):
		return string(reason), service.ErrorCategoryValidation, false
	case string(service.ReasonQuotaExceeded):
		return string(reason), service.ErrorCategoryQuota, false
	case string(service.ReasonRateLimited):
		return "request_throttled", service.ErrorCategoryQuota, true
	default:
		return handlerErrorMetadata(status)
	}
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
