package condense

func _() {
	// no change on single line
	_, _, _ = test("test", 1, true)

	// condense onto a single line
	_, _, _ = test( // want "condense call expression"
		"test",
		1,
		true,
	)

	// don't condense function calls with line comments next to arguments
	_, _, _ = test(
		"test", // this is a string
		1,      // this is an int
		true,   // this is a bool
	)
}
