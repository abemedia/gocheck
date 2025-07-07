package test

// ExportedDirect should be detected as tested (direct call)
func ExportedDirect() string {
	return "direct"
}

// ExportedIndirect should be detected as tested (called by helper)
func ExportedIndirect() string {
	return "indirect"
}

// ExportedViaAssign requires SSA analysis to detect interface assignment testing
type MyReader struct{}

func (m MyReader) Read(p []byte) (n int, err error) { // want "exported method \"MyReader.Read\" has no test"
	return 0, nil
}

// ExportedChained should be detected as tested (called by CallChained, CallChained tested)
func ExportedChained() string {
	return "chained"
}

func CallChained() string {
	return ExportedChained()
}

// ExportedUntested should be flagged as untested
func ExportedUntested() string { // want "exported function \"ExportedUntested\" has no test"
	return "untested"
}
