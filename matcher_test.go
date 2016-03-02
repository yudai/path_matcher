package pmatcher

import (
	"testing"
)

func TestMatcher(t *testing.T) {
	var (
		matcher *Matcher
		match   bool
		pattern string
		params  map[string]string
	)

	matcher = New()
	matcher.Add("/foo/bar")
	match, pattern, params = matcher.Match("/foo/bar")
	if !match {
		t.Fail()
	}
	if pattern != "/foo/bar" {
		t.Fail()
	}
	if len(params) != 0 {
		t.Fail()
	}
	match, pattern, params = matcher.Match("/foo/baz")
	if match {
		t.Fail()
	}
	match, pattern, params = matcher.Match("/foo/bar/baz")
	if match {
		t.Fail()
	}

	matcher = New()
	matcher.Add("/foo/:p1/:p2")
	match, pattern, params = matcher.Match("/foo/bar/baz")
	if !match {
		t.Fail()
	}
	if pattern != "/foo/:p1/:p2" {
		t.Fail()
	}
	if len(params) != 2 || params["p1"] != "bar" || params["p2"] != "baz" {
		t.Fail()
	}
	match, pattern, params = matcher.Match("/foo/bar")
	if match {
		t.Fail()
	}
	match, pattern, params = matcher.Match("/foo")
	if match {
		t.Fail()
	}
	match, pattern, params = matcher.Match("/foo/bar/baz/foo")
	if match {
		t.Fail()
	}

	matcher = New()
	matcher.Add("/foo/bar")
	matcher.Add("/foo/:p1/:p2")
	matcher.Add("/foo/bar/baz")
	matcher.Add("/baz/:p1")
	matcher.Add("/baz/one")
	matcher.Add("/baz/:p1/two")
	matcher.Add("/baz/:p2/three")

	match, pattern, params = matcher.Match("/foo")
	if match {
		t.Fail()
	}
	match, pattern, params = matcher.Match("/foo/bar/baz")
	if !match {
		t.Fail()
	}
	if pattern != "/foo/bar/baz" {
		t.Fail()
	}
	if len(params) != 0 {
		t.Fail()
	}
	match, pattern, params = matcher.Match("/baz/some/two")
	if !match {
		t.Fail()
	}
	if pattern != "/baz/:p1/two" {
		t.Fail()
	}
	if len(params) != 1 || params["p1"] != "some" {
		t.Fail()
	}
	match, pattern, params = matcher.Match("/baz/some/three")
	if !match {
		t.Fail()
	}
	if pattern != "/baz/:p2/three" {
		t.Fail()
	}
	if len(params) != 1 || params["p2"] != "some" {
		t.Fail()
	}
	match, pattern, params = matcher.Match("/baz/some/two/three")
	if match {
		t.Fail()
	}
	match, pattern, params = matcher.Match("/baz/some/thing")
	if match {
		t.Fail()
	}
}
