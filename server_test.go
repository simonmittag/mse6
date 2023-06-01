package mse6

import "testing"

func TestGetCert(t *testing.T) {
	_, e := getCert()
	if e != nil {
		t.Errorf("built-in tls cert not working")
	}
}
