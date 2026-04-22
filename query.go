package ufo

import (
	"fmt"
	"strings"
)

// ParseQuery parses a query string (with or without leading "?") into a map.
// Repeated keys produce a []string value. The keys "__proto__" and "constructor"
// are silently dropped.
func ParseQuery(parametersString string) map[string]interface{} {
	result := make(map[string]interface{})
	if parametersString == "" {
		return result
	}
	if parametersString[0] == '?' {
		parametersString = parametersString[1:]
	}
	for _, param := range strings.Split(parametersString, "&") {
		if param == "" {
			continue
		}
		eqIdx := strings.Index(param, "=")
		var rawKey, rawVal string
		if eqIdx == -1 {
			rawKey = param
			rawVal = ""
		} else {
			rawKey = param[:eqIdx]
			rawVal = param[eqIdx+1:]
		}
		key := DecodeQueryKey(rawKey)
		if key == "__proto__" || key == "constructor" {
			continue
		}
		value := DecodeQueryValue(rawVal)
		existing, exists := result[key]
		if !exists {
			result[key] = value
		} else if arr, ok := existing.([]string); ok {
			result[key] = append(arr, value)
		} else {
			result[key] = []string{existing.(string), value}
		}
	}
	return result
}

// EncodeQueryItem encodes a key-value pair for use in a query string.
// If value is a []interface{}, it emits multiple key=value pairs joined by "&".
// If value is nil, it emits just the key.
func EncodeQueryItem(key string, value interface{}) string {
	if value == nil {
		return EncodeQueryKey(key)
	}
	switch v := value.(type) {
	case []interface{}:
		parts := make([]string, 0, len(v))
		for _, item := range v {
			parts = append(parts, EncodeQueryKey(key)+"="+EncodeQueryValue(fmt.Sprintf("%v", item)))
		}
		return strings.Join(parts, "&")
	default:
		return EncodeQueryKey(key) + "=" + EncodeQueryValue(fmt.Sprintf("%v", v))
	}
}

// StringifyQuery encodes a QueryObject into a query string (without leading "?").
// Keys with nil values are omitted.
func StringifyQuery(query QueryObject) string {
	parts := make([]string, 0, len(query))
	for k, v := range query {
		if v == nil {
			continue
		}
		item := EncodeQueryItem(k, v)
		if item != "" {
			parts = append(parts, item)
		}
	}
	return strings.Join(parts, "&")
}
