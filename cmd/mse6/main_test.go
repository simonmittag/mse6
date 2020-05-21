package main

import (
	"testing"
)

func TestParsePrefix(t *testing.T) {
	want := "/"
	got := parsePrefix("")
	if want != got {
		t.Errorf("prefix parse error, want %v, got: %v", want, got)
	}

	want = "/mse6/"
	got = parsePrefix("/mse6")
	if want != got {
		t.Errorf("prefix parse error, want %v, got: %v", want, got)
	}

	want = "/mse6/"
	got = parsePrefix("/mse6/")
	if want != got {
		t.Errorf("prefix parse error, want %v, got: %v", want, got)
	}

	want = "/mse6/"
	got = parsePrefix("mse6/")
	if want != got {
		t.Errorf("prefix parse error, want %v, got: %v", want, got)
	}

	want = "/mse6/"
	got = parsePrefix("mse6")
	if want != got {
		t.Errorf("prefix parse error, want %v, got: %v", want, got)
	}
}
