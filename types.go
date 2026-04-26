package ufo

// ParsedURL is the result of ParseURL.
type ParsedURL struct {
	Protocol         string // e.g. "http:"
	Auth             string // e.g. "user:pass"
	Host             string // e.g. "foo.com" or "foo.com:8080"
	Pathname         string // e.g. "/foo"
	Search           string // e.g. "?test=123"
	Hash             string // e.g. "#token"
	ProtocolRelative bool   // true when input had no protocol (e.g. "//foo.com")
}

// ParsedPath holds the path components split from a URL string.
type ParsedPath struct {
	Pathname string
	Search   string
	Hash     string
}

// ParsedAuth holds decoded username and password from a URL auth segment.
type ParsedAuth struct {
	Username string
	Password string
}

// ParsedHost holds hostname and port parsed from a host string.
type ParsedHost struct {
	Hostname string
	Port     string
}

// CompareURLOptions controls the behaviour of IsEqual.
type CompareURLOptions struct {
	TrailingSlash bool
	LeadingSlash  bool
	Encoding      bool
}

// ParseFilenameOptions controls the behaviour of ParseFilename.
type ParseFilenameOptions struct {
	// Strict requires the final segment to contain a dot (extension).
	Strict bool
}

// QueryValue is any scalar or composite value that can appear in a query string.
type QueryValue any

// QueryObject is a map of query key to one or more QueryValues.
type QueryObject map[string]any
