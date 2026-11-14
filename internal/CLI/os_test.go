package cli

import "testing"

func TestOSProbeDispatch(t *testing.T) {
	for _, name := range []string{Linux, Darwin, FreeBSD, OpenBSD, NetBSD, DragonFly} {
		if got := GetOS(name, ""); got != name {
			t.Fatalf("%s probes would be dispatched as %s", name, got)
		}
	}
	if got := GetOS("Windows_NT", ""); got != Windows {
		t.Fatalf("Windows dispatch: %s", got)
	}
	if got := GetOS("UnknownOS", ""); got != "UnknownOS" {
		t.Fatalf("unknown OS should degrade gracefully: %s", got)
	}
}
