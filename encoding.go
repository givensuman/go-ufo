package ufo

import (
	"net/url"
	"strings"

	"golang.org/x/net/idna"
)

// Encode encodes characters using encodeURI semantics: everything except the
// standard unreserved set plus ; , / ? : @ & = + $ # ~ ! * ' ( ) and |.
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
		{"%7C", "|"}, // pipe must NOT be encoded per ufo
	} {
		encoded = strings.ReplaceAll(encoded, pair[0], pair[1])
	}
	return encoded
}

// Decode decodes text using url.PathUnescape. Returns the original text on failure.
func Decode(text string) string {
	result, err := url.PathUnescape(text)
	if err != nil {
		return text
	}
	return result
}

// DecodePath decodes path section of URL. %2F is preserved (not decoded to /).
func DecodePath(text string) string {
	// Protect encoded slashes before decoding
	protected := strings.ReplaceAll(text, "%2F", "%252F")
	protected = strings.ReplaceAll(protected, "%2f", "%252F")
	return Decode(protected)
}

// DecodeQueryKey decodes a query key, converting + to space first.
func DecodeQueryKey(text string) string {
	return Decode(strings.ReplaceAll(text, "+", " "))
}

// DecodeQueryValue decodes a query value, converting + to space first.
func DecodeQueryValue(text string) string {
	return Decode(strings.ReplaceAll(text, "+", " "))
}

// EncodeHash encodes characters for the hash section of a URL.
// Like Encode but { } ^ are preserved.
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

// EncodePath encodes characters for the path section of a URL.
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

// EncodeParam encodes like EncodePath but also encodes the slash character.
func EncodeParam(text string) string {
	return strings.ReplaceAll(EncodePath(text), "/", "%2F")
}

// EncodeQueryValue encodes characters for a query string value.
// Space is encoded as +, literal + is encoded as %2B.
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

// EncodeQueryKey encodes like EncodeQueryValue but also encodes =.
func EncodeQueryKey(text string) string {
	return strings.ReplaceAll(EncodeQueryValue(text), "=", "%3D")
}

// EncodeHost encodes a hostname using punycode (IDNA) for internationalized domains.
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
