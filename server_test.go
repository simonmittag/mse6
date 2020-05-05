package mse6

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHttpClientSocketTimeout(t *testing.T) {
	go Bootstrap(65534, 1)
}

func TestGetResponds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(get))
	defer srv.Close()

	res, err := http.Get(srv.URL)
	if err!=nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	body, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2!=nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}

	if res.StatusCode!=200 {
		t.Errorf("response status code want 200, got %v", res.StatusCode)
	}

	ce := res.Header["Content-Encoding"][0]
	if ce != "identity" {
		t.Errorf("response Content Encoding should be identity got %v", ce)
	}

	sbody := string(body)
	if !strings.Contains(sbody, "mse6") {
		t.Errorf("invalid response, wanted mse6 json, got %v", sbody)
	}
}

func TestPostResponds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(post))
	defer srv.Close()

	rBody := []byte(`{"hello":"world"}`)
	res, err := http.Post(srv.URL, "application/json", bytes.NewBuffer(rBody))
	if err!=nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	body, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2!=nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}

	if res.StatusCode!=200 {
		t.Errorf("response status code want 200, got %v", res.StatusCode)
	}

	ce := res.Header["Content-Encoding"][0]
	if ce != "identity" {
		t.Errorf("response Content Encoding should be identity got %v", ce)
	}

	sbody := string(body)
	if !strings.Contains(sbody, "mse6") {
		t.Errorf("invalid response, wanted mse6 json, got %v", sbody)
	}
}

func TestSlowBodyResponds(t *testing.T) {
	waitDuration = 3
	//TODO: will not test the delay only the response because of httptest?
	srv := httptest.NewServer(http.HandlerFunc(slowbody))
	defer srv.Close()

	res, err := http.Get(srv.URL)
	if err!=nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	body, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2!=nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}

	if res.StatusCode!=200 {
		t.Errorf("response status code want 200, got %v", res.StatusCode)
	}

	ce := res.Header["Content-Encoding"][0]
	if ce != "identity" {
		t.Errorf("response Content Encoding should be identity got %v", ce)
	}

	sbody := string(body)
	if !strings.Contains(sbody, "mse6") {
		t.Errorf("invalid response, wanted mse6 json, got %v", sbody)
	}
}

func TestSlowHeaderResponds(t *testing.T) {
	//TODO: will not test the delay only the response because of httptest?
	waitDuration = 3
	srv := httptest.NewServer(http.HandlerFunc(slowheader))
	defer srv.Close()

	res, err := http.Get(srv.URL)
	if err!=nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	body, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2!=nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}

	if res.StatusCode!=200 {
		t.Errorf("response status code want 200, got %v", res.StatusCode)
	}

	ce := res.Header["Content-Encoding"][0]
	if ce != "identity" {
		t.Errorf("response Content Encoding should be identity got %v", ce)
	}

	sbody := string(body)
	if !strings.Contains(sbody, "mse6") {
		t.Errorf("invalid response, wanted mse6 json, got %v", sbody)
	}
}

func TestGzipResponds(t *testing.T) {
	//TODO: will not test the delay only the response because of httptest?
	waitDuration = 3
	srv := httptest.NewServer(http.HandlerFunc(gzipf))
	defer srv.Close()

	res, err := http.Get(srv.URL)
	if err!=nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	body, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2!=nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}

	if res.StatusCode!=200 {
		t.Errorf("response status code want 200, got %v", res.StatusCode)
	}

	//TODO test golang res does not have gzip header after automatic decode.?

	sbody := string(body)
	if !strings.Contains(sbody, "mse6") {
		t.Errorf("invalid response, wanted mse6 json, got %v", sbody)
	}
}