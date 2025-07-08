// Package skip provides efficient file skipping strategies for Go static analysis tools.
package skip

import (
	"cmp"
	"go/ast"
	"go/token"
	"slices"

	"golang.org/x/tools/go/analysis"
)

type (
	// FileFilter is a function that determines whether a file should be skipped.
	FileFilter func(*ast.File) bool

	// NodeFilter is a function that determines whether a node should be skipped.
	NodeFilter func(ast.Node) bool
)

// NewFileStrategy creates a node filter that skips nodes from files matching the given filter.
func NewFileStrategy(pass *analysis.Pass, filter FileFilter) NodeFilter {
	skip := make([]fileRange, 0, len(pass.Files))
	for _, file := range pass.Files {
		if filter(file) {
			skip = append(skip, fileRange{start: file.Pos(), end: file.End()})
		}
	}

	slices.SortFunc(skip, func(a, b fileRange) int { return cmp.Compare(a.start, b.start) })

	return func(node ast.Node) bool {
		_, found := slices.BinarySearchFunc(skip, node.Pos(), func(r fileRange, pos token.Pos) int {
			if pos < r.start {
				return +1
			}
			if pos > r.end {
				return -1
			}
			return 0
		})
		return found
	}
}

type fileRange struct{ start, end token.Pos }
