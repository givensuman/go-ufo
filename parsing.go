package ufo

import (
	"fmt"
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

	if m := specialProtoRe.FindStringSubmatch(input); m != nil {
		proto := strings.ToLower(m[1])
		pathname := m[2]
		return ParsedURL{
			Protocol: proto,
			Pathname: pathname,
		}
	}

	hasProto := hasProtocolRe.MatchString(input) || protocolRelRe.MatchString(input)
	if !hasProto {
		if defaultProto != "" {
			return ParseURL(defaultProto+input, "")
		}

		parse := ParsePath(input)
		return ParsedURL{
			Pathname: parse.Pathname,
			Search:   parse.Search,
			Hash:     parse.Hash,
		}
	}

	// Normalise backslashes
	normalized := strings.ReplaceAll(input, "\\", "/")
	m := fullURLRe.FindStringSubmatch(normalized)

	var protocol, auth, hostAndPath string

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

	// 'file:' protocol strips leading slash before drive letter
	if strings.ToLower(protocol) == "file:" {
		path = driveLetterRe.ReplaceAllString(path, "$1")
	}

	parse := ParsePath(path)
	// strip trailing @ from auth
	if len(auth) > 0 && auth[len(auth)-1] == '@' {
		auth = auth[:len(auth)-1]
	}

	return ParsedURL{
		Protocol:         strings.ToLower(protocol),
		Auth:             auth,
		Host:             host,
		Pathname:         parse.Pathname,
		Search:           parse.Search,
		Hash:             parse.Hash,
		ProtocolRelative: protocol == "",
	}
}

// ParseAuth splits "username:password" into its components.
func ParseAuth(input string) ParsedAuth {
	if input == "" {
		return ParsedAuth{}
	}

	username, password, found := strings.Cut(input, ":")
	if !found {
		return ParsedAuth{Username: Decode(input)}
	}

	return ParsedAuth{
		Username: Decode(username),
		Password: Decode(password),
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

// String converts a ParsedURL back to a URL string.
func (url ParsedURL) String() string {
	if url.Search != "" && !strings.HasPrefix(url.Search, "?") {
		url.Search = "?" + url.Search
	}

	if url.Auth != "" {
		url.Auth = url.Auth + "@"
	}

	if url.Protocol != "" || url.ProtocolRelative {
		url.Protocol = url.Protocol + "//"
	}

	return fmt.Sprintf("%s%s%s%s%s%s",
		url.Protocol,
		url.Auth,
		url.Host,
		url.Pathname,
		url.Search,
		url.Hash,
	)
}

var _ fmt.Stringer = (*ParsedURL)(nil)

// ParseFilename returns the last path segment from the URL's pathname.
// With `opts.Strict` = true, only returns a name that contains a dot.
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
