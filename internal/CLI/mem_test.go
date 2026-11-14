package cli

import "testing"

func TestMemoryBarBounds(t *testing.T) {
	c := &MemoryConfig{BarWidth: 10, BarCharFilled: "-", BarCharEmpty: "="}
	for _, tc := range []struct {
		used, total int64
		want        string
	}{{5, 10, "[-----=====]"}, {20, 10, "[----------]"}, {-5, 10, "[==========]"}, {1, 0, ""}} {
		if got := generateBar(tc.used, tc.total, c); got != tc.want {
			t.Fatalf("%d/%d: %q != %q", tc.used, tc.total, got, tc.want)
		}
	}
}
