package service

import "fmt"

const (
	ErrorCategoryValidation     = "validation"
	ErrorCategoryAuthentication = "authentication"
	ErrorCategoryAuthorization  = "authorization"
	ErrorCategoryNotFound       = "not_found"
	ErrorCategoryConflict       = "conflict"
	ErrorCategoryQuota          = "quota"
	ErrorCategoryCapability     = "capability"
	ErrorCategoryPricing        = "pricing"
	ErrorCategoryRouting        = "routing"
	ErrorCategoryProvider       = "provider"
	ErrorCategoryTimeout        = "timeout"
	ErrorCategoryInternal       = "internal"
)

// AppError 是 service 层对外公开的结构化错误。
// Message 必须可安全展示给用户，Cause 仅用于保留内部诊断链路，不得直接写入 HTTP 响应。
type AppError struct {
	Status     int
	Code       int
	PublicCode string
	Category   string
	Message    string
	Retryable  bool
	Cause      error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func NewAppError(status int, message string) *AppError {
	publicCode, category, retryable := defaultAppErrorMetadata(status)
	return &AppError{Status: status, Code: status, PublicCode: publicCode, Category: category, Message: message, Retryable: retryable}
}

func WrapAppError(status int, message string, cause error) *AppError {
	error := NewAppError(status, message)
	error.Cause = cause
	return error
}

func defaultAppErrorMetadata(status int) (string, string, bool) {
	switch status {
	case 400:
		return "bad_request", ErrorCategoryValidation, false
	case 401:
		return "authentication_required", ErrorCategoryAuthentication, false
	case 403:
		return "access_denied", ErrorCategoryAuthorization, false
	case 404:
		return "not_found", ErrorCategoryNotFound, false
	case 408, 504:
		return "request_timeout", ErrorCategoryTimeout, true
	case 409:
		return "state_conflict", ErrorCategoryConflict, false
	case 425, 429:
		return "request_throttled", ErrorCategoryQuota, true
	case 502:
		return "upstream_unavailable", ErrorCategoryProvider, true
	case 503:
		return "service_unavailable", ErrorCategoryInternal, true
	default:
		return "internal_error", ErrorCategoryInternal, status >= 500
	}
}

// ModelErrorCode 定义模型相关的错误码
type ModelErrorCode string

const (
	// ErrCodeModelCapabilityNotSupported 当前模型能力不支持请求
	ErrCodeModelCapabilityNotSupported ModelErrorCode = "model_capability_not_supported"
	// ErrCodeModelPriceNotConfigured 当前模型未配置该能力组合价格
	ErrCodeModelPriceNotConfigured ModelErrorCode = "model_price_not_configured"
	// ErrCodeModelRouteUnavailable 没有可用供应线路
	ErrCodeModelRouteUnavailable ModelErrorCode = "model_route_unavailable"
	// ErrCodeProviderRequestFailed 供应商异常、响应格式错误或上游失败
	ErrCodeProviderRequestFailed ModelErrorCode = "provider_request_failed"
	// ErrCodeModelCatalogMismatch 模型目录已更新，请重新选择
	ErrCodeModelCatalogMismatch ModelErrorCode = "model_catalog_mismatch"
	// ErrCodeInvalidModelSelection 无效的模型选择
	ErrCodeInvalidModelSelection ModelErrorCode = "invalid_model_selection"
)

// ModelError 模型相关的错误，包含错误码和详细信息
// 继承 AppError 以保持兼容性
type ModelError struct {
	*AppError
	ErrorCode ModelErrorCode
	Details   map[string]any
}

func (e *ModelError) Error() string {
	if e.AppError != nil {
		return e.AppError.Error()
	}
	return string(e.ErrorCode)
}

func (e *ModelError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.AppError
}

// NewModelError 创建模型错误
func NewModelError(code ModelErrorCode, message string) *ModelError {
	status := 400
	category := ErrorCategoryValidation
	retryable := false
	switch code {
	case ErrCodeModelCapabilityNotSupported:
		category = ErrorCategoryCapability
	case ErrCodeModelPriceNotConfigured:
		category = ErrorCategoryPricing
	case ErrCodeModelRouteUnavailable:
		status = 503
		category = ErrorCategoryRouting
		retryable = true
	case ErrCodeProviderRequestFailed:
		status = 502
		category = ErrorCategoryProvider
		retryable = true
	case ErrCodeModelCatalogMismatch:
		status = 409
		category = ErrorCategoryConflict
	case ErrCodeInvalidModelSelection:
		category = ErrorCategoryValidation
	}
	appError := NewAppError(status, message)
	appError.PublicCode = string(code)
	appError.Category = category
	appError.Retryable = retryable
	return &ModelError{
		AppError:  appError,
		ErrorCode: code,
		Details:   make(map[string]any),
	}
}

// WithDetails 添加错误详情
func (e *ModelError) WithDetails(details map[string]any) *ModelError {
	e.Details = details
	return e
}

// ModelCapabilityNotSupported 当前模型能力不支持请求
func ModelCapabilityNotSupported(message string) error {
	if message == "" {
		message = "当前模型不支持该能力"
	}
	return NewModelError(ErrCodeModelCapabilityNotSupported, message)
}

// ModelPriceNotConfigured 当前模型未配置价格
func ModelPriceNotConfigured(message string) error {
	if message == "" {
		message = "当前模型未配置该能力组合价格"
	}
	return NewModelError(ErrCodeModelPriceNotConfigured, message)
}

// ModelRouteUnavailable 没有可用供应线路
func ModelRouteUnavailable(message string) error {
	if message == "" {
		message = "没有可用供应线路"
	}
	return NewModelError(ErrCodeModelRouteUnavailable, message)
}

// ProviderRequestFailed 供应商请求失败
func ProviderRequestFailed(message string) error {
	if message == "" {
		message = "模型服务返回失败，请检查请求内容或渠道配置"
	}
	return NewModelError(ErrCodeProviderRequestFailed, message)
}

// ModelCatalogMismatch 模型目录已更新
func ModelCatalogMismatch(message string) error {
	if message == "" {
		message = "模型目录已更新，请重新选择"
	}
	return NewModelError(ErrCodeModelCatalogMismatch, message)
}

// InvalidModelSelection 无效的模型选择
func InvalidModelSelection(message string) error {
	if message == "" {
		message = "无效的模型选择"
	}
	return NewModelError(ErrCodeInvalidModelSelection, message)
}

// IsModelError 判断是否为模型错误
func IsModelError(err error) bool {
	_, ok := err.(*ModelError)
	return ok
}

// GetModelErrorCode 获取模型错误码
func GetModelErrorCode(err error) ModelErrorCode {
	if modelErr, ok := err.(*ModelError); ok {
		return modelErr.ErrorCode
	}
	return ""
}

// FormatModelError 格式化模型错误消息
func FormatModelError(err error) string {
	if modelErr, ok := err.(*ModelError); ok {
		if len(modelErr.Details) > 0 {
			return fmt.Sprintf("%s (错误码: %s)", modelErr.Message, modelErr.ErrorCode)
		}
		return modelErr.Message
	}
	return err.Error()
}
