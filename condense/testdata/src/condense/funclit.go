package condense

func _() {
	// Anonymous function literals
	f1 := func() { println("hello") }

	// condense function literal onto single line
	f2 := func() { // want "condense function literal"
		println("hello")
	}

	// condense function literal with parameters
	f3 := func( // want "condense function literal"
		x int,
		y string,
	) {
		println(x, y)
	}

	// condense function literal with return values
	f4 := func() ( // want "condense function literal"
		int,
		string,
	) {
		return 42, "test"
	}

	// don't condense multi-statement function literals
	f5 := func() {
		x := 42
		println(x)
	}

	// complex function literal with multiple parameters and returns
	f6 := func( // want "condense function literal"
		a int,
		b string,
		c bool,
	) (
		result int,
		message string,
		ok bool,
	) {
		return a, b, c
	}

	// function literal with variadic parameters
	f7 := func( // want "condense function literal"
		prefix string,
		values ...int,
	) string {
		return prefix
	}

	// function literal in variable assignment with type
	var f8 func(int) int = func( // want "condense function literal"
		x int,
	) int {
		return x * 2
	}

	// function literal as immediate call
	result := func( // want "condense function literal"
		x,
		y int,
	) int {
		return x + y
	}(1, 2)

	// function literal with named return parameters (should condense now)
	f9Named := func( // want "condense function literal"
		input string,
	) (
		output string,
	) {
		return input
	}

	// function literal with named return parameters (single statement - should condense)
	f9 := func( // want "condense function literal"
		input string,
	) string {
		return input
	}

	// function literal that captures variables
	captured := 100
	f10 := func( // want "condense function literal"
		x int,
	) int {
		return x + captured
	}

	// don't condense function literals with comments inside
	f11 := func() {
		// this has a comment
		println("hello")
	}

	// don't condense function literals with line comments next to arguments
	f12 := func(
		x int, // this is x
		y string, // this is y
	) int {
		return x
	}

	// function literal with line comments on returns only (currently doesn't condense, but ideally should condense arguments)
	f13LineCommentsOnReturns := func( // want "condense function literal"
		x int,
		y string,
	) (
		result int, // this is the result
		err error, // this is the error
	) {
		return x, nil
	}

	// function literal with line comments on arguments only (currently doesn't condense, but ideally should condense returns)
	f14LineCommentsOnArgs := func( // want "condense function literal"
		x int, // this is x
		y string, // this is y
	) (
		int,
		error,
	) {
		return x, nil
	}

	// function literal with line comments on both (should not condense anything)
	f15LineCommentsOnBoth := func(
		x int, // this is x
		y string, // this is y
	) (
		result int, // this is the result
		err error, // this is the error
	) {
		return x, nil
	}

	// function literal returning function literal
	f16 := func() func() { // want "condense function literal"
		return func() {
			println("nested")
		}
	}

	// use the functions to avoid unused variable errors
	f1()
	f2()
	f3(1, "test")
	_, _ = f4()
	f5()
	_, _, _ = f6(1, "test", true)
	_ = f7("prefix", 1, 2, 3)
	_ = f8(5)
	_ = result
	_ = f9Named("test")
	_ = f9("test")
	_ = f10(10)
	f11()
	_ = f12(1, "test")
	_, _ = f13LineCommentsOnReturns(1, "test")
	_, _ = f14LineCommentsOnArgs(1, "test")
	_, _ = f15LineCommentsOnBoth(1, "test")
	_ = f16()
}
