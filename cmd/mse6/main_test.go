package main

import (
	"os"
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

func TestMainFunc(t *testing.T) {
	os.Setenv("LOGCOLOR", "TRUE")
	os.Args = append([]string{"-v"}, "-v")
	main()
}

// tests below only need to execute to work out
func TestPrintVersion(t *testing.T) {
	printVersion()
}

func TestPrintSelfTest(t *testing.T) {
	printSelftest(61234)
}

func TestInitLogger(t *testing.T) {
	initLogger()
}
