package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	requestIDHeader = "X-Request-ID"
	traceIDHeader   = "X-Canvas-Trace-ID"
	requestIDKey    = "canvas.request_id"
	traceIDKey      = "canvas.trace_id"
	// requestIDContextKey 保留给同包测试和旧内部调用；唯一真相仍是 requestIDKey。
	requestIDContextKey = requestIDKey
)

var correlationIDPattern = regexp.MustCompile(`^[A-Za-z0-9._:-]{1,96}$`)
var correlationFallbackCounter atomic.Uint64

// RequestCorrelationMiddleware 为每个请求生成服务端 requestId，并保留一次业务操作的 traceId。
// requestId 不能被客户端覆盖；traceId 只接受有限字符集，避免把任意请求头内容写入日志和诊断包。
func RequestCorrelationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := newCorrelationID("req")
		traceID := normalizeCorrelationID(c.GetHeader(traceIDHeader))
		if traceID == "" {
			traceID = newCorrelationID("trace")
		}
		c.Set(requestIDKey, requestID)
		c.Set(traceIDKey, traceID)
		c.Request.Header.Set(requestIDHeader, requestID)
		c.Request.Header.Set(traceIDHeader, traceID)
		c.Header(requestIDHeader, requestID)
		c.Header(traceIDHeader, traceID)
		c.Next()
	}
}

// RequestIDMiddleware 是旧内部名称的兼容别名；请求编号仍由服务端生成。
func RequestIDMiddleware() gin.HandlerFunc {
	return RequestCorrelationMiddleware()
}

func RequestID(c *gin.Context) string {
	if c != nil {
		if value, ok := c.Get(requestIDKey); ok {
			if result, ok := value.(string); ok && normalizeCorrelationID(result) != "" {
				return result
			}
		}
	}
	result := newCorrelationID("req")
	setCorrelationID(c, requestIDKey, requestIDHeader, result)
	return result
}

func TraceID(c *gin.Context) string {
	if c != nil {
		if value, ok := c.Get(traceIDKey); ok {
			if result, ok := value.(string); ok && normalizeCorrelationID(result) != "" {
				return result
			}
		}
	}
	result := newCorrelationID("trace")
	setCorrelationID(c, traceIDKey, traceIDHeader, result)
	return result
}

func setCorrelationID(c *gin.Context, key string, header string, value string) {
	if c == nil {
		return
	}
	c.Set(key, value)
	if c.Request != nil {
		c.Request.Header.Set(header, value)
	}
	c.Header(header, value)
}

func normalizeCorrelationID(value string) string {
	value = strings.TrimSpace(value)
	if !correlationIDPattern.MatchString(value) {
		return ""
	}
	return value
}

func newCorrelationID(prefix string) string {
	var raw [12]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return fmt.Sprintf("%s_fallback_%x_%x", prefix, time.Now().UnixNano(), correlationFallbackCounter.Add(1))
	}
	return prefix + "_" + hex.EncodeToString(raw[:])
}
