package condense

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// handleFuncDecl processes function declarations.
func handleFuncDecl(pass *analysis.Pass, f *ast.File, fd *ast.FuncDecl, maxLen int) bool {
	// Skip function declarations without names (shouldn't happen, but be safe)
	if fd.Name == nil {
		return true
	}

	fd.Doc = nil // Clear the doc comment as it can confuse the output.

	var edits []analysis.TextEdit

	if edit, ok := condenseFieldList(pass.Fset, f, fd.Type.TypeParams); ok {
		edits = append(edits, edit)
	}

	if edit, ok := condenseFieldList(pass.Fset, f, fd.Type.Params); ok {
		edits = append(edits, edit)
	}

	if edit, ok := condenseFieldList(pass.Fset, f, fd.Type.Results); ok {
		edits = append(edits, edit)
	}

	if edit, ok := condenseNode(pass.Fset, f, fd.Body, maxLen); ok {
		edits = append(edits, edit)
	}

	if len(edits) > 0 {
		pass.Report(analysis.Diagnostic{
			Pos:            fd.Pos(),
			End:            fd.End(),
			Message:        "condense function declaration",
			SuggestedFixes: []analysis.SuggestedFix{{Message: "Condense function declaration", TextEdits: edits}},
		})
		return false
	}

	return true
}
