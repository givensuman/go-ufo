package ufo

import (
	"net/url"
	"strings"

	"golang.org/x/net/idna"
)

// Encode encodes characters that need to be encoded in the path, search
// and hash sections of the URL.
//
// Attempts to follow the behavior of JavaScript's `encodeURI`.
func Encode(text string) string {
	encoded := url.PathEscape(text)

	// url.PathEscape encodes some chars that encodeURI does NOT encode. Restore them.
	for _, pair := range [][2]string{
		{"%3B", ";"},
		{"%2C", ","},
		{"%3F", "?"},
		{"%23", "#"},
		{"%40", "@"},
		{"%26", "&"},
		{"%3D", "="},
		{"%2B", "+"},
		{"%24", "$"},
		{"%21", "!"},
		{"%27", "'"},
		{"%28", "("},
		{"%29", ")"},
		{"%2A", "*"},
		{"%7C", "|"},
	} {
		encoded = strings.ReplaceAll(encoded, pair[0], pair[1])
	}

	return encoded
}

// Decode decodes text, returning the original text if
// decoding fails.
func Decode(text string) string {
	result, err := url.PathUnescape(text)
	if err != nil {
		return text
	}

	return result
}

// DecodePath decodes path section of the URL, consistent with
// `EncodePath`.
func DecodePath(text string) string {
	// protect encoded slashes before decoding
	protected := strings.ReplaceAll(text, "%2F", "%252F")
	protected = strings.ReplaceAll(protected, "%2f", "%252F")
	return Decode(protected)
}

// DecodeQueryKey decodes query key, consistent with
// `EncodeQueryKey`
func DecodeQueryKey(text string) string {
	return Decode(strings.ReplaceAll(text, "+", " "))
}

// DecodeQueryValue decodes query value, consistent with
// `EncodeQueryValue`
func DecodeQueryValue(text string) string {
	return Decode(strings.ReplaceAll(text, "+", " "))
}

// EncodeHash encodes characters that need to be encoded in the
// hash section of the URL.
func EncodeHash(text string) string {
	encoded := Encode(text)

	encoded = strings.ReplaceAll(encoded, "%7B", "{")
	encoded = strings.ReplaceAll(encoded, "%7b", "{")
	encoded = strings.ReplaceAll(encoded, "%7D", "}")
	encoded = strings.ReplaceAll(encoded, "%7d", "}")
	encoded = strings.ReplaceAll(encoded, "%5E", "^")
	encoded = strings.ReplaceAll(encoded, "%5e", "^")

	return encoded
}

// EncodePath encodes characters that need to be encoded in the path
// section of the URL.
func EncodePath(text string) string {
	encoded := Encode(text)

	encoded = strings.ReplaceAll(encoded, "#", "%23")
	encoded = strings.ReplaceAll(encoded, "?", "%3F")
	// preserve double-encoded slashes
	encoded = strings.ReplaceAll(encoded, "%252f", "%2F")
	encoded = strings.ReplaceAll(encoded, "&", "%26")
	encoded = strings.ReplaceAll(encoded, "+", "%2B")

	return encoded
}

// EncodeParam encodes characters that need to be encoded in the
// path section of the URL and also the "/" character.
func EncodeParam(text string) string {
	return strings.ReplaceAll(EncodePath(text), "/", "%2F")
}

// EncodeQueryValue encodes characters that need to be encoded
// for query values in the query section of the URL.
func EncodeQueryValue(input string) string {
	encoded := Encode(input)

	encoded = strings.ReplaceAll(encoded, "+", "%2B")
	encoded = strings.ReplaceAll(encoded, "%20", "+")
	encoded = strings.ReplaceAll(encoded, "#", "%23")
	encoded = strings.ReplaceAll(encoded, "&", "%26")
	encoded = strings.ReplaceAll(encoded, "%60", "`")
	encoded = strings.ReplaceAll(encoded, "%5E", "^")
	encoded = strings.ReplaceAll(encoded, "%5e", "^")
	encoded = strings.ReplaceAll(encoded, "/", "%2F")

	return encoded
}

// EncodeQueryKey encodes characters that need to be encoded
// for query values in the query section of the URL and also
// encodes the "=" character.
func EncodeQueryKey(text string) string {
	return strings.ReplaceAll(EncodeQueryValue(text), "=", "%3D")
}

// EncodeHost encodes hostname with punycode encoding, returning
// the original name if encoding fails.
func EncodeHost(name string) string {
	if name == "" {
		return ""
	}

	result, err := idna.Lookup.ToASCII(name)
	if err != nil {
		return name
	}

	return result
}
