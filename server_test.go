package mse6

import "testing"

func TestHttpClientSocketTimeout(t *testing.T) {
	go Bootstrap(65534, 1)
}