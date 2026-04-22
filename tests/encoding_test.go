package ufo_test

import (
	"testing"

	ufo "github.com/givensuman/go-ufo"
)

func TestEncode(t *testing.T) {
	if got := ufo.Encode("hello world"); got != "hello%20world" {
		t.Errorf("Encode space = %q, want %q", got, "hello%20world")
	}
	// pipe character must NOT be encoded per encodeURI semantics
	if got := ufo.Encode("a|b"); got != "a|b" {
		t.Errorf("Encode pipe = %q, want %q", got, "a|b")
	}
}

func TestDecode(t *testing.T) {
	if got := ufo.Decode("hello%20world"); got != "hello world" {
		t.Errorf("Decode = %q, want %q", got, "hello world")
	}
	// Returns original on failure
	if got := ufo.Decode("%zz"); got != "%zz" {
		t.Errorf("Decode bad input = %q, want %q", got, "%zz")
	}
}

func TestEncodeHash(t *testing.T) {
	// { } ^ should be preserved (not encoded)
	got := ufo.EncodeHash("foo{bar}^baz")
	if got != "foo{bar}^baz" {
		t.Errorf("EncodeHash = %q, want %q", got, "foo{bar}^baz")
	}
}

func TestEncodePath(t *testing.T) {
	got := ufo.EncodePath("foo#bar?baz")
	if got != "foo%23bar%3Fbaz" {
		t.Errorf("EncodePath = %q, want %q", got, "foo%23bar%3Fbaz")
	}
}

func TestEncodeParam(t *testing.T) {
	if got := ufo.EncodeParam("foo/bar"); got != "foo%2Fbar" {
		t.Errorf("EncodeParam = %q, want %q", got, "foo%2Fbar")
	}
}

func TestEncodeQueryValue(t *testing.T) {
	// space → +
	if got := ufo.EncodeQueryValue("hello world"); got != "hello+world" {
		t.Errorf("EncodeQueryValue space = %q, want %q", got, "hello+world")
	}
	// literal + → %2B
	if got := ufo.EncodeQueryValue("a+b"); got != "a%2Bb" {
		t.Errorf("EncodeQueryValue plus = %q, want %q", got, "a%2Bb")
	}
}

func TestEncodeQueryKey(t *testing.T) {
	if got := ufo.EncodeQueryKey("foo=bar"); got != "foo%3Dbar" {
		t.Errorf("EncodeQueryKey = %q, want %q", got, "foo%3Dbar")
	}
}

func TestDecodePath(t *testing.T) {
	// %2F in path should NOT become / — it stays as %2F
	if got := ufo.DecodePath("foo%2Fbar"); got != "foo%2Fbar" {
		t.Errorf("DecodePath %%2F = %q, want %q", got, "foo%2Fbar")
	}
	if got := ufo.DecodePath("foo%20bar"); got != "foo bar" {
		t.Errorf("DecodePath space = %q, want %q", got, "foo bar")
	}
}

func TestDecodeQueryKey(t *testing.T) {
	if got := ufo.DecodeQueryKey("hello+world"); got != "hello world" {
		t.Errorf("DecodeQueryKey = %q, want %q", got, "hello world")
	}
}

func TestDecodeQueryValue(t *testing.T) {
	if got := ufo.DecodeQueryValue("hello+world"); got != "hello world" {
		t.Errorf("DecodeQueryValue = %q, want %q", got, "hello world")
	}
}

func TestEncodeHost(t *testing.T) {
	// ASCII hostname unchanged
	if got := ufo.EncodeHost("example.com"); got != "example.com" {
		t.Errorf("EncodeHost ASCII = %q, want %q", got, "example.com")
	}
	// IDN hostname gets punycoded
	if got := ufo.EncodeHost("münchen.de"); got != "xn--mnchen-3ya.de" {
		t.Errorf("EncodeHost IDN = %q, want %q", got, "xn--mnchen-3ya.de")
	}
}
