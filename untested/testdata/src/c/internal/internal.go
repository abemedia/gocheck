package test

// InternalFunction should trigger warning when internal flag is enabled
func InternalFunction() { // want "exported function \"InternalFunction\" has no test"
	// internal function
}

// InternalFunctionWithTest should not trigger warning when internal flag is enabled
func InternalFunctionWithTest() {
	// internal function with test
}
