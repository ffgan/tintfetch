package cli

import "testing"

func TestFormatPackages(t *testing.T) {
	packages := []PkgInfo{{"dpkg", 100}, {"flatpak-user", 3}}
	for _, tc := range []struct{ mode, want string }{{"on", "100 (dpkg), 3 (flatpak-user)"}, {"tiny", "103 (dpkg, flatpak-user)"}, {"off", "103"}} {
		if got := formatPackages(packages, tc.mode); got != tc.want {
			t.Fatalf("%s: %q != %q", tc.mode, got, tc.want)
		}
	}
	if got := formatPackages(nil, "on"); got != "" {
		t.Fatalf("missing packages should be omitted: %q", got)
	}
}
