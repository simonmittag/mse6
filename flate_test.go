package mse6

import "testing"

func TestDeflateAndReInflate(t *testing.T) {
	const want = "MaryHadALittleLamb"
	got := Inflate(*Deflate([]byte(want)))
	if string(*got) != want {
		t.Errorf("not reinflated. want %v, got %v", want, string(*got))
	}
}
