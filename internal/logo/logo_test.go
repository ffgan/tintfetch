package logo

import "testing"

func TestRegistryLongestPrefix(t *testing.T) {
	tree := registry{"Linux": 1, "Linux Mint": 2, "linux mint": 3}
	for i := 0; i < 100; i++ {
		key, value, ok := tree.LongestPrefix("LINUX MINT 22")
		if !ok || key != "Linux Mint" || value != 2 {
			t.Fatalf("non-deterministic prefix match: %s %v", key, value)
		}
	}
}
func TestPaletteEncoding(t *testing.T) {
	if got := color(4); got != "\x1b[38;5;4m" {
		t.Fatalf("palette index converted incorrectly: %q", got)
	}
	if color(256) != Reset {
		t.Fatal("invalid palette index accepted")
	}
}
