package condense

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// handleCompositeLit processes composite literals (struct/slice/array literals).
func handleCompositeLit(pass *analysis.Pass, f *ast.File, cl *ast.CompositeLit, maxLen int) bool {
	edit, ok := condenseNode(pass.Fset, f, cl, maxLen)
	if ok {
		pass.Report(analysis.Diagnostic{
			Pos:            cl.Pos(),
			End:            cl.End(),
			Message:        "condense declaration",
			SuggestedFixes: []analysis.SuggestedFix{{Message: "Condense declaration", TextEdits: []analysis.TextEdit{edit}}},
		})
		return false
	}

	// Handle large slices by condensing elements individually
	if arrayType, isSlice := cl.Type.(*ast.ArrayType); isSlice && arrayType != nil {
		return handleSliceElements(pass, f, cl, maxLen)
	}

	return true
}

// handleSliceElements processes slice literals with individual element condensing.
func handleSliceElements(pass *analysis.Pass, f *ast.File, cl *ast.CompositeLit, maxLen int) bool {
	edits := make([]analysis.TextEdit, 0, len(cl.Elts)*2+2)

	// Condense each element and insert newlines between them
	for i, elt := range cl.Elts {
		edit, ok := condenseNode(pass.Fset, f, elt, maxLen)
		if !ok {
			continue
		}

		// Insert newline after opening brace if needed
		if i == 0 {
			if e, ok := insertNewline(pass.Fset, cl.Lbrace, elt.Pos()); ok {
				edits = append(edits, e)
			}
		}

		edits = append(edits, edit)

		// Insert newline after each element explicitly if next element is inline
		if i < len(cl.Elts)-1 {
			if e, ok := insertNewline(pass.Fset, elt.End(), cl.Elts[i+1].Pos()); ok {
				edits = append(edits, e)
			}
		}

		// Insert newline before closing brace if last element inline
		if i == len(cl.Elts)-1 {
			if pass.Fset.Position(elt.End()).Line == pass.Fset.Position(cl.Rbrace).Line {
				edits = append(edits, analysis.TextEdit{Pos: cl.Rbrace, End: cl.Rbrace, NewText: []byte(",\n")})
			}
		}
	}

	if len(edits) > 0 {
		pass.Report(analysis.Diagnostic{
			Pos:     cl.Pos(),
			End:     cl.End(),
			Message: "condense declaration",
			SuggestedFixes: []analysis.SuggestedFix{
				{Message: "Condense elements individually", TextEdits: edits},
			},
		})
	}

	return false
}
