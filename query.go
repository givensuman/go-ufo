package ufo

import (
	"fmt"
	"strings"
)

// ParseQuery parses and decodes a query string into
// a map. The input can be a query string with or
// without the leading "?".
//
// The `__proto__` and `constructor` keys are ignored
// to prevent prototype pollution.
//
// Examples:
//
//	ParseQuery("?foo=bar&baz=qux")
//	// map[string]any{ foo: "bar", baz: "qux" }
//
//	ParseQuery("tags=javascript&tags=web&tags=dev")
//	// map[string]any{ tags: ["javascript", "web", "dev"]}
func ParseQuery(parametersString string) map[string]any {
	result := make(map[string]any)
	if parametersString == "" {
		return result
	}

	if parametersString[0] == '?' {
		parametersString = parametersString[1:]
	}

	for param := range strings.SplitSeq(parametersString, "&") {
		if param == "" {
			continue
		}

		var rawKey, rawVal string
		rawKey, rawVal, found := strings.Cut(param, "=")

		if !found {
			rawKey = param
			rawVal = ""
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

// EncodeQueryItem encodes a key-value pair into
// a URL query string value. If the value is an array,
// it will be encoded as multiple key-value pairs with
// the same key.
//
// Examples:
//
//	EncodeQueryItem("message", "Hello World")
//	// "message=Hello+World"
//
//	EncodeQueryItem("tags", "[javascript", "web", "dev"])
//	// "tags=javascript&tags=web&tags=dev"
func EncodeQueryItem(key string, value any) string {
	if value == nil {
		return EncodeQueryKey(key)
	}

	switch v := value.(type) {
	case []any:
		parts := make([]string, 0, len(v))
		for _, item := range v {
			parts = append(parts, EncodeQueryKey(key)+"="+EncodeQueryValue(fmt.Sprintf("%v", item)))
		}

		return strings.Join(parts, "&")

	default:
		return EncodeQueryKey(key) + "=" + EncodeQueryValue(fmt.Sprintf("%v", v))
	}
}

func (query QueryObject) String() string {
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
