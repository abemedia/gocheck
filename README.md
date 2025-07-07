# gocheck

<img align="right" width="200" alt="" src="assets/logo.png">

[![Go Reference](https://pkg.go.dev/badge/github.com/abemedia/gocheck.svg)](https://pkg.go.dev/github.com/abemedia/gocheck)
[![Codecov](https://codecov.io/gh/abemedia/gocheck/branch/master/graph/badge.svg)](https://codecov.io/gh/abemedia/gocheck)
[![CI](https://github.com/abemedia/gocheck/actions/workflows/test.yml/badge.svg)](https://github.com/abemedia/gocheck/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/abemedia/gocheck)](https://goreportcard.com/report/github.com/abemedia/gocheck)

A collection of Go static analysis tools to improve code quality and consistency.

## Linters

### `fieldorder`

Ensures struct literal fields are in the same order as their type declaration, promoting consistency and readability.

<details>
<summary>More details</summary>

#### Example

```go
type Person struct {
    Name string
    Age  int
    City string
}

// ❌ Fields out of order
person := Person{
    Age:  30,
    Name: "John",  // Should come before Age
    City: "NYC",
}

// ✅ Fields in correct order
person := Person{
    Name: "John",
    Age:  30,
    City: "NYC",
}
```

</details>

### `untested`

Checks that exported functions and methods have corresponding tests, including indirect testing through helper functions.

<details>
<summary>More details</summary>

#### Example

```go
// main.go
package main

// ✅ Function has tests
func ExportedFunction() string {
    return "hello"
}

// ❌ Function has no tests
func ExportedUntested() string {
    return "not tested"
}

// main_test.go
package main

import "testing"

func TestExportedFunction(t *testing.T) {
    result := ExportedFunction()
    if result != "hello" {
        t.Errorf("expected hello, got %s", result)
    }
}
```

</details>

## Installation

Install the latest version:

```bash
go install github.com/abemedia/gocheck/cmd/gocheck@latest
```

## Usage

### Basic Usage

Run all linters:

```bash
gocheck ./...
```

Run specific linters:

```bash
gocheck -analyzers=untested,fieldorder ./...
```

### Available Flags

| Flag                  | Description                              | Default                 |
| --------------------- | ---------------------------------------- | ----------------------- |
| `-analyzers`          | Comma-separated list of analyzers to run | All available analyzers |
| `-untested.internal`  | Check functions in internal packages     | `false`                 |
| `-untested.generated` | Check functions in generated files       | `false`                 |

### Examples

Run only the `untested` linter with default settings:

```bash
gocheck -analyzers=untested ./...
```

Run `untested` linter including both internal packages and generated files:

```bash
gocheck -analyzers=untested -untested.internal=true -untested.generated=true ./...
```

Run both linters:

```bash
gocheck -analyzers=untested,fieldorder ./...
```
