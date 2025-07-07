// Package fieldorder checks that struct literal fields are in the same order
// as the type declaration.
package fieldorder

import (
	"bytes"
	"cmp"
	"go/ast"
	"go/printer"
	"go/types"
	"slices"

	"golang.org/x/tools/go/analysis"
)

// NewAnalyzer creates a new analysis.Analyzer that checks struct literal
// fields are in the same order as the type declaration.
func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "fieldorder",
		Doc:  "check that struct literal fields are in the same order as the type declaration",
		Run:  run,
	}
}

//nolint:funlen,gocognit
func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if ast.IsGenerated(file) {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			cl, ok := n.(*ast.CompositeLit)
			if !ok || cl.Type == nil || len(cl.Elts) == 0 {
				return true
			}

			typ := pass.TypesInfo.TypeOf(cl.Type)
			if typ == nil {
				return true
			}

			st, ok := typ.Underlying().(*types.Struct)
			if !ok {
				return true
			}

			// Build declared field order map
			fieldOrder := make(map[string]int)
			for i := range st.NumFields() {
				fieldOrder[st.Field(i).Name()] = i
			}
			// Check if fields are in the correct order
			needsReordering := false
			keyValueExprs := make([]*ast.KeyValueExpr, 0, len(cl.Elts))

			for _, elt := range cl.Elts {
				kv, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}

				ident, ok := kv.Key.(*ast.Ident)
				if !ok {
					continue
				}

				if _, exists := fieldOrder[ident.Name]; !exists {
					continue
				}

				keyValueExprs = append(keyValueExprs, kv)
			}

			// Check if current order matches declared order
			for i := range len(keyValueExprs) - 1 {
				ident1 := keyValueExprs[i].Key.(*ast.Ident)
				ident2 := keyValueExprs[i+1].Key.(*ast.Ident)

				idx1 := fieldOrder[ident1.Name]
				idx2 := fieldOrder[ident2.Name]

				if idx1 > idx2 {
					needsReordering = true
					break
				}
			}

			var edits []analysis.TextEdit

			if needsReordering {
				// Sort keyValueExprs by field order
				sortedExprs := slices.Clone(keyValueExprs)
				slices.SortFunc(sortedExprs, func(a, b *ast.KeyValueExpr) int {
					return cmp.Compare(fieldOrder[a.Key.(*ast.Ident).Name], fieldOrder[b.Key.(*ast.Ident).Name])
				})

				// Create edits for each pair that needs swapping
				for i, unsorted := range keyValueExprs {
					sorted := sortedExprs[i]
					if unsorted != sorted {
						var buf bytes.Buffer

						_ = printer.Fprint(&buf, pass.Fset, sorted)
						edits = append(edits, analysis.TextEdit{Pos: unsorted.Pos(), End: unsorted.End(), NewText: buf.Bytes()})
					}
				}
			}

			if len(edits) > 0 {
				pass.Report(analysis.Diagnostic{
					Pos:      cl.Lbrace,
					End:      cl.Rbrace + 1,
					Category: "style",
					Message:  "struct literal fields are out of order",
					SuggestedFixes: []analysis.SuggestedFix{
						{Message: "Reorder fields to match declaration order", TextEdits: edits},
					},
				})
			}

			return true
		})
	}

	return nil, nil
}
