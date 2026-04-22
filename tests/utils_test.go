package ufo_test

import (
	"testing"

	ufo "github.com/givensuman/go-ufo"
)

func TestHasProtocol(t *testing.T) {
	if !ufo.HasProtocol("https://example.com", ufo.HasProtocolOptions{}) {
		t.Error("HasProtocol https should be true")
	}
	if ufo.HasProtocol("//example.com", ufo.HasProtocolOptions{}) {
		t.Error("HasProtocol // without acceptRelative should be false")
	}
	if !ufo.HasProtocol("//example.com", ufo.HasProtocolOptions{AcceptRelative: true}) {
		t.Error("HasProtocol // with acceptRelative should be true")
	}
	if !ufo.HasProtocol("ftp://example.com", ufo.HasProtocolOptions{}) {
		t.Error("HasProtocol ftp should be true")
	}
	if ufo.HasProtocol("data:text/plain", ufo.HasProtocolOptions{Strict: true}) {
		t.Error("HasProtocol data strict should be false")
	}
	if !ufo.HasProtocol("data:text/plain", ufo.HasProtocolOptions{}) {
		t.Error("HasProtocol data non-strict should be true")
	}
}

func TestIsScriptProtocol(t *testing.T) {
	for _, proto := range []string{"javascript:", "data:", "blob:", "vbscript:"} {
		if !ufo.IsScriptProtocol(proto) {
			t.Errorf("IsScriptProtocol %q should be true", proto)
		}
	}
	if ufo.IsScriptProtocol("https:") {
		t.Error("IsScriptProtocol https should be false")
	}
}

func TestHasTrailingSlash(t *testing.T) {
	if !ufo.HasTrailingSlash("/foo/", false) {
		t.Error("/foo/ should have trailing slash")
	}
	if ufo.HasTrailingSlash("/foo", false) {
		t.Error("/foo should not have trailing slash")
	}
	if ufo.HasTrailingSlash("/foo?q=1", true) {
		t.Error("/foo?q=1 respectQuery — no trailing slash")
	}
	if !ufo.HasTrailingSlash("/foo/?q=1", true) {
		t.Error("/foo/?q=1 respectQuery — has trailing slash")
	}
}

func TestWithoutTrailingSlash(t *testing.T) {
	if got := ufo.WithoutTrailingSlash("/foo/", false); got != "/foo" {
		t.Errorf("WithoutTrailingSlash = %q, want %q", got, "/foo")
	}
	if got := ufo.WithoutTrailingSlash("/path/?query=true", true); got != "/path?query=true" {
		t.Errorf("WithoutTrailingSlash respect = %q, want %q", got, "/path?query=true")
	}
	// no trailing slash → unchanged
	if got := ufo.WithoutTrailingSlash("/foo", false); got != "/foo" {
		t.Errorf("WithoutTrailingSlash no-op = %q, want %q", got, "/foo")
	}
}

func TestWithTrailingSlash(t *testing.T) {
	if got := ufo.WithTrailingSlash("/foo", false); got != "/foo/" {
		t.Errorf("WithTrailingSlash = %q, want %q", got, "/foo/")
	}
	if got := ufo.WithTrailingSlash("/path?q=1", true); got != "/path/?q=1" {
		t.Errorf("WithTrailingSlash respect = %q, want %q", got, "/path/?q=1")
	}
	// already has trailing slash → unchanged
	if got := ufo.WithTrailingSlash("/foo/", false); got != "/foo/" {
		t.Errorf("WithTrailingSlash already = %q, want %q", got, "/foo/")
	}
}

func TestHasLeadingSlash(t *testing.T) {
	if !ufo.HasLeadingSlash("/foo") {
		t.Error("/foo should have leading slash")
	}
	if ufo.HasLeadingSlash("foo") {
		t.Error("foo should not have leading slash")
	}
}

func TestWithLeadingSlash(t *testing.T) {
	if got := ufo.WithLeadingSlash("foo"); got != "/foo" {
		t.Errorf("WithLeadingSlash = %q, want %q", got, "/foo")
	}
	if got := ufo.WithLeadingSlash("/foo"); got != "/foo" {
		t.Errorf("WithLeadingSlash already = %q, want %q", got, "/foo")
	}
}

func TestWithoutLeadingSlash(t *testing.T) {
	if got := ufo.WithoutLeadingSlash("/foo"); got != "foo" {
		t.Errorf("WithoutLeadingSlash = %q, want %q", got, "foo")
	}
}

func TestCleanDoubleSlashes(t *testing.T) {
	if got := ufo.CleanDoubleSlashes("//foo//bar//"); got != "/foo/bar/" {
		t.Errorf("CleanDoubleSlashes = %q, want %q", got, "/foo/bar/")
	}
	got := ufo.CleanDoubleSlashes("http://example.com/analyze//http://localhost:3000//")
	want := "http://example.com/analyze/http://localhost:3000/"
	if got != want {
		t.Errorf("CleanDoubleSlashes protocol = %q, want %q", got, want)
	}
}

func TestWithBase(t *testing.T) {
	// input already starts with base → unchanged
	if got := ufo.WithBase("/foo/bar", "/foo"); got != "/foo/bar" {
		t.Errorf("WithBase already has base = %q, want %q", got, "/foo/bar")
	}
	// input does not start with base → prepend base
	if got := ufo.WithBase("/foo/bar", "/baz"); got != "/baz/foo/bar" {
		t.Errorf("WithBase adds base = %q, want %q", got, "/baz/foo/bar")
	}
	// no false positives: /admin-dashboard should not match base /admin
	if got := ufo.WithBase("/admin-dashboard", "/admin"); got != "/admin/admin-dashboard" {
		t.Errorf("WithBase partial = %q, want %q", got, "/admin/admin-dashboard")
	}
}

func TestWithoutBase(t *testing.T) {
	if got := ufo.WithoutBase("/foo/bar", "/foo"); got != "/bar" {
		t.Errorf("WithoutBase = %q, want %q", got, "/bar")
	}
	// no base → unchanged
	if got := ufo.WithoutBase("/foo/bar", ""); got != "/foo/bar" {
		t.Errorf("WithoutBase empty base = %q", got)
	}
}

func TestJoinURL(t *testing.T) {
	if got := ufo.JoinURL("a", "/b", "/c"); got != "a/b/c" {
		t.Errorf("JoinURL = %q, want %q", got, "a/b/c")
	}
	if got := ufo.JoinURL("", "/foo"); got != "/foo" {
		t.Errorf("JoinURL empty base = %q, want %q", got, "/foo")
	}
}

func TestJoinRelativeURL(t *testing.T) {
	if got := ufo.JoinRelativeURL("/a", "../b", "./c"); got != "/b/c" {
		t.Errorf("JoinRelativeURL = %q, want %q", got, "/b/c")
	}
}

func TestWithProtocol(t *testing.T) {
	if got := ufo.WithProtocol("http://example.com", "ftp://"); got != "ftp://example.com" {
		t.Errorf("WithProtocol = %q, want %q", got, "ftp://example.com")
	}
}

func TestWithHTTP(t *testing.T) {
	if got := ufo.WithHTTP("https://example.com"); got != "http://example.com" {
		t.Errorf("WithHTTP = %q, want %q", got, "http://example.com")
	}
}

func TestWithHTTPS(t *testing.T) {
	if got := ufo.WithHTTPS("http://example.com"); got != "https://example.com" {
		t.Errorf("WithHTTPS = %q, want %q", got, "https://example.com")
	}
}

func TestWithoutProtocol(t *testing.T) {
	if got := ufo.WithoutProtocol("http://example.com"); got != "example.com" {
		t.Errorf("WithoutProtocol = %q, want %q", got, "example.com")
	}
}

func TestWithQuery(t *testing.T) {
	got := ufo.WithQuery("/foo?page=a", ufo.QueryObject{"token": "secret"})
	// order not guaranteed, check both params present
	if got != "/foo?page=a&token=secret" && got != "/foo?token=secret&page=a" {
		t.Errorf("WithQuery = %q", got)
	}
}

func TestGetQuery(t *testing.T) {
	q := ufo.GetQuery("http://foo.com/foo?test=123&unicode=%E5%A5%BD")
	if q["test"] != "123" {
		t.Errorf("GetQuery test = %v", q["test"])
	}
	if q["unicode"] != "好" {
		t.Errorf("GetQuery unicode = %v", q["unicode"])
	}
}

func TestFilterQuery(t *testing.T) {
	got := ufo.FilterQuery("/foo?bar=1&baz=2", func(key string, value interface{}) bool {
		return key != "bar"
	})
	if got != "/foo?baz=2" {
		t.Errorf("FilterQuery = %q, want %q", got, "/foo?baz=2")
	}
	// no query → unchanged
	if got2 := ufo.FilterQuery("/foo", func(k string, v interface{}) bool { return true }); got2 != "/foo" {
		t.Errorf("FilterQuery no query = %q", got2)
	}
}

func TestNormalizeURL(t *testing.T) {
	got := ufo.NormalizeURL("test?query=123 123#hash, test")
	want := "test?query=123%20123#hash,%20test"
	if got != want {
		t.Errorf("NormalizeURL = %q, want %q", got, want)
	}
}

func TestResolveURL(t *testing.T) {
	got := ufo.ResolveURL("http://foo.com/foo?test=123#token", "bar", "baz")
	want := "http://foo.com/foo/bar/baz?test=123#token"
	if got != want {
		t.Errorf("ResolveURL = %q, want %q", got, want)
	}
}

func TestIsSamePath(t *testing.T) {
	if !ufo.IsSamePath("/foo", "/foo/") {
		t.Error("IsSamePath /foo /foo/ should be equal")
	}
	if ufo.IsSamePath("/foo", "/bar") {
		t.Error("IsSamePath /foo /bar should not be equal")
	}
}

func TestIsEqual(t *testing.T) {
	if !ufo.IsEqual("/foo", "foo", ufo.CompareURLOptions{}) {
		t.Error("IsEqual /foo foo default")
	}
	if !ufo.IsEqual("/foo bar", "/foo%20bar", ufo.CompareURLOptions{}) {
		t.Error("IsEqual encoding default")
	}
	if ufo.IsEqual("/foo", "foo", ufo.CompareURLOptions{LeadingSlash: true}) {
		t.Error("IsEqual strict leading should differ")
	}
}

func TestWithFragment(t *testing.T) {
	if got := ufo.WithFragment("/foo", "bar"); got != "/foo#bar" {
		t.Errorf("WithFragment = %q, want %q", got, "/foo#bar")
	}
	if got := ufo.WithFragment("/foo#bar", "baz"); got != "/foo#baz" {
		t.Errorf("WithFragment replace = %q, want %q", got, "/foo#baz")
	}
	if got := ufo.WithFragment("/foo#bar", ""); got != "/foo" {
		t.Errorf("WithFragment empty = %q, want %q", got, "/foo")
	}
}

func TestWithoutFragment(t *testing.T) {
	got := ufo.WithoutFragment("http://example.com/foo?q=123#bar")
	if got != "http://example.com/foo?q=123" {
		t.Errorf("WithoutFragment = %q, want %q", got, "http://example.com/foo?q=123")
	}
}

func TestWithoutHost(t *testing.T) {
	got := ufo.WithoutHost("http://example.com/foo?q=123#bar")
	if got != "/foo?q=123#bar" {
		t.Errorf("WithoutHost = %q, want %q", got, "/foo?q=123#bar")
	}
}

func TestIsRelative(t *testing.T) {
	if !ufo.IsRelative("./foo") {
		t.Error("./foo should be relative")
	}
	if !ufo.IsRelative("../foo") {
		t.Error("../foo should be relative")
	}
	if ufo.IsRelative("/foo") {
		t.Error("/foo should not be relative")
	}
}

func TestIsEmptyURL(t *testing.T) {
	if !ufo.IsEmptyURL("") {
		t.Error("empty string should be empty URL")
	}
	if !ufo.IsEmptyURL("/") {
		t.Error("/ should be empty URL")
	}
	if ufo.IsEmptyURL("/foo") {
		t.Error("/foo should not be empty URL")
	}
}
