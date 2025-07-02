package condense

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// handleCallExpr processes function calls.
func handleCallExpr(pass *analysis.Pass, f *ast.File, ce *ast.CallExpr, maxLen int) bool {
	// Skip function calls that are immediate invocations of function literals
	if _, isFuncLit := ce.Fun.(*ast.FuncLit); isFuncLit {
		return true
	}

	edit, ok := condenseNode(pass.Fset, f, ce, maxLen)
	if !ok {
		return true
	}

	pass.Report(analysis.Diagnostic{
		Pos:     ce.Pos(),
		End:     ce.End(),
		Message: "condense call expression",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message:   "Condense call expression",
				TextEdits: []analysis.TextEdit{edit},
			},
		},
	})
	return false
}
