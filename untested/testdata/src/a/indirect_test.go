package test

import (
	"io"
	"testing"
)

func TestExportedDirect(t *testing.T) {
	// Direct call
	result := ExportedDirect()
	if result != "direct" {
		t.Errorf("got %s, want direct", result)
	}
}

func TestExportedIndirect(t *testing.T) {
	// Indirect call through helper
	runHelper(t)
}

func runHelper(t *testing.T) {
	result := ExportedIndirect()
	if result != "indirect" {
		t.Errorf("got %s, want indirect", result)
	}
}

func TestExportedViaAssign(t *testing.T) {
	// Call after assignment to interface
	tests := []struct {
		name string
		r    io.Reader
	}{
		{name: "MyReader", r: MyReader{}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			buf := make([]byte, 10)
			test.r.Read(buf)
		})
	}
}

func TestExportedChained(t *testing.T) {
	// Test Bar which calls ExportedChained
	result := CallChained()
	if result != "chained" {
		t.Errorf("got %s, want chained", result)
	}
}
