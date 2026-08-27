package service

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var safeTaskLogTokenPattern = regexp.MustCompile(`^[A-Za-z0-9._:= -]{1,512}$`)

func safeTaskLogPayload(payload string) string {
	payload = strings.TrimSpace(payload)
	if payload == "" {
		return ""
	}
	if safeTaskLogTokenPattern.MatchString(payload) {
		return truncateTaskLogPayload(payload)
	}
	if json.Valid([]byte(payload)) {
		return truncateTaskLogPayload(SanitizeAPICallPayload([]byte(payload), "application/json"))
	}
	return fmt.Sprintf("[诊断内容已隐藏，原始长度 %d 字节]", len(payload))
}
