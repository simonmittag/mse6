package mse6

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHttpClientSocketTimeout(t *testing.T) {
	go Bootstrap(65534, 1, "/mse6/", false)
}

func TestGetResponds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(get))
	defer srv.Close()

	res, err := http.Get(srv.URL)
	if err != nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	body, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2 != nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}

	if res.StatusCode != 200 {
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

func TestPutResponds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(put))
	defer srv.Close()

	jsonData := map[string]string{"scandi": "grind", "convex": "grind", "concave": "grind"}
	jsonValue, _ := json.Marshal(jsonData)
	buf := bytes.NewBuffer(jsonValue)

	client := http.Client{}
	req, _ := http.NewRequest("PUT", srv.URL, buf)
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	_, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2 != nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("response status code want 200, got %v", res.StatusCode)
	}

	ce := res.Header["Content-Encoding"][0]
	if ce != "identity" {
		t.Errorf("response Content Encoding should be identity got %v", ce)
	}
}

func TestPatchResponds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(patch))
	defer srv.Close()

	jsonData := map[string]string{"scandi": "grind", "convex": "grind", "concave": "grind"}
	jsonValue, _ := json.Marshal(jsonData)
	buf := bytes.NewBuffer(jsonValue)

	client := http.Client{}
	req, _ := http.NewRequest("PATCH", srv.URL, buf)
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	_, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2 != nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("response status code want 200, got %v", res.StatusCode)
	}

	ce := res.Header["Content-Encoding"][0]
	if ce != "identity" {
		t.Errorf("response Content Encoding should be identity got %v", ce)
	}
}

func TestDeleteResponds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(delete))
	defer srv.Close()

	jsonData := map[string]string{"scandi": "grind", "convex": "grind", "concave": "grind"}
	jsonValue, _ := json.Marshal(jsonData)
	buf := bytes.NewBuffer(jsonValue)

	client := http.Client{}
	req, _ := http.NewRequest("DELETE", srv.URL, buf)
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	_, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2 != nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}

	if res.StatusCode != 204 {
		t.Errorf("response status code want 204, got %v", res.StatusCode)
	}
}

func TestTraceResponds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(trace))
	defer srv.Close()

	client := http.Client{}
	req, _ := http.NewRequest("TRACE", srv.URL, nil)
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	_, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2 != nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("response status code want 200, got %v", res.StatusCode)
	}

	ce := res.Header["Content-Type"][0]
	if ce != "message/http" {
		t.Errorf("response Content Encoding should be identity got %v", ce)
	}
}

func TestPostResponds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(post))
	defer srv.Close()

	rBody := []byte(`{"hello":"world"}`)
	res, err := http.Post(srv.URL, "application/json", bytes.NewBuffer(rBody))
	if err != nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	body, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2 != nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}

	if res.StatusCode != 201 {
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

func TestBadContentLengthResponds(t *testing.T) {
	waitDuration = 3
	srv := httptest.NewServer(http.HandlerFunc(badcontentlength))
	defer srv.Close()

	res, err := http.Get(srv.URL)
	if err != nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	_, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2 == nil {
		t.Errorf("body parsing did return nil error, want unexecpted EOF")
	}

	if res.StatusCode != 200 {
		t.Errorf("response status code want 200, got %v", res.StatusCode)
	}

	cl := res.Header["Content-Length"][0]
	if cl != "2048" {
		t.Errorf("response Content Length should be 2048 got %v", cl)
	}
}

func TestSlowHeaderResponds(t *testing.T) {
	//TODO: will not test the delay only the response because of httptest?
	waitDuration = 3
	srv := httptest.NewServer(http.HandlerFunc(slowheader))
	defer srv.Close()

	res, err := http.Get(srv.URL)
	if err != nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	body, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2 != nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}

	if res.StatusCode != 200 {
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
	if err != nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	body, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2 != nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("response status code want 200, got %v", res.StatusCode)
	}

	//TODO test golang res does not have gzip header after automatic decode.?

	sbody := string(body)
	if !strings.Contains(sbody, "mse6") {
		t.Errorf("invalid response, wanted mse6 json, got %v", sbody)
	}
}

func TestBadGzipResponds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(badgzipf))
	defer srv.Close()

	res, err := http.Get(srv.URL)
	if err != nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	_, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if !strings.Contains(err2.Error(), "gzip: invalid header") {
		t.Errorf("body parsing should fail due to invalid gzip header, does not: %v", err2)
	}

	if res.StatusCode != 200 {
		t.Errorf("response status code want 200, got %v", res.StatusCode)
	}
}

func TestSend404Responds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(send404))
	defer srv.Close()

	res, err := http.Get(srv.URL)
	if err != nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	body, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2 != nil {
		t.Errorf("body parsing did not return ok cause %v", err2)
	}
	sbody := string(body)
	if !strings.Contains(sbody, "404") {
		t.Errorf("invalid response, wanted mse6 json, got %v", sbody)
	}

	if res.StatusCode != 404 {
		t.Errorf("response status code want 200, got %v", res.StatusCode)
	}
}

func TestSendResponds(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(send))
	defer srv.Close()

	res, err := http.Get(srv.URL + "?code=201")
	if err != nil {
		t.Errorf("server did not return ok cause %v", err)
	}

	if res.StatusCode != 201 {
		t.Errorf("response status code want 201, got %v", res.StatusCode)
	}
}
