# gocheck

<img align="right" width="160" alt="" src="assets/logo.png">

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
// example.go
package example

// ✅ Function has tests
func ExportedFunction() string {
    return "hello"
}

// ❌ Function has no tests
func ExportedUntested() string {
    return "not tested"
}

// example_test.go
package example

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

Run all linters:

```bash
gocheck ./...
```

### Available Flags

> [!NOTE]
> When you explicitly enable one analyzer (e.g., `-fieldorder`), it disables others unless they're also explicitly enabled.

| Flag                  | Description                                          | Default |
| --------------------- | ---------------------------------------------------- | ------- |
| `-fieldorder`         | Enable fieldorder analysis                           | `true`  |
| `-untested`           | Enable untested analysis                             | `true`  |
| `-untested.internal`  | Check functions in internal packages                 | `false` |
| `-untested.generated` | Check functions in generated files                   | `false` |
| `-fix`                | Apply all suggested fixes                            | `false` |
| `-json`               | Emit JSON output                                     | `false` |
| `-test`               | Indicates whether test files should be analyzed, too | `true`  |

### Examples

Run only the `fieldorder` linter:

```bash
gocheck -fieldorder ./...
```

Run the `untested` linter including both internal packages and generated files:

```bash
gocheck -untested -untested.internal -untested.generated ./...
```

Show available options:

```bash
gocheck help
```
