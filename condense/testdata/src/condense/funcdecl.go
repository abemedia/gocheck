package condense

// condense function declaration onto a single line
func test( // want "condense function declaration"
	a string,
	b int,
	c bool,
) (
	string,
	int,
	bool,
) {
	return a, b, c
}

// no change on single line
func _(a string, b int, c bool) (string, int, bool) { return a, b, c }

// condense function type parameters onto a single line
func _[ // want "condense function declaration"
	A string,
	B int,
	C bool,
](
	a string,
	b int, // this is an int
	c bool,
) (
	string,
	int, // this is an int
	bool,
) {
	return a, b, c
}

// condense function parameters onto a single line
func _( // want "condense function declaration"
	a string,
	b int,
	c bool,
) (
	string,
	int, // this is an int
	bool,
) {
	return a, b, c
}

// condense function returns onto a single line
func _( // want "condense function declaration"
	a string,
	b int, // this is an int
	c bool,
) (
	string,
	int,
	bool,
) {
	return a, b, c
}

// condense function parameters and results onto a single line
func _( // want "condense function declaration"
	a string,
	b int,
	c bool,
) (
	string,
	int,
	bool,
) {
	return a, b, c // result
}
