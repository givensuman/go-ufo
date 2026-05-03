package ufo

import (
	"maps"
	"regexp"
	"strings"
)

var (
	protocolStrictRe   = regexp.MustCompile(`^[\s\w\x00+.-]{2,}:([/\\]{1,2})`)
	protocolRe         = regexp.MustCompile(`^[\s\w\x00+.-]{2,}:([/\\]{2})?`)
	protocolRelativeRe = regexp.MustCompile(`^([/\\]\s*){2,}[^/\\]`)
	protocolScriptRe   = regexp.MustCompile(`(?i)^[\s\x00]*(blob|data|javascript|vbscript):$`)
	trailingSlashRe    = regexp.MustCompile(`/$|/\?|/#`)
	joinLeadingSlashRe = regexp.MustCompile(`^\.?/`)
	doubleSlashRe      = regexp.MustCompile(`/{2,}`)
	leadDoubleSlashRe  = regexp.MustCompile(`^/{2,}`)
)

// IsRelative checks if a path starts with "./" or "../".
//
// Example:
//
//	IsRelative("./foo")  // true
func IsRelative(inputString string) bool {
	return strings.HasPrefix(inputString, "./") || strings.HasPrefix(inputString, "../")
}

// HasProtocol checks if the input has a protocol. You can
// use `opts.AcceptRelative = true` to accept relative URLs
// as valid protocol. Use `opts.Strict = true` to strictly check
// for RFC 3986 protocols.
//
// Examples:
//
//	HasProtocol("https://example.com", nil) // true
//	HasProtocol("//example.com", nil) // false
//
// The same examples, allowing relative URLs:
//
//	HasProtocol("https://example.com", &HasProtocolOptions{ AcceptRelative: true }) // true
//	HasProtocol("//example.com", &HasProtocolOptions{ AcceptRelative: true }) // true
//
// Example distinction of strict protocol check:
//
//	HasProtocol("data:text/plain", nil) // true
//	HasProtocol("data:text/plain", &HasProtocolOptions{ Strict: true }) // false
func HasProtocol(inputString string, opts *HasProtocolOptions) bool {
	if opts == nil {
		opts = &HasProtocolOptions{}
	}

	if opts.Strict {
		return protocolStrictRe.MatchString(inputString)
	}

	return protocolRe.MatchString(inputString) ||
		(opts.AcceptRelative && protocolRelativeRe.MatchString(inputString))
}

// HasProtocolOptions configures HasProtocol.
type HasProtocolOptions struct {
	// AcceptRelative treats relative URLs (e.g. "//foo.com") as having a protocol.
	AcceptRelative bool
	// Strict requires RFC 3986 protocol syntax (two slashes after colon).
	Strict bool
}

// IsScriptProtocol checks if the input protocol is any of
// the dangerous "blob:", "data:", "javascript:" or
// "vbscript:" protocols.
//
// Examples:
//
//	IsScriptProtocol("javascript:alert(1)") // true
//	IsScriptProtocol("data:text/html,hello") // true
//	IsScriptProtocol("blob:hello") // true
//	IsScriptProtocol("vbscript:alert(1)") // true
//	IsScriptProtocol("https://examples.com") // false
func IsScriptProtocol(protocol string) bool {
	return protocol != "" && protocolScriptRe.MatchString(protocol)
}

// RespectQueryAndFragmentOption controls whether a slash immediately
// before a query string or fragment is treated as a trailing slash.
type RespectQueryAndFragmentOption struct {
	// Count a slash before a query string or hash as a trailing slash.
	RespectQueryAndFragment bool
}

var dontRespect = &RespectQueryAndFragmentOption{false}

func (r *RespectQueryAndFragmentOption) isTrue() bool {
	return r != nil && r.RespectQueryAndFragment
}

// HasTrailingSlash checks if input has a trailing slash. You
// can use the `respectQueryAndFragment` parameter to
// count a slash before a query string or hash as a trailing slash.
//
// Examples:
//
//	HasTrailingSlash("/foo/", nil) // true
//	HasTrailingSlash("/foo", nil) // false
//	HasTrailingSlash("/foo?query=true", &RespectQueryAndFragmentOption{ true }) // false
//	HasTrailingSlash("/foo/?query=true", &RespectQueryAndFragmentOption{ true }) // true
func HasTrailingSlash(input string, respectQueryAndFragment *RespectQueryAndFragmentOption) bool {
	if input == "" {
		return false
	}

	if respectQueryAndFragment.isTrue() {
		return strings.HasSuffix(input, "/")
	}

	return trailingSlashRe.MatchString(input)
}

// WithoutTrailingSlash removes the trailing slash from a URL.
//
// Examples:
//
//	WithoutTrailingSlash("/foo/", nil) // "/foo"
//	WithoutTrailingSlash("/path/?query=true", &RespectQueryAndFragmentOption{ true }) // "/path?query=true"
func WithoutTrailingSlash(input string, respectQueryAndFragment *RespectQueryAndFragmentOption) string {
	if input == "" {
		return "/"
	}

	if !respectQueryAndFragment.isTrue() {
		if HasTrailingSlash(input, dontRespect) {
			r := input[:len(input)-1]
			if r == "" {
				return "/"
			}

			return r
		}

		return input
	}

	if !HasTrailingSlash(input, dontRespect) {
		return input
	}

	path := input
	fragment := ""
	if idx := strings.Index(input, "#"); idx != -1 {
		path = input[:idx]
		fragment = input[idx:]
	}

	parts := strings.SplitN(path, "?", 2)
	cleanPath := strings.TrimSuffix(parts[0], "/")

	if cleanPath == "" {
		cleanPath = "/"
	}

	query := ""
	if len(parts) > 1 {
		query = "?" + parts[1]
	}

	return cleanPath + query + fragment
}

// WithTrailingSlash ensures the URL ends with a trailing slash.
//
// Examples:
//
//	WithTrailingSlash("/foo", nil) // "/foo/"
//	WithTrailingSlash("/path?query=true", &RespectQueryAndFragmentOption{ true }) // "/path/?query=true"
func WithTrailingSlash(input string, respectQueryAndFragment *RespectQueryAndFragmentOption) string {
	if input == "" {
		return "/"
	}

	if !respectQueryAndFragment.isTrue() {
		if strings.HasSuffix(input, "/") {
			return input
		}

		return input + "/"
	}

	if HasTrailingSlash(input, dontRespect) {
		return input
	}

	path := input
	fragment := ""
	if idx := strings.Index(input, "#"); idx != -1 {
		path = input[:idx]
		fragment = input[idx:]

		if path == "" {
			return fragment
		}
	}

	parts := strings.SplitN(path, "?", 2)
	query := ""
	if len(parts) > 1 {
		query = "?" + parts[1]
	}

	return parts[0] + "/" + query + fragment
}

// HasLeadingSlash checks if the input has a leading slash.
//
// Example:
//
//	HasLeadingSlash("/foo") // true
func HasLeadingSlash(input string) bool {
	return strings.HasPrefix(input, "/")
}

// WithoutLeadingSlash removes the leading slash
// from the URL or pathname.
//
// Example:
//
//	WithoutLeadingSlash("/foo") // "foo"
func WithoutLeadingSlash(input string) string {
	if HasLeadingSlash(input) {
		return input[1:]
	}

	return input
}

// WithLeadingSlash ensures the URL or pathname
// has a leading slash.
//
// Example:
//
//	WithLeadingSlash("foo") // "/foo"
func WithLeadingSlash(input string) string {
	if HasLeadingSlash(input) {
		return input
	}

	return "/" + input
}

// CleanDoubleSlashes removes double slashes from the URL.
//
// Examples:
//
//	CleanDoubleSlashes("//foo//bar//") // "/foo/bar/"
//	CleanDoubleSlashes("http://example.com/analyze//http://localhost:3000/")
//	// Returns "http://example.com/analyze/http://localhost:3000/"
func CleanDoubleSlashes(input string) string {
	parts := strings.Split(input, "://")
	for i, part := range parts {
		parts[i] = doubleSlashRe.ReplaceAllString(part, "/")
	}

	return strings.Join(parts, "://")
}

// WithBase ensures the URL or pathname starts with base.
//
// Examples:
//
//	WithBase("/foo/bar", "/foo") // "/foo/bar"
//	WithBase("/foo/bar", "/baz") // "/baz/foo/bar"
func WithBase(input string, base string) string {
	if IsEmptyURL(base) || HasProtocol(input, nil) {
		return input
	}

	b := WithoutTrailingSlash(base, dontRespect)
	if strings.HasPrefix(input, b) {
		next := ""
		if len(input) > len(b) {
			next = string(input[len(b)])
		}

		if next == "" || next == "/" || next == "?" {
			return input
		}
	}

	return JoinURL(b, input)
}

// WithoutBase removes the base from the URL or pathname.
//
// Example:
//
//	WithoutBase("/foo/bar", "/foo") // "/bar"
func WithoutBase(input string, base string) string {
	if IsEmptyURL(base) {
		return input
	}

	b := WithoutTrailingSlash(base, dontRespect)
	if !strings.HasPrefix(input, b) {
		return input
	}

	next := ""
	if len(input) > len(b) {
		next = string(input[len(b)])
	}

	if next != "" && next != "/" && next != "?" {
		return input
	}

	trimmed := input[len(b):]
	if trimmed == "" || trimmed[0] != '/' {
		return "/" + trimmed
	}

	return trimmed
}

// WithQuery adds or replaces query params of the URL.
//
// Example:
//
//	WithQuery("/foo?page=a", QueryObject{ "token": []string{"secret"} })
//	// Returns "/foo?page=a&token=secret"
func WithQuery(input string, query QueryObject) string {
	parsed := ParseURL(input, "")
	existing := ParseQuery(parsed.Search)
	maps.Copy(existing, query)

	qs := existing.String()
	if qs != "" {
		parsed.Search = "?" + qs
	} else {
		parsed.Search = ""
	}

	return parsed.String()
}

// FilterQuery filters the query section of the URL using a callback.
// Keys for which the callback returns false are removed.
//
// Example:
//
//	FilterQuery("/foo?bar=1&baz=2", func(k string, v []string) bool {
//		return k != "bar"
//	}) // "/foo?baz=2"
func FilterQuery(input string, callback func(k string, v []string) bool) string {
	if !strings.Contains(input, "?") {
		return input
	}

	parsed := ParseURL(input, "")
	q := ParseQuery(parsed.Search)
	filtered := make(QueryObject)

	for k, v := range q {
		if callback(k, v) {
			filtered[k] = v
		}
	}

	qs := filtered.String()
	if qs != "" {
		parsed.Search = "?" + qs
	} else {
		parsed.Search = ""
	}

	return parsed.String()
}

// GetQuery parses and decodes the query object of an input URL.
//
// Example:
//
//	GetQuery("http://foo.com/foo?test=123&unicode=%E5%A5%BD")
//	// Returns QueryObject{ "test": ["123"], "unicode": ["好"] }
func GetQuery(input string) QueryObject {
	return ParseQuery(ParseURL(input, "").Search)
}

// IsEmptyURL checks if the input URL is empty or "/".
func IsEmptyURL(url string) bool {
	return url == "" || url == "/"
}

// IsNonEmptyURL checks if the URL is neither empty nor "/".
func IsNonEmptyURL(url string) bool {
	return url != "" && url != "/"
}

// JoinURL joins multiple URL segments into a single URL.
//
// Example:
//
//	JoinURL("a", "/b", "/c") // "a/b/c"
func JoinURL(base string, input ...string) string {
	url := base
	for _, segment := range input {
		if !IsNonEmptyURL(segment) {
			continue
		}

		if url != "" {
			seg := joinLeadingSlashRe.ReplaceAllString(segment, "")
			url = WithTrailingSlash(url, dontRespect) + seg
		} else {
			url = segment
		}
	}

	return url
}

func splitURLPath(path string) []string {
	var result []string
	start := 0

	for j := 0; j < len(path); j++ {
		if path[j] == '/' {
			// skip if next char is also '/'
			if j+1 < len(path) && path[j+1] == '/' {
				continue
			}
			result = append(result, path[start:j])
			start = j + 1
		}
	}

	result = append(result, path[start:])
	return result
}

// JoinRelativeURL joins multiple URL segments into
// a single URL and also handles relative paths with "./" and "../"
//
// Example:
//
//	JoinRelativeURL("/a", "../b", "./c") // "/b/c"
func JoinRelativeURL(input ...string) string {
	var segments []string
	segmentsDepth := 0

	var filtered []string
	for _, s := range input {
		if s != "" {
			filtered = append(filtered, s)
		}
	}

	for _, i := range filtered {
		if i == "" || i == "/" {
			continue
		}

		for sindex, s := range splitURLPath(i) {
			if s == "" || s == "." {
				continue
			}

			if s == ".." {
				if len(segments) == 1 && HasProtocol(segments[0], &HasProtocolOptions{}) {
					continue
				}

				if len(segments) > 0 {
					segments = segments[:len(segments)-1]
				}

				segmentsDepth--
				continue
			}

			if sindex == 1 && len(segments) > 0 && strings.HasSuffix(segments[len(segments)-1], ":/") {
				segments[len(segments)-1] += "/" + s
				continue
			}

			segments = append(segments, s)
			segmentsDepth++
		}
	}

	url := strings.Join(segments, "/")
	if segmentsDepth >= 0 {
		if len(filtered) > 0 && strings.HasPrefix(filtered[0], "/") && !strings.HasPrefix(url, "/") {
			url = "/" + url
		} else if len(filtered) > 0 && strings.HasPrefix(filtered[0], "./") && !strings.HasPrefix(url, "./") {
			url = "./" + url
		}
	} else {
		prefix := strings.Repeat("../", -1*segmentsDepth)
		url = prefix + url
	}

	if len(filtered) > 0 && strings.HasSuffix(filtered[len(filtered)-1], "/") && !strings.HasSuffix(url, "/") {
		url += "/"
	}

	return url
}

// WithHTTP adds or replaces the URL protocol with "http://".
//
// Example:
//
//	WithHTTP("https://example.com") // http://example.com
func WithHTTP(input string) string {
	return WithProtocol(input, "http://")
}

// WithHTTPS adds or replaces the URL protocol with "https://".
//
// Example:
//
//	WithHTTPS("http://example.com") // https://example.com
func WithHTTPS(input string) string {
	return WithProtocol(input, "https://")
}

// WithoutProtocol removes the protocol from the input.
//
// Example:
//
//	WithoutProtocol("http://example.com") // example.com
func WithoutProtocol(input string) string {
	return WithProtocol(input, "")
}

// WithProtocol adds or replaces the protocol from the input.
//
// Example:
//
//	WithProtocol("http://example.com", "ftp://") // ftp://example.com
func WithProtocol(input string, protocol string) string {
	match := protocolRe.FindString(input)
	if match == "" {
		match = leadDoubleSlashRe.FindString(input)
	}

	if match == "" {
		return protocol + input
	}

	return protocol + input[len(match):]
}

// NormalizeURL normalizes the URL:
//   - Ensures the URL is properly encoded
//   - Ensures pathname starts with a slash
//   - Preserves protocol/host if provided
//
// Example:
//
//	NormalizeURL("test?query=123 123#hash, test")
//	// Returns "test?query=123%20123#hash,%20test"
func NormalizeURL(input string) string {
	parsed := ParseURL(input, "")
	parsed.Pathname = EncodePath(DecodePath(parsed.Pathname))

	// Hash includes the leading "#"; encode only the content after it
	if strings.HasPrefix(parsed.Hash, "#") {
		parsed.Hash = "#" + EncodeHash(Decode(parsed.Hash[1:]))
	} else {
		parsed.Hash = EncodeHash(Decode(parsed.Hash))
	}

	parsed.Host = EncodeHost(Decode(parsed.Host))
	qs := ParseQuery(parsed.Search).String()

	// NormalizeURL uses %20 encoding rather than + for spaces
	qs = strings.ReplaceAll(qs, "+", "%20")
	if qs != "" {
		parsed.Search = "?" + qs
	} else {
		parsed.Search = ""
	}

	return parsed.String()
}

// ResolveURL resolves multiple URL segments into a single URL.
//
// Example:
//
//	ResolveURL("http://foo.com/foo?test=123#token", "bar", "baz")
//	// Returns "http://foo.com/foo/bar/baz?test=123#token"
func ResolveURL(base string, inputs ...string) string {
	var filtered []string
	for _, inp := range inputs {
		if IsNonEmptyURL(inp) {
			filtered = append(filtered, inp)
		}
	}

	if len(filtered) == 0 {
		return base
	}

	url := ParseURL(base, "")
	for _, segment := range filtered {
		seg := ParseURL(segment, "")

		// Append path
		if seg.Pathname != "" {
			url.Pathname = WithTrailingSlash(
				url.Pathname,
				dontRespect,
			) + WithoutLeadingSlash(seg.Pathname)
		}

		// Override hash
		if seg.Hash != "" && seg.Hash != "#" {
			url.Hash = seg.Hash
		}

		// Append search
		if seg.Search != "" && seg.Search != "?" {
			if url.Search != "" && url.Search != "?" {
				m := make(QueryObject)
				maps.Copy(m, ParseQuery(url.Search))
				maps.Copy(m, ParseQuery(seg.Search))
				merged := m.String()

				if merged != "" {
					url.Search = "?" + merged
				} else {
					url.Search = ""
				}
			} else {
				url.Search = seg.Search
			}
		}
	}

	return url.String()
}

// IsSamePath checks if two paths are equal or not.
// Trailing slash and encoding are normalized before comparison.
//
// Example:
//
//	IsSamePath("/foo", "/foo/") // true
func IsSamePath(a string, b string) bool {
	return Decode(WithoutTrailingSlash(a, dontRespect)) == Decode(WithoutTrailingSlash(b, dontRespect))
}

// IsEqual checks if two paths are equal regardless of encoding,
// trailing slash, and leading slash differences. You can make
// encoding, trailing slash, or leading slash checks strict by
// setting the `opts` parameter.
//
// Examples:
//
//	IsEqual("/foo", "foo", nil) // true
//	IsEqual("foo/", "foo", nil) // true
//	IsEqual("/foo bar", "/foo%20bar", nil) // true
//
// The same examples, with strict checks:
//
//	IsEqual("/foo", "foo", &IsEqualOptions{ StrictLeadingSlash: true }) // false
//	IsEqual("foo/", "foo", &IsEqualOptions{ StrictTrailingSlash: true }) // false
//	IsEqual("/foo bar", "/foo%20bar", &IsEqualOptions{ StrictEncoding: true }) // false
func IsEqual(a string, b string, opts *IsEqualOptions) bool {
	if opts == nil {
		opts = &IsEqualOptions{}
	}

	if !opts.StrictTrailingSlash {
		a = WithTrailingSlash(a, dontRespect)
		b = WithTrailingSlash(b, dontRespect)
	}
	if !opts.StrictLeadingSlash {
		a = WithLeadingSlash(a)
		b = WithLeadingSlash(b)
	}
	if !opts.StrictEncoding {
		a = Decode(a)
		b = Decode(b)
	}

	return a == b
}

// IsEqualOptions configures IsEqual.
type IsEqualOptions struct {
	// StrictTrailingSlash disables trailing-slash normalisation.
	StrictTrailingSlash bool
	// StrictLeadingSlash disables leading-slash normalisation.
	StrictLeadingSlash bool
	// StrictEncoding disables percent-encoding normalisation.
	StrictEncoding bool
}

// WithFragment adds or replaces the fragment section
// of the URL.
//
// Examples:
//
//	WithFragment("/foo", "bar") // "/foo#bar"
//	WithFragment("/foo#bar", "baz") // "/foo#baz"
//	WithFragment("/foo#bar", "") // "/foo"
func WithFragment(input string, hash string) string {
	if hash == "#" {
		return input
	}

	parsed := ParseURL(input, "")
	if hash == "" {
		parsed.Hash = ""
	} else {
		parsed.Hash = "#" + EncodeHash(hash)
	}

	return parsed.String()
}

// WithoutFragment removes the fragment section from the URL.
//
// Example:
//
//	WithoutFragment("http://example.com/foo?q=123#bar")
//	// Returns "http://example.com/foo?q=123"
func WithoutFragment(input string) string {
	parsed := ParseURL(input, "")
	parsed.Hash = ""

	return parsed.String()
}

// WithoutHost removes the host from the URL while
// preserving everything else.
//
// Example:
//
//	WithoutHost("http://example.com/foo?q=123#bar")
//	// Returns "/foo?q=123#bar"
func WithoutHost(input string) string {
	parsed := ParseURL(input, "")
	pathname := parsed.Pathname
	if pathname == "" {
		pathname = "/"
	}

	return pathname + parsed.Search + parsed.Hash
}
