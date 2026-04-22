package ufo

import (
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
)

// IsRelative returns true if the input starts with "./" or "../".
func IsRelative(inputString string) bool {
	return strings.HasPrefix(inputString, "./") || strings.HasPrefix(inputString, "../")
}

// HasProtocol checks if the input has a URL protocol.
func HasProtocol(inputString string, opts HasProtocolOptions) bool {
	if opts.Strict {
		return protocolStrictRe.MatchString(inputString)
	}
	return protocolRe.MatchString(inputString) ||
		(opts.AcceptRelative && protocolRelativeRe.MatchString(inputString))
}

// IsScriptProtocol checks if the protocol is one of the dangerous script protocols.
func IsScriptProtocol(protocol string) bool {
	return protocol != "" && protocolScriptRe.MatchString(protocol)
}

// HasTrailingSlash checks if input has a trailing slash.
// If respectQueryAndFragment is true, a slash before ? or # also counts.
func HasTrailingSlash(input string, respectQueryAndFragment bool) bool {
	if input == "" {
		return false
	}
	if !respectQueryAndFragment {
		return strings.HasSuffix(input, "/")
	}
	return trailingSlashRe.MatchString(input)
}

// WithoutTrailingSlash removes the trailing slash from a URL.
func WithoutTrailingSlash(input string, respectQueryAndFragment bool) string {
	if input == "" {
		return "/"
	}
	if !respectQueryAndFragment {
		if HasTrailingSlash(input, false) {
			r := input[:len(input)-1]
			if r == "" {
				return "/"
			}
			return r
		}
		return input
	}
	if !HasTrailingSlash(input, true) {
		return input
	}
	path := input
	fragment := ""
	if idx := strings.Index(input, "#"); idx != -1 {
		path = input[:idx]
		fragment = input[idx:]
	}
	parts := strings.SplitN(path, "?", 2)
	cleanPath := parts[0]
	if strings.HasSuffix(cleanPath, "/") {
		cleanPath = cleanPath[:len(cleanPath)-1]
	}
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
func WithTrailingSlash(input string, respectQueryAndFragment bool) string {
	if input == "" {
		return "/"
	}
	if !respectQueryAndFragment {
		if strings.HasSuffix(input, "/") {
			return input
		}
		return input + "/"
	}
	if HasTrailingSlash(input, true) {
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

// HasLeadingSlash checks if input starts with "/".
func HasLeadingSlash(input string) bool {
	return strings.HasPrefix(input, "/")
}

// WithoutLeadingSlash removes the leading slash from input.
func WithoutLeadingSlash(input string) string {
	if HasLeadingSlash(input) {
		return input[1:]
	}
	return input
}

// WithLeadingSlash ensures input starts with "/".
func WithLeadingSlash(input string) string {
	if HasLeadingSlash(input) {
		return input
	}
	return "/" + input
}

// CleanDoubleSlashes removes consecutive slashes within a URL (preserving ://).
func CleanDoubleSlashes(input string) string {
	re := regexp.MustCompile(`/{2,}`)
	parts := strings.Split(input, "://")
	for i, part := range parts {
		parts[i] = re.ReplaceAllString(part, "/")
	}
	return strings.Join(parts, "://")
}

// WithBase ensures input starts with base (prefix).
func WithBase(input string, base string) string {
	if IsEmptyURL(base) || HasProtocol(input, HasProtocolOptions{}) {
		return input
	}
	b := WithoutTrailingSlash(base, false)
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

// WithoutBase removes base prefix from input.
func WithoutBase(input string, base string) string {
	if IsEmptyURL(base) {
		return input
	}
	b := WithoutTrailingSlash(base, false)
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

// WithQuery adds or replaces query params on a URL.
func WithQuery(input string, query QueryObject) string {
	parsed := ParseURL(input, "")
	existing := ParseQuery(parsed.Search)
	for k, v := range query {
		existing[k] = v
	}
	qs := StringifyQuery(QueryObject(existing))
	if qs != "" {
		parsed.Search = "?" + qs
	} else {
		parsed.Search = ""
	}
	return StringifyParsedURL(parsed)
}

// FilterQuery filters query parameters by a predicate.
func FilterQuery(input string, predicate func(key string, value interface{}) bool) string {
	if !strings.Contains(input, "?") {
		return input
	}
	parsed := ParseURL(input, "")
	q := ParseQuery(parsed.Search)
	filtered := make(QueryObject)
	for k, v := range q {
		if predicate(k, v) {
			filtered[k] = v
		}
	}
	qs := StringifyQuery(filtered)
	if qs != "" {
		parsed.Search = "?" + qs
	} else {
		parsed.Search = ""
	}
	return StringifyParsedURL(parsed)
}

// GetQuery parses and returns the query object of a URL.
func GetQuery(input string) map[string]interface{} {
	return ParseQuery(ParseURL(input, "").Search)
}

// IsEmptyURL checks if the URL is empty or just "/".
func IsEmptyURL(url string) bool {
	return url == "" || url == "/"
}

// IsNonEmptyURL checks if the URL is neither empty nor "/".
func IsNonEmptyURL(url string) bool {
	return !IsEmptyURL(url)
}

// JoinURL joins multiple URL segments.
func JoinURL(base string, input ...string) string {
	url := base
	for _, segment := range input {
		if !IsNonEmptyURL(segment) {
			continue
		}
		if url != "" {
			seg := joinLeadingSlashRe.ReplaceAllString(segment, "")
			url = WithTrailingSlash(url, false) + seg
		} else {
			url = segment
		}
	}
	return url
}

// splitURLPath splits a path on single slashes (not double slashes).
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

// JoinRelativeURL joins URL segments and resolves "." and ".." components.
func JoinRelativeURL(input ...string) string {
	segments := []string{}
	segmentsDepth := 0

	filtered := []string{}
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
				if len(segments) == 1 && HasProtocol(segments[0], HasProtocolOptions{}) {
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
func WithHTTP(input string) string {
	return WithProtocol(input, "http://")
}

// WithHTTPS adds or replaces the URL protocol with "https://".
func WithHTTPS(input string) string {
	return WithProtocol(input, "https://")
}

// WithoutProtocol removes the protocol from the URL.
func WithoutProtocol(input string) string {
	return WithProtocol(input, "")
}

// WithProtocol adds or replaces the URL protocol.
func WithProtocol(input string, protocol string) string {
	match := protocolRe.FindString(input)
	if match == "" {
		leadRe := regexp.MustCompile(`^/{2,}`)
		match = leadRe.FindString(input)
	}
	if match == "" {
		return protocol + input
	}
	return protocol + input[len(match):]
}

// NormalizeURL normalises a URL (encode pathname, hash, host; re-stringify query).
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
	qs := StringifyQuery(QueryObject(ParseQuery(parsed.Search)))
	// NormalizeURL uses %20 encoding rather than + for spaces
	qs = strings.ReplaceAll(qs, "+", "%20")
	if qs != "" {
		parsed.Search = "?" + qs
	} else {
		parsed.Search = ""
	}
	return StringifyParsedURL(parsed)
}

// ResolveURL resolves URL segments onto a base URL.
func ResolveURL(base string, inputs ...string) string {
	filtered := []string{}
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
		if seg.Pathname != "" {
			url.Pathname = WithTrailingSlash(url.Pathname, false) + WithoutLeadingSlash(seg.Pathname)
		}
		if seg.Hash != "" && seg.Hash != "#" {
			url.Hash = seg.Hash
		}
		if seg.Search != "" && seg.Search != "?" {
			if url.Search != "" && url.Search != "?" {
				merged := StringifyQuery(func() QueryObject {
					m := QueryObject{}
					for k, v := range ParseQuery(url.Search) {
						m[k] = v
					}
					for k, v := range ParseQuery(seg.Search) {
						m[k] = v
					}
					return m
				}())
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
	return StringifyParsedURL(url)
}

// IsSamePath checks if two paths are equal after decoding and removing trailing slash.
func IsSamePath(p1 string, p2 string) bool {
	return Decode(WithoutTrailingSlash(p1, false)) == Decode(WithoutTrailingSlash(p2, false))
}

// IsEqual checks if two URLs are equal, with configurable strictness.
func IsEqual(a string, b string, options CompareURLOptions) bool {
	if !options.TrailingSlash {
		a = WithTrailingSlash(a, false)
		b = WithTrailingSlash(b, false)
	}
	if !options.LeadingSlash {
		a = WithLeadingSlash(a)
		b = WithLeadingSlash(b)
	}
	if !options.Encoding {
		a = Decode(a)
		b = Decode(b)
	}
	return a == b
}

// WithFragment adds or replaces the hash fragment of a URL.
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
	return StringifyParsedURL(parsed)
}

// WithoutFragment removes the hash fragment from a URL.
func WithoutFragment(input string) string {
	parsed := ParseURL(input, "")
	parsed.Hash = ""
	return StringifyParsedURL(parsed)
}

// WithoutHost removes the host from a URL, preserving path/query/hash.
func WithoutHost(input string) string {
	parsed := ParseURL(input, "")
	pathname := parsed.Pathname
	if pathname == "" {
		pathname = "/"
	}
	return pathname + parsed.Search + parsed.Hash
}
