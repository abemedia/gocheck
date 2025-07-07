package test

import "testing"

func TestExportedWithTest(t *testing.T) {
	result := ExportedWithTest()
	if result != "tested" {
		t.Errorf("expected 'tested', got %s", result)
	}
}

func TestExportedMethodWithTest(t *testing.T) {
	m := MyStruct{}
	m.ExportedMethodWithTest() // This references the method
}

func TestExportedPointerMethodWithTest(t *testing.T) {
	m := &MyStruct{}
	m.ExportedPointerMethodWithTest() // This references the pointer method
}
