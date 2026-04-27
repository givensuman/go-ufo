<div align="center">
  <img src="./assets/logo.svg" alt="go-ufo" width="200" />
</div>

# go-ufo

This is a port of the [unjs.io/ufo](https://unjs.io/packages/ufo) package to Go.

## Install

```bash
go get github.com/givensuman/go-ufo
```

Import:

```go
import (
  "github.com/givensuman/go-ufo"
)
```

## Encoding Utils

### `Encode(text string) string`

Encodes characters that need to be encoded in the path, search and hash sections of the URL.

### `Decode(text string) string`

Decodes text, returning the original text if decoding fails.

### `DecodePath(text string) string`

Decodes path section of the URL, consistent with `EncodePath`.

### `DecodeQueryKey(text string) string`

Decodes query key, consistent with `EncodeQueryKey`.

### `DecodeQueryValue(text string) string`

Decodes query value, consistent with `EncodeQueryValue`.

### `EncodeHash(text string) string`

Encodes characters that need to be encoded in the hash section of the URL.

### `EncodePath(text string) string`

Encodes characters that need to be encoded in the path section of the URL.

### `EncodeParam(text string) string`

Encodes characters that need to be encoded in the path section of the URL as a param. Also encodes the slash (`/`) character.

### `EncodeQueryValue(input string) string`

Encodes characters that need to be encoded for query values in the query section of the URL.

### `EncodeQueryKey(text string) string`

Encodes characters that need to be encoded for query keys in the query section of the URL and also encodes the `=` character.

### `EncodeHost(name string) string`

Encodes hostname with punycode encoding, returning the original name if encoding fails.

## Parsing Utils

### `ParsePath(input string) ParsedPath`

Splits an input string into pathname, search, and hash:

```go
ParsePath("http://foo.com/foo?test=123#token")
// ParsedPath{ Pathname: "http://foo.com/foo", Search: "?test=123", Hash: "#token" }
```

### `ParseURL(input string, defaultProto string) ParsedURL`

Parses a URL string into a `ParsedURL`. If the input has no protocol and `defaultProto` is non-empty, it is prepended:

```go
ParseURL("http://foo.com/foo?test=123#token", "")
// ParsedURL{ Protocol: "http:", Host: "foo.com", Pathname: "/foo", Search: "?test=123", Hash: "#token" }

ParseURL("foo.com/foo?test=123#token", "")
// ParsedURL{ Pathname: "foo.com/foo", Search: "?test=123", Hash: "#token" }

ParseURL("foo.com/foo?test=123#token", "https://")
// ParsedURL{ Protocol: "https:", Host: "foo.com", Pathname: "/foo", Search: "?test=123", Hash: "#token" }
```

`ParsedURL` implements `fmt.Stringer`, so you can call `.String()` to reconstruct the URL.

### `ParseAuth(input string) ParsedAuth`

Splits a `username:password` string into its components:

```go
ParseAuth("user:pass")
// ParsedAuth{ Username: "user", Password: "pass" }
```

### `ParseHost(input string) ParsedHost`

Splits a `hostname:port` string into its components:

```go
ParseHost("foo.com:8080")
// ParsedHost{ Hostname: "foo.com", Port: "8080" }
```

### `ParseFilename(input string, opts ParseFilenameOptions) string`

Returns the last path segment from the URL's pathname. With `opts.Strict = true`, only returns a name that contains a dot (extension):

```go
ParseFilename("http://example.com/path/to/filename.ext", ParseFilenameOptions{})
// "filename.ext"

ParseFilename("/path/to/.hidden-file", ParseFilenameOptions{ Strict: true })
// ""
```

## Query Utils

### `ParseQuery(parametersString string) map[string]any`

Parses and decodes a query string into a map. The input can include or omit the leading `?`:

```go
ParseQuery("?foo=bar&baz=qux")
// map[string]any{ "foo": "bar", "baz": "qux" }

ParseQuery("tags=javascript&tags=web&tags=dev")
// map[string]any{ "tags": []string{"javascript", "web", "dev"} }
```

### `EncodeQueryItem(key string, value any) string`

Encodes a key-value pair into a URL query string. If the value is a `[]any`, it is encoded as multiple key-value pairs with the same key:

```go
EncodeQueryItem("message", "Hello World")
// "message=Hello+World"

EncodeQueryItem("tags", []any{"javascript", "web", "dev"})
// "tags=javascript&tags=web&tags=dev"
```

## Utils

### `CleanDoubleSlashes(input string) string`

Removes double slashes from the URL:

```go
CleanDoubleSlashes("//foo//bar//") // "/foo/bar/"

CleanDoubleSlashes("http://example.com/analyze//http://localhost:3000//")
// "http://example.com/analyze/http://localhost:3000/"
```

### `FilterQuery(input string, callback func(k string, v any) bool) string`

Filters query params of the URL using a predicate:

```go
FilterQuery("/foo?bar=1&baz=2", func(k string, v any) bool {
    return k != "bar"
}) // "/foo?baz=2"
```

### `GetQuery(input string) map[string]any`

Parses and decodes the query object of an input URL into a map:

```go
GetQuery("http://foo.com/foo?test=123&unicode=%E5%A5%BD")
// map[string]any{ "test": "123", "unicode": "好" }
```

### `HasLeadingSlash(input string) bool`

Checks if the input has a leading slash:

```go
HasLeadingSlash("/foo") // true
```

### `HasProtocol(inputString string, opts *HasProtocolOptions) bool`

Checks if the input has a protocol. Use `opts.AcceptRelative = true` to accept relative URLs, and `opts.Strict = true` for strict RFC 3986 protocol matching:

```go
HasProtocol("https://example.com", nil) // true
HasProtocol("//example.com", nil) // false
HasProtocol("//example.com", &HasProtocolOptions{ AcceptRelative: true }) // true
HasProtocol("data:text/plain", nil) // true
HasProtocol("data:text/plain", &HasProtocolOptions{ Strict: true }) // false
```

### `HasTrailingSlash(input string, respectQueryAndFragment *RespectQueryAndFragmentOption) bool`

Checks if the input has a trailing slash:

```go
HasTrailingSlash("/foo/", nil) // true
HasTrailingSlash("/foo", nil) // false
HasTrailingSlash("/foo?query=true", &RespectQueryAndFragmentOption{ true }) // false
HasTrailingSlash("/foo/?query=true", &RespectQueryAndFragmentOption{ true }) // true
```

### `IsEmptyURL(url string) bool`

Checks if the input URL is empty or `/`.

### `IsEqual(a string, b string, opts *IsEqualOptions) bool`

Checks if two paths are equal regardless of encoding, trailing slash, and leading slash differences. Use `opts` for strict comparisons:

```go
IsEqual("/foo", "foo", nil) // true
IsEqual("foo/", "foo", nil) // true
IsEqual("/foo bar", "/foo%20bar", nil) // true

// Strict comparisons
IsEqual("/foo", "foo", &IsEqualOptions{ StrictLeadingSlash: true }) // false
IsEqual("foo/", "foo", &IsEqualOptions{ StrictTrailingSlash: true }) // false
IsEqual("/foo bar", "/foo%20bar", &IsEqualOptions{ StrictEncoding: true }) // false
```

### `IsNonEmptyURL(url string) bool`

Checks if the input URL is neither empty nor `/`.

### `IsRelative(inputString string) bool`

Checks if a path starts with `./` or `../`:

```go
IsRelative("./foo") // true
```

### `IsSamePath(a string, b string) bool`

Checks if two paths are equal. Trailing slash and encoding are normalized before comparison:

```go
IsSamePath("/foo", "/foo/") // true
```

### `IsScriptProtocol(protocol string) bool`

Checks if the input is any of the dangerous `blob:`, `data:`, `javascript:`, or `vbscript:` protocols:

```go
IsScriptProtocol("javascript:alert(1)") // true
IsScriptProtocol("https://example.com") // false
```

### `JoinRelativeURL(input ...string) string`

Joins multiple URL segments into a single URL, handling relative paths with `./` and `../`:

```go
JoinRelativeURL("/a", "../b", "./c") // "/b/c"
```

### `JoinURL(base string, input ...string) string`

Joins multiple URL segments into a single URL:

```go
JoinURL("a", "/b", "/c") // "a/b/c"
```

### `NormalizeURL(input string) string`

Normalizes the input URL by ensuring proper encoding and that the pathname starts with a slash:

```go
NormalizeURL("test?query=123 123#hash, test")
// "test?query=123%20123#hash,%20test"
```

### `ResolveURL(base string, inputs ...string) string`

Resolves multiple URL segments into a single URL:

```go
ResolveURL("http://foo.com/foo?test=123#token", "bar", "baz")
// "http://foo.com/foo/bar/baz?test=123#token"
```

### `WithBase(input string, base string) string`

Ensures the URL or pathname starts with base:

```go
WithBase("/foo/bar", "/foo") // "/foo/bar"
WithBase("/foo/bar", "/baz") // "/baz/foo/bar"
```

### `WithFragment(input string, hash string) string`

Adds or replaces the fragment section of the URL:

```go
WithFragment("/foo", "bar") // "/foo#bar"
WithFragment("/foo#bar", "baz") // "/foo#baz"
WithFragment("/foo#bar", "") // "/foo"
```

### `WithHTTP(input string) string`

Adds or replaces the URL protocol with `http://`:

```go
WithHTTP("https://example.com") // "http://example.com"
```

### `WithHTTPS(input string) string`

Adds or replaces the URL protocol with `https://`:

```go
WithHTTPS("http://example.com") // "https://example.com"
```

### `WithLeadingSlash(input string) string`

Ensures the URL or pathname has a leading slash:

```go
WithLeadingSlash("foo") // "/foo"
```

### `WithoutBase(input string, base string) string`

Removes the base from the URL or pathname:

```go
WithoutBase("/foo/bar", "/foo") // "/bar"
```

### `WithoutFragment(input string) string`

Removes the fragment section from the URL:

```go
WithoutFragment("http://example.com/foo?q=123#bar")
// "http://example.com/foo?q=123"
```

### `WithoutHost(input string) string`

Removes the host from the URL while preserving everything else:

```go
WithoutHost("http://example.com/foo?q=123#bar")
// "/foo?q=123#bar"
```

### `WithoutLeadingSlash(input string) string`

Removes the leading slash from the URL or pathname:

```go
WithoutLeadingSlash("/foo") // "foo"
```

### `WithoutProtocol(input string) string`

Removes the protocol from the input:

```go
WithoutProtocol("http://example.com") // "example.com"
```

### `WithoutTrailingSlash(input string, respectQueryAndFragment *RespectQueryAndFragmentOption) string`

Removes the trailing slash from the URL or pathname:

```go
WithoutTrailingSlash("/foo/", nil) // "/foo"
WithoutTrailingSlash("/path/?query=true", &RespectQueryAndFragmentOption{ true }) // "/path?query=true"
```

### `WithProtocol(input string, protocol string) string`

Adds or replaces the protocol of the input URL:

```go
WithProtocol("http://example.com", "ftp://") // "ftp://example.com"
```

### `WithQuery(input string, query QueryObject) string`

Adds or replaces query params of the URL:

```go
WithQuery("/foo?page=a", QueryObject{ "token": "secret" })
// "/foo?page=a&token=secret"
```

### `WithTrailingSlash(input string, respectQueryAndFragment *RespectQueryAndFragmentOption) string`

Ensures the URL ends with a trailing slash:

```go
WithTrailingSlash("/foo", nil) // "/foo/"
WithTrailingSlash("/path?query=true", &RespectQueryAndFragmentOption{ true }) // "/path/?query=true"
```

## Why?

I'm a fan of the [unjs.io](https://unjs.io) ecosystem, particularly their excellent API design. So I thought why not Go?

## License

[MIT](./LICENSE)
