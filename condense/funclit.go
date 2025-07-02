package condense

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// handleFuncLit processes function literals.
func handleFuncLit(pass *analysis.Pass, f *ast.File, fl *ast.FuncLit, maxLen int) bool {
	var edits []analysis.TextEdit

	if edit, ok := condenseNode(pass.Fset, f, fl, maxLen); ok && fl.Body != nil && len(fl.Body.List) <= 1 {
		edits = append(edits, edit)
	} else {
		if edit, ok := condenseFieldList(pass.Fset, f, fl.Type.TypeParams); ok {
			edits = append(edits, edit)
		}
		if edit, ok := condenseFieldList(pass.Fset, f, fl.Type.Params); ok {
			edits = append(edits, edit)
		}
		if edit, ok := condenseFieldList(pass.Fset, f, fl.Type.Results); ok {
			edits = append(edits, edit)
		}
	}

	if len(edits) > 0 {
		pass.Report(analysis.Diagnostic{
			Pos:     fl.Pos(),
			End:     fl.End(),
			Message: "condense function literal",
			SuggestedFixes: []analysis.SuggestedFix{
				{Message: "Condense function literal", TextEdits: edits},
			},
		})
		return false
	}

	return true
}
