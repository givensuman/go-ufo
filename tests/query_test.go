package ufo_test

import (
	"strings"
	"testing"

	ufo "github.com/givensuman/go-ufo"
)

func TestParseQuery(t *testing.T) {
	q := ufo.ParseQuery("?foo=bar&baz=qux")
	if len(q["foo"]) == 0 || q["foo"][0] != "bar" || len(q["baz"]) == 0 || q["baz"][0] != "qux" {
		t.Errorf("ParseQuery basic = %v", q)
	}
	// repeated key collects into a slice
	q2 := ufo.ParseQuery("tags=js&tags=web&tags=dev")
	tags := q2["tags"]
	if len(tags) != 3 || tags[0] != "js" {
		t.Errorf("ParseQuery repeated key = %v", q2["tags"])
	}
	// __proto__ is ignored
	q3 := ufo.ParseQuery("__proto__=bad")
	if _, exists := q3["__proto__"]; exists {
		t.Error("ParseQuery must ignore __proto__")
	}
	// empty string
	q4 := ufo.ParseQuery("")
	if len(q4) != 0 {
		t.Errorf("ParseQuery empty = %v", q4)
	}
}

func TestStringifyQuery(t *testing.T) {
	q := ufo.QueryObject{"foo": {"bar"}, "baz": {"qux"}}
	got := q.String()
	// Order is not guaranteed for maps, so check both parts present
	if !strings.Contains(got, "foo=bar") || !strings.Contains(got, "baz=qux") {
		t.Errorf("StringifyQuery = %q", got)
	}
	// nil/empty values are omitted
	q2 := ufo.QueryObject{"a": {"1"}, "b": nil}
	got2 := q2.String()
	if got2 != "a=1" {
		t.Errorf("StringifyQuery nil value = %q, want %q", got2, "a=1")
	}
}

func TestEncodeQueryItem(t *testing.T) {
	if got := ufo.EncodeQueryItem("message", []string{"Hello World"}); got != "message=Hello+World" {
		t.Errorf("EncodeQueryItem = %q, want %q", got, "message=Hello+World")
	}
	// multiple values
	got := ufo.EncodeQueryItem("tags", []string{"js", "web", "dev"})
	if got != "tags=js&tags=web&tags=dev" {
		t.Errorf("EncodeQueryItem array = %q, want %q", got, "tags=js&tags=web&tags=dev")
	}
	// nil/empty value slice — just key
	if got := ufo.EncodeQueryItem("key", nil); got != "key" {
		t.Errorf("EncodeQueryItem nil = %q, want %q", got, "key")
	}
}
