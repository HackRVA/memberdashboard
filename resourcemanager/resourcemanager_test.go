package resourcemanager

import (
	"testing"
)

// TestHash test the hash function
func TestHash(t *testing.T) {
	reg := []string{"a", "b", "c"}
	h1 := hash(reg)
	h2 := hash([]string{"a", "b", "c", "d"})
	println(h1 == h2)
	if len(h1) > 0 {
		t.Errorf("%s\n", h1)
	}
}
