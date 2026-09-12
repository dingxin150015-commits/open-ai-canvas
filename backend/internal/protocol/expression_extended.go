package protocol

import (
	"encoding/base64"
	"fmt"
	"math"
	"net"
	"net/url"
	"strconv"
	"strings"
	"unicode"
)

func evaluateExtendedManifestOperator(operator string, operand any, env map[string]any) (any, error) {
	value, err := evaluateManifestValue(operand, env)
	if err != nil {
		return nil, err
	}
	spec := manifestObject(value)
	switch operator {
	case "$flatten":
		result := make([]any, 0)
		for _, item := range manifestArray(value) {
			result = append(result, manifestArray(item)...)
		}
		return result, nil
	case "$strictInt":
		input := spec["value"]
		if manifestEmpty(input) {
			input = spec["default"]
		}
		if input == nil {
			return nil, nil
		}
		text := manifestString(input)
		number, err := strconv.ParseInt(strings.TrimSpace(text), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("必须是整数")
		}
		for _, allowed := range manifestArray(spec["allow"]) {
			if float64(number) == manifestFloat(allowed) {
				return number, nil
			}
		}
		if lower, ok := spec["min"]; ok && float64(number) < manifestFloat(lower) {
			return nil, fmt.Errorf("整数小于允许范围")
		}
		if upper, ok := spec["max"]; ok && float64(number) > manifestFloat(upper) {
			return nil, fmt.Errorf("整数超过允许范围")
		}
		return number, nil
	case "$strictBool":
		input := spec["value"]
		if manifestEmpty(input) {
			input = spec["default"]
		}
		if input == nil {
			return nil, nil
		}
		if b, ok := input.(bool); ok {
			return b, nil
		}
		if text, ok := input.(string); ok {
			switch strings.ToLower(strings.TrimSpace(text)) {
			case "true", "1":
				return true, nil
			case "false", "0":
				return false, nil
			}
		}
		return nil, fmt.Errorf("必须是布尔值")
	case "$replace":
		args := manifestArray(value)
		if len(args) != 3 {
			return nil, fmt.Errorf("$replace requires three operands")
		}
		return strings.ReplaceAll(manifestString(args[0]), manifestString(args[1]), manifestString(args[2])), nil
	case "$sum":
		total := 0.0
		for _, item := range manifestArray(value) {
			total += manifestFloat(item)
		}
		return normalizeManifestNumber(total), nil
	case "$weightedLength":
		total := 0
		for _, r := range manifestString(spec["value"]) {
			if unicode.Is(unicode.Han, r) {
				total += manifestInt(spec["hanWeight"])
			} else {
				total++
			}
		}
		return total, nil
	case "$selectByIds":
		items := manifestArray(spec["from"])
		ids := manifestArray(spec["ids"])
		wanted := make([]string, 0)
		for _, id := range ids {
			if text := strings.TrimSpace(manifestString(id)); text != "" {
				wanted = append(wanted, text)
			}
		}
		if len(wanted) == 0 {
			return items, nil
		}
		byID := make(map[string]any, len(items))
		for _, item := range items {
			id := manifestString(manifestObject(item)["id"])
			if id == "" || byID[id] != nil {
				return nil, fmt.Errorf("帧素材 ID 必须非空且唯一")
			}
			byID[id] = item
		}
		result := make([]any, 0, len(wanted))
		seen := make(map[string]bool)
		for _, id := range wanted {
			item, ok := byID[id]
			if !ok || seen[id] {
				return nil, fmt.Errorf("首尾帧选择缺失或重复")
			}
			seen[id] = true
			result = append(result, item)
		}
		if len(result) != len(items) {
			return nil, fmt.Errorf("帧模式不能混入普通参考素材")
		}
		return result, nil
	case "$pixelSize":
		return resolveManifestPixelSize(spec)
	case "$validateMedia":
		return validateManifestMedia(spec)
	default:
		return nil, fmt.Errorf("unsupported manifest operator %q", operator)
	}
}

// All model limits and presets are supplied by the plugin. The evaluator only
// understands generic pixel/ratio notation, never vendor or model identities.
func resolveManifestPixelSize(spec map[string]any) (any, error) {
	value := strings.ToLower(strings.TrimSpace(manifestString(spec["value"])))
	if value == "" || value == "auto" {
		return nil, nil
	}
	value = strings.ReplaceAll(value, "×", "x")
	if preset := manifestString(manifestObject(spec["presets"])[value]); preset != "" {
		value = preset
	}
	var width, height int64
	if strings.Contains(value, ":") {
		parts := strings.Split(value, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("比例格式无效")
		}
		w, e1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		h, e2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if e1 != nil || e2 != nil || w <= 0 || h <= 0 || math.IsInf(w, 0) || math.IsInf(h, 0) || math.IsNaN(w) || math.IsNaN(h) {
			return nil, fmt.Errorf("比例必须是有限正数")
		}
		ratio := w / h
		if ratio < manifestFloat(spec["minRatio"]) || ratio > manifestFloat(spec["maxRatio"]) {
			return nil, fmt.Errorf("宽高比超出支持范围")
		}
		x := math.Sqrt(manifestFloat(spec["targetPixels"]) * ratio)
		y := math.Sqrt(manifestFloat(spec["targetPixels"]) / ratio)
		if side := manifestFloat(spec["ratioMaxSide"]); side > 0 {
			scale := math.Min(1, side/math.Max(x, y))
			x *= scale
			y *= scale
		}
		width, height = int64(math.Round(x)), int64(math.Round(y))
	} else {
		parts := strings.Split(strings.ReplaceAll(value, "*", "x"), "x")
		if len(parts) != 2 {
			return nil, fmt.Errorf("尺寸请使用宽x高或宽:高")
		}
		var e1, e2 error
		width, e1 = strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 32)
		height, e2 = strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 32)
		if e1 != nil || e2 != nil {
			return nil, fmt.Errorf("尺寸必须是整数")
		}
	}
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("宽高必须是正整数")
	}
	pixels := float64(width) * float64(height)
	ratio := float64(width) / float64(height)
	if pixels < manifestFloat(spec["minPixels"]) || pixels > manifestFloat(spec["maxPixels"]) {
		return nil, fmt.Errorf("总像素超出支持范围")
	}
	if ratio < manifestFloat(spec["minRatio"]) || ratio > manifestFloat(spec["maxRatio"]) {
		return nil, fmt.Errorf("宽高比超出支持范围")
	}
	return fmt.Sprintf("%d%s%d", width, manifestString(spec["separator"]), height), nil
}

func validateManifestMedia(spec map[string]any) (any, error) {
	items := manifestArray(spec["items"])
	if minimum, ok := spec["minCount"]; ok && len(items) < manifestInt(minimum) {
		return nil, fmt.Errorf("参考素材不足")
	}
	if maximum, ok := spec["maxCount"]; ok && len(items) > manifestInt(maximum) {
		return nil, fmt.Errorf("参考素材数量超过限制")
	}
	for _, raw := range items {
		item := manifestObject(raw)
		meta := manifestObject(item["metadata"])
		value := manifestString(item["value"])
		mime := strings.ToLower(strings.TrimSpace(strings.Split(manifestString(item["mimeType"]), ";")[0]))
		size := manifestFloat(meta["bytes"])
		if strings.HasPrefix(value, "data:") {
			if !manifestTruthy(spec["allowData"]) {
				return nil, fmt.Errorf("参考素材必须使用公网 URL")
			}
			header, encoded, ok := strings.Cut(value, ",")
			if !ok || !strings.HasSuffix(header, ";base64") {
				return nil, fmt.Errorf("无效 base64 data URL")
			}
			data, e := base64.StdEncoding.DecodeString(encoded)
			if e != nil || len(data) == 0 {
				return nil, fmt.Errorf("无效 base64 媒体")
			}
			size = float64(len(data))
			mime = strings.TrimSuffix(strings.TrimPrefix(header, "data:"), ";base64")
		} else {
			u, e := url.Parse(value)
			if e != nil || u.Hostname() == "" || u.User != nil {
				return nil, fmt.Errorf("参考素材地址无效")
			}
			allowedOrigin := false
			for _, origin := range manifestArray(spec["allowedOrigins"]) {
				allowedOrigin = allowedOrigin || u.Scheme+"://"+u.Host == manifestString(origin)
			}
			if !allowedOrigin {
				if u.Scheme != "https" && u.Scheme != "http" {
					return nil, fmt.Errorf("参考素材必须使用 HTTP(S) URL")
				}
				ip := net.ParseIP(u.Hostname())
				if strings.EqualFold(u.Hostname(), "localhost") || (ip != nil && (ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsUnspecified())) {
					return nil, fmt.Errorf("参考素材不能使用私网地址")
				}
			}
		}
		if max := manifestFloat(spec["maxBytes"]); max > 0 && size > max {
			return nil, fmt.Errorf("参考素材文件超过大小限制")
		}
		if allowed, ok := spec["mimes"]; ok && mime != "" {
			found := false
			for _, candidate := range manifestArray(allowed) {
				found = found || mime == manifestString(candidate)
			}
			if !found {
				return nil, fmt.Errorf("参考素材格式不受支持")
			}
		}
		w, h := manifestFloat(meta["width"]), manifestFloat(meta["height"])
		if w > 0 && h > 0 {
			if min := manifestFloat(spec["minSide"]); min > 0 && math.Min(w, h) < min {
				return nil, fmt.Errorf("参考素材尺寸过小")
			}
			if max := manifestFloat(spec["maxSide"]); max > 0 && math.Max(w, h) > max {
				return nil, fmt.Errorf("参考素材尺寸过大")
			}
			if max := manifestFloat(spec["maxAspect"]); max > 0 && math.Max(w, h)/math.Min(w, h) > max {
				return nil, fmt.Errorf("参考素材宽高比超出限制")
			}
		}
		duration := manifestFloat(meta["durationMs"])
		if duration > 0 {
			if min := manifestFloat(spec["minDurationMs"]); min > 0 && duration < min {
				return nil, fmt.Errorf("参考素材时长过短")
			}
			if max := manifestFloat(spec["maxDurationMs"]); max > 0 && duration > max {
				return nil, fmt.Errorf("参考素材时长过长")
			}
		}
	}
	return true, nil
}
