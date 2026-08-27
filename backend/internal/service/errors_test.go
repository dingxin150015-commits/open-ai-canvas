package service

import (
	"errors"
	"testing"
)

func TestAppErrorPreservesSafeProjectionAndCause(t *testing.T) {
	cause := errors.New("database password=secret")
	err := WrapAppError(503, "服务暂时不可用", cause)
	err.Code = 10001
	err.Retryable = true

	if err.Error() != "服务暂时不可用" {
		t.Fatalf("AppError.Error() = %q", err.Error())
	}
	if err.Status != 503 || err.Code != 10001 || !err.Retryable {
		t.Fatalf("AppError fields = %#v", err)
	}
	if err.PublicCode != "service_unavailable" || err.Category != ErrorCategoryInternal {
		t.Fatalf("AppError metadata = %#v", err)
	}
	if !errors.Is(err, cause) {
		t.Fatal("AppError should unwrap its internal cause")
	}
}

func TestModelErrorPreservesPublicCodeCategoryAndAppErrorChain(t *testing.T) {
	err := ProviderRequestFailed("模型服务暂时不可用")
	var modelErr *ModelError
	var appErr *AppError
	if !errors.As(err, &modelErr) || !errors.As(err, &appErr) {
		t.Fatalf("model error chain = %#v", err)
	}
	if appErr.Status != 502 || appErr.PublicCode != string(ErrCodeProviderRequestFailed) || appErr.Category != ErrorCategoryProvider || !appErr.Retryable {
		t.Fatalf("model app error = %#v", appErr)
	}
}

func TestAuthErrorAliasRemainsCompatible(t *testing.T) {
	var err *AuthError = BadAuthRequest("请求参数错误")
	var appErr *AppError
	if !errors.As(err, &appErr) || appErr.Status != 400 || appErr.Message != "请求参数错误" {
		t.Fatalf("AuthError compatibility = %#v", err)
	}
}
