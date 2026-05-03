package ufo

import (
	"strings"
)

// QueryObject is a map of query key to one or more string values.
type QueryObject map[string][]string

// ParseQuery parses and decodes a query string into a map. The input
// can be a query string with or without the leading "?".
//
// The `__proto__` and `constructor` keys are ignored to prevent
// prototype pollution.
//
// Examples:
//
//	ParseQuery("?foo=bar&baz=qux")
//	// QueryObject{ "foo": ["bar"], "baz": ["qux"] }
//
//	ParseQuery("tags=javascript&tags=web&tags=dev")
//	// QueryObject{ "tags": ["javascript", "web", "dev"] }
func ParseQuery(parametersString string) QueryObject {
	result := make(QueryObject)
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
		result[key] = append(result[key], value)
	}

	return result
}

// EncodeQueryItem encodes a key and its values into URL query string
// segments. Multiple values produce multiple key=value pairs joined
// with "&". A nil or empty slice produces only the key with no value.
//
// Examples:
//
//	EncodeQueryItem("message", []string{"Hello World"})
//	// "message=Hello+World"
//
//	EncodeQueryItem("tags", []string{"javascript", "web", "dev"})
//	// "tags=javascript&tags=web&tags=dev"
func EncodeQueryItem(key string, values []string) string {
	if len(values) == 0 {
		return EncodeQueryKey(key)
	}

	parts := make([]string, 0, len(values))
	for _, v := range values {
		parts = append(parts, EncodeQueryKey(key)+"="+EncodeQueryValue(v))
	}

	return strings.Join(parts, "&")
}

// String encodes the QueryObject into a URL query string without a
// leading "?". Keys with nil or empty value slices are omitted.
func (query QueryObject) String() string {
	parts := make([]string, 0, len(query))
	for k, vals := range query {
		if len(vals) == 0 {
			continue
		}

		item := EncodeQueryItem(k, vals)
		if item != "" {
			parts = append(parts, item)
		}
	}

	return strings.Join(parts, "&")
}
