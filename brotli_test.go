package mse6

import "testing"

func TestBrotliEncodeAndDecode(t *testing.T) {
	want := "ResistanceIsFutile"
	got := BrotliDecode(*BrotliEncode([]byte(want)))

	if string(*got) != want {
		t.Errorf("brotli decode/encode failed. want %v got %v", want, got)
	}
}
