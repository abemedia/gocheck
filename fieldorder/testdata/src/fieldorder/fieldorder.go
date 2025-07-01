package fieldorder

func _() {
	type Person struct {
		Name    string
		Age     int
		Email   string
		Address string
	}

	type Simple struct {
		A string
		B int
	}

	// correct order - no change needed
	_ = Person{
		Name:    "John",
		Age:     30,
		Email:   "john@example.com",
		Address: "123 Main St",
	}

	// out of order - should be reordered
	_ = Person{ // want "struct literal fields are out of order"
		Age:     30,
		Name:    "John",
		Email:   "john@example.com",
		Address: "123 Main St",
	}

	// multiple fields out of order
	_ = Person{ // want "struct literal fields are out of order"
		Email:   "john@example.com",
		Name:    "John",
		Address: "123 Main St",
		Age:     30,
	}

	// multiple fields out of order on shared line
	_ = Person{ // want "struct literal fields are out of order"
		Email: "john@example.com", Name: "John",
		Address: "123 Main St", Age: 30,
	}

	// partial fields in wrong order
	_ = Person{ // want "struct literal fields are out of order"
		Email: "john@example.com",
		Name:  "John",
	}

	// single field - no issue
	_ = Person{
		Name: "John",
	}

	// empty struct literal - no issue
	_ = Person{}

	// anonymous struct with correct order
	_ = struct {
		X int
		Y string
	}{
		X: 1,
		Y: "test",
	}

	// anonymous struct with wrong order
	_ = struct {
		X int
		Y string
	}{ // want "struct literal fields are out of order"
		Y: "test",
		X: 1,
	}

	// nested struct - only outer should be flagged
	_ = Person{ // want "struct literal fields are out of order"
		Age:  30,
		Name: "John",
		Address: func() string {
			inner := Simple{A: "nested", B: 123} // correct order
			return inner.A
		}(),
	}

	// struct with embedded types
	type Embedded struct {
		ID int
	}

	type WithEmbedded struct {
		Embedded
		Name string
		Age  int
	}

	// correct order with embedded field
	_ = WithEmbedded{
		Embedded: Embedded{ID: 1},
		Name:     "John",
		Age:      30,
	}

	// wrong order with embedded field
	_ = WithEmbedded{ // want "struct literal fields are out of order"
		Name:     "John",
		Embedded: Embedded{ID: 1},
		Age:      30,
	}
}
