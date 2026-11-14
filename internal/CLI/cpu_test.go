package cli

import "testing"

func TestPhysicalCoreCount(t *testing.T) {
	data := "processor : 0\nphysical id : 0\ncore id : 0\n\nprocessor : 1\nphysical id : 0\ncore id : 0\n\nprocessor : 2\nphysical id : 1\ncore id : 0\n"
	if got := physicalCoreCount(data); got != 2 {
		t.Fatalf("expected two physical cores across sockets, got %d", got)
	}
	if got := physicalCoreCount("processor : 0\n\nprocessor : 1\n"); got != 2 {
		t.Fatalf("ARM fallback: got %d", got)
	}
}
