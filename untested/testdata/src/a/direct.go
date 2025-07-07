package test

// ExportedWithTest should not trigger warning (has test)
func ExportedWithTest() string {
	return "tested"
}

// ExportedWithoutTest should trigger warning (no test)
func ExportedWithoutTest() string { // want "exported function \"ExportedWithoutTest\" has no test"
	return "untested"
}

// unexportedFunction should not trigger warning (not exported)
func unexportedFunction() string {
	return "private"
}

type MyStruct struct{}

// ExportedMethod should trigger warning (no test)
func (m MyStruct) ExportedMethod() { // want "exported method \"MyStruct.ExportedMethod\" has no test"
	// method code
}

// ExportedMethodWithTest should not trigger warning (has test)
func (m MyStruct) ExportedMethodWithTest() {
	// method code
}

// ExportedPointerMethod should trigger warning (no test)
func (m *MyStruct) ExportedPointerMethod() { // want "exported method \"MyStruct.ExportedPointerMethod\" has no test"
	// pointer method code
}

// ExportedPointerMethodWithTest should not trigger warning (has test)
func (m *MyStruct) ExportedPointerMethodWithTest() {
	// pointer method code
}
