package dbstore

import (
	"testing"
)

func Test_encodeRFID_SingleLeadingZero(t *testing.T) {
	id := "0101436029"
	want := "7dca0b06"
	result := encodeRFID(id)
	if want != result {
		t.Errorf(`encodeRFID(%q) = %q, want %q`, id, result, want)
	}
}

func Test_encodeRFID_EmbeddedLeadingZero(t *testing.T) {
	id := "2949385415"
	want := "c70cccaf"
	result := encodeRFID(id)
	if want != result {
		t.Errorf(`encodeRFID(%q) = %q, want %q`, id, result, want)
	}
}

func Test_encodeRFID_MultipleLeadingZeros(t *testing.T) {
	id := "0097304109"
	want := "2dbecc05"
	result := encodeRFID(id)
	if want != result {
		t.Errorf(`encodeRFID(%q) = %q, want %q`, id, result, want)
	}
}
