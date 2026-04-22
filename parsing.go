package ufo

import (
	"regexp"
	"strings"
)

var (
	specialProtoRe   = regexp.MustCompile(`(?i)^[\s\x00]*(blob:|data:|javascript:|vbscript:)(.*)`)
	hasProtocolRe    = regexp.MustCompile(`^[\s\w\x00+.-]{2,}:([/\\]{2})?`)
	protocolRelRe    = regexp.MustCompile(`^([/\\]\s*){2,}[^/\\]`)
	fullURLRe        = regexp.MustCompile(`(?i)^[\s\x00]*([\w+.-]{2,}:)?\/\/([^/@]+@)?(.*)`)
	hostPathRe       = regexp.MustCompile(`([^#/?]*)(.*)`)
	pathSplitRe      = regexp.MustCompile(`([^#?]*)(\?[^#]*)?(#.*)?`)
	hostPortRe       = regexp.MustCompile(`([^/:]*):?(\d+)?`)
	driveLetterRe    = regexp.MustCompile(`/([A-Za-z]:)`)
	filenameRe       = regexp.MustCompile(`/([^/]+)$`)
	filenameStrictRe = regexp.MustCompile(`/([^/]+\.[^/]+)$`)
)

// ParsePath splits a path string into pathname, search, and hash.
func ParsePath(input string) ParsedPath {
	if input == "" {
		return ParsedPath{}
	}
	m := pathSplitRe.FindStringSubmatch(input)
	if m == nil {
		return ParsedPath{Pathname: input}
	}
	return ParsedPath{
		Pathname: m[1],
		Search:   m[2],
		Hash:     m[3],
	}
}

// ParseURL parses a URL string into a ParsedURL.
// If the input has no protocol and defaultProto is non-empty, it is prepended.
func ParseURL(input string, defaultProto string) ParsedURL {
	if input == "" {
		return ParsedURL{}
	}
	// Handle special protocols (blob:, data:, javascript:, vbscript:)
	if m := specialProtoRe.FindStringSubmatch(input); m != nil {
		proto := strings.ToLower(m[1])
		pathname := m[2]
		return ParsedURL{
			Protocol: proto,
			Pathname: pathname,
		}
	}
	// Check if input has a protocol or is protocol-relative
	hasProto := hasProtocolRe.MatchString(input) || protocolRelRe.MatchString(input)
	if !hasProto {
		if defaultProto != "" {
			return ParseURL(defaultProto+input, "")
		}
		pp := ParsePath(input)
		return ParsedURL{
			Pathname: pp.Pathname,
			Search:   pp.Search,
			Hash:     pp.Hash,
		}
	}
	// Normalise backslashes
	normalized := strings.ReplaceAll(input, "\\", "/")
	m := fullURLRe.FindStringSubmatch(normalized)
	protocol := ""
	auth := ""
	hostAndPath := ""
	if m != nil {
		protocol = m[1]
		auth = m[2]
		hostAndPath = m[3]
	}
	hm := hostPathRe.FindStringSubmatch(hostAndPath)
	host := ""
	path := ""
	if hm != nil {
		host = hm[1]
		path = hm[2]
	}
	// file: protocol strips leading slash before drive letter
	if strings.ToLower(protocol) == "file:" {
		path = driveLetterRe.ReplaceAllString(path, "$1")
	}
	pp := ParsePath(path)
	// strip trailing @ from auth
	if len(auth) > 0 && auth[len(auth)-1] == '@' {
		auth = auth[:len(auth)-1]
	}
	return ParsedURL{
		Protocol:         strings.ToLower(protocol),
		Auth:             auth,
		Host:             host,
		Pathname:         pp.Pathname,
		Search:           pp.Search,
		Hash:             pp.Hash,
		ProtocolRelative: protocol == "",
	}
}

// ParseAuth splits "username:password" into its components (URL-decoded).
func ParseAuth(input string) ParsedAuth {
	if input == "" {
		return ParsedAuth{}
	}
	idx := strings.Index(input, ":")
	if idx == -1 {
		return ParsedAuth{Username: Decode(input)}
	}
	return ParsedAuth{
		Username: Decode(input[:idx]),
		Password: Decode(input[idx+1:]),
	}
}

// ParseHost splits "hostname:port" into its components.
func ParseHost(input string) ParsedHost {
	if input == "" {
		return ParsedHost{}
	}
	m := hostPortRe.FindStringSubmatch(input)
	if m == nil {
		return ParsedHost{Hostname: Decode(input)}
	}
	return ParsedHost{
		Hostname: Decode(m[1]),
		Port:     m[2],
	}
}

// StringifyParsedURL converts a ParsedURL back to a URL string.
func StringifyParsedURL(parsed ParsedURL) string {
	search := parsed.Search
	if search != "" && !strings.HasPrefix(search, "?") {
		search = "?" + search
	}
	auth := ""
	if parsed.Auth != "" {
		auth = parsed.Auth + "@"
	}
	proto := ""
	if parsed.Protocol != "" || parsed.ProtocolRelative {
		proto = parsed.Protocol + "//"
	}
	return proto + auth + parsed.Host + parsed.Pathname + search + parsed.Hash
}

// ParseFilename returns the last path segment from the URL's pathname.
// With opts.Strict = true, only returns a name that contains a dot (extension).
func ParseFilename(input string, opts ParseFilenameOptions) string {
	parsed := ParseURL(input, "")
	var m []string
	if opts.Strict {
		m = filenameStrictRe.FindStringSubmatch(parsed.Pathname)
	} else {
		m = filenameRe.FindStringSubmatch(parsed.Pathname)
	}
	if m == nil {
		return ""
	}
	return m[1]
}
