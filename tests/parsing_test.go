package ufo_test

import (
	"testing"

	ufo "github.com/givensuman/go-ufo"
)

func TestParsePath(t *testing.T) {
	p := ufo.ParsePath("/foo?test=123#token")
	if p.Pathname != "/foo" || p.Search != "?test=123" || p.Hash != "#token" {
		t.Errorf("ParsePath = %+v", p)
	}
	p2 := ufo.ParsePath("")
	if p2.Pathname != "" || p2.Search != "" || p2.Hash != "" {
		t.Errorf("ParsePath empty = %+v", p2)
	}
}

func TestParseURL(t *testing.T) {
	// Full URL
	p := ufo.ParseURL("http://foo.com/foo?test=123#token", "")
	if p.Protocol != "http:" || p.Host != "foo.com" || p.Pathname != "/foo" ||
		p.Search != "?test=123" || p.Hash != "#token" {
		t.Errorf("ParseURL full = %+v", p)
	}
	// No protocol — returns path-only result
	p2 := ufo.ParseURL("foo.com/foo?test=123#token", "")
	if p2.Pathname != "foo.com/foo" || p2.Search != "?test=123" || p2.Hash != "#token" {
		t.Errorf("ParseURL no proto = %+v", p2)
	}
	// Default protocol supplied
	p3 := ufo.ParseURL("foo.com/foo?test=123#token", "https://")
	if p3.Protocol != "https:" || p3.Host != "foo.com" {
		t.Errorf("ParseURL default proto = %+v", p3)
	}
	// Protocol-relative
	p4 := ufo.ParseURL("//foo.com/path", "")
	if p4.Host != "foo.com" || !p4.ProtocolRelative {
		t.Errorf("ParseURL protocol-relative = %+v", p4)
	}
	// Empty input
	p5 := ufo.ParseURL("", "")
	if p5.Pathname != "" {
		t.Errorf("ParseURL empty = %+v", p5)
	}
}

func TestParseAuth(t *testing.T) {
	a := ufo.ParseAuth("user:pass")
	if a.Username != "user" || a.Password != "pass" {
		t.Errorf("ParseAuth = %+v", a)
	}
	a2 := ufo.ParseAuth("user")
	if a2.Username != "user" || a2.Password != "" {
		t.Errorf("ParseAuth no pass = %+v", a2)
	}
	a3 := ufo.ParseAuth("")
	if a3.Username != "" || a3.Password != "" {
		t.Errorf("ParseAuth empty = %+v", a3)
	}
}

func TestParseHost(t *testing.T) {
	h := ufo.ParseHost("foo.com:8080")
	if h.Hostname != "foo.com" || h.Port != "8080" {
		t.Errorf("ParseHost with port = %+v", h)
	}
	h2 := ufo.ParseHost("foo.com")
	if h2.Hostname != "foo.com" || h2.Port != "" {
		t.Errorf("ParseHost no port = %+v", h2)
	}
}

func TestStringifyParsedURL(t *testing.T) {
	p := ufo.ParseURL("http://foo.com/foo?test=123#token", "")
	p.Host = "bar.com"
	got := p.String()
	if got != "http://bar.com/foo?test=123#token" {
		t.Errorf("StringifyParsedURL = %q, want %q", got, "http://bar.com/foo?test=123#token")
	}
	// Protocol-relative
	p2 := ufo.ParseURL("//foo.com/path", "")
	got2 := p2.String()
	if got2 != "//foo.com/path" {
		t.Errorf("StringifyParsedURL proto-relative = %q, want %q", got2, "//foo.com/path")
	}
}

func TestParseFilename(t *testing.T) {
	got := ufo.ParseFilename("http://example.com/path/to/filename.ext", ufo.ParseFilenameOptions{})
	if got != "filename.ext" {
		t.Errorf("ParseFilename = %q, want %q", got, "filename.ext")
	}
	// strict mode: file with extension
	got2 := ufo.ParseFilename("http://example.com/path/to/filename.ext", ufo.ParseFilenameOptions{Strict: true})
	if got2 != "filename.ext" {
		t.Errorf("ParseFilename strict = %q, want %q", got2, "filename.ext")
	}
	// strict mode: hidden file (no extension after dot at start)
	got3 := ufo.ParseFilename("/path/to/.hidden-file", ufo.ParseFilenameOptions{Strict: true})
	if got3 != "" {
		t.Errorf("ParseFilename strict hidden = %q, want %q", got3, "")
	}
	// non-strict: returns last segment even without extension
	got4 := ufo.ParseFilename("/path/to/noext", ufo.ParseFilenameOptions{})
	if got4 != "noext" {
		t.Errorf("ParseFilename noext = %q, want %q", got4, "noext")
	}
}
