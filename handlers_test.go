package mse6

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlers(t *testing.T) {
	tests := []struct {
		h ServerHandler
		b bool
		r int
	}{
		{h: ServerHandler{Methods: []string{"GET"}, Pattern: Prefix + "/badcontentlength", Handler: badcontentlength}, b: true, r: 200},
		{h: ServerHandler{Methods: []string{"GET"}, Pattern: Prefix + "/badgzip", Handler: badgzipf}, b: true, r: 200},
		{h: ServerHandler{Methods: []string{"CONNECT"}, Pattern: Prefix + "/connect", Handler: connect}, b: false, r: 200},
	}

	for _, tt := range tests {
		t.Run(tt.h.Pattern, func(t *testing.T) {
			doTestHandler(t, tt.h, tt.b, tt.r)
		})
	}
}

func doTestHandler(t *testing.T, h ServerHandler, bodyParsingErrWant bool, r int) {
	srv := httptest.NewServer(http.HandlerFunc(h.Handler))
	defer srv.Close()

	for _, m := range h.Methods {
		client := http.Client{}
		req, _ := http.NewRequest(m, srv.URL, nil)
		res, err := client.Do(req)
		if err != nil {
			t.Errorf("server did not return ok cause %v", err)
		}

		_, err2 := ioutil.ReadAll(res.Body)
		res.Body.Close()
		bodyParsingErrGet := err2 != nil
		if bodyParsingErrWant != bodyParsingErrGet {
			t.Errorf("body parsing err want %v got %v", bodyParsingErrWant, bodyParsingErrGet)
		}

		if res.StatusCode != r {
			t.Errorf("response status code want %v, got %v", r, res.StatusCode)
		}
	}
}
