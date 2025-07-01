// Package condense provides an analyzer that condenses struct literal
// declarations.
package condense

import (
	"bytes"
	"flag"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// NewAnalyzer creates a new analysis.Analyzer that condenses struct literal
// declarations.
func NewAnalyzer() *analysis.Analyzer {
	analyzer := &analysis.Analyzer{
		Name: "condense",
		Doc:  "Condense struct literal declarations to a single line if they fit within the specified maximum line length.",
		Run:  run,
	}
	analyzer.Flags.Int("max-len", 120, "maximum line length for collapsing declarations")
	return analyzer
}

func run(pass *analysis.Pass) (any, error) {
	maxLen, _ := pass.Analyzer.Flags.Lookup("max-len").Value.(flag.Getter).Get().(int)

	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			cl, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}

			if hasLineComment(cl, f) {
				return true
			}
			var sb bytes.Buffer
			_ = printer.Fprint(&sb, pass.Fset, cl)
			line, ok := collapseLine(sb.Bytes())
			if !ok {
				return false
			}

			var edits []analysis.TextEdit
			if len(line) <= maxLen {
				edits = append(edits, analysis.TextEdit{Pos: cl.Pos(), End: cl.End(), NewText: line})
			}

			if len(line) <= maxLen {
				pass.Report(analysis.Diagnostic{
					Pos:            cl.Pos(),
					End:            cl.End(),
					Message:        "condense declaration onto a single line",
					SuggestedFixes: []analysis.SuggestedFix{{Message: "Condense declaration", TextEdits: edits}},
				})
				return false
			}

			// Large slice: collapse elements individually and ensure line breaks around each.
			arrayType, isSlice := cl.Type.(*ast.ArrayType)
			if isSlice && arrayType != nil {
				var edits []analysis.TextEdit

				// Insert newline after opening brace if needed
				if e, ok := insertNewline(pass.Fset, cl.Lbrace, cl.Elts[0].Pos()); ok {
					edits = append(edits, e)
				}

				// Condense each element and insert newlines between them
				for i, elt := range cl.Elts {
					if hasLineComment(elt, f) {
						continue
					}

					var eltBuf bytes.Buffer
					_ = printer.Fprint(&eltBuf, pass.Fset, elt)
					collapsedElt, ok := collapseLine(eltBuf.Bytes())
					if !ok {
						continue
					}

					if len(collapsedElt) <= maxLen {
						edits = append(edits, analysis.TextEdit{
							Pos:     elt.Pos(),
							End:     elt.End(),
							NewText: collapsedElt,
						})
					}

					// Insert newline after each element explicitly if next element is inline
					if i < len(cl.Elts)-1 {
						if e, ok := insertNewline(pass.Fset, elt.End(), cl.Elts[i+1].Pos()); ok {
							edits = append(edits, e)
						}
					}
				}

				// Insert newline before closing brace if last element inline
				lastElt := cl.Elts[len(cl.Elts)-1]
				if pass.Fset.Position(lastElt.End()).Line == pass.Fset.Position(cl.Rbrace).Line {
					edits = append(edits, analysis.TextEdit{Pos: cl.Rbrace, End: cl.Rbrace, NewText: []byte(",\n")})
				}

				if len(edits) > 0 {
					pass.Report(analysis.Diagnostic{
						Pos:     cl.Pos(),
						End:     cl.End(),
						Message: "condense element declarations onto a single line",
						SuggestedFixes: []analysis.SuggestedFix{
							{Message: "Condense elements individually", TextEdits: edits},
						},
					})
				}
				return false
			}

			return true
		})
	}

	return nil, nil //nolint:nilnil
}

func hasLineComment(node ast.Node, file *ast.File) bool {
	for _, group := range file.Comments {
		if group.Pos() >= node.Pos() && group.End() <= node.End() {
			for _, comment := range group.List {
				if strings.HasPrefix(comment.Text, "//") &&
					(!testing.Testing() || !strings.HasPrefix(comment.Text, "// want")) {
					return true
				}
			}
		}
	}
	return false
}

func collapseLine(b []byte) ([]byte, bool) {
	if bytes.IndexByte(b, '\n') == -1 {
		return nil, false
	}

	out := make([]byte, 0, len(b))
	i := 0

	// skip leading whitespace
	for ; i < len(b) && (b[i] == ' ' || b[i] == '\n' || b[i] == '\r' || b[i] == '\t'); i++ {
		out = append(out, b[i])
	}

	var last byte
	for ; i < len(b); i++ {
		c := b[i]

		// collapse whitespace
		if c == ' ' || c == '\n' || c == '\r' || c == '\t' {
			j := i + 1
			for j < len(b) && (b[j] == ' ' || b[j] == '\n' || b[j] == '\r' || b[j] == '\t') {
				j++
			}
			if j < len(b) && (b[j] == '}' || b[j] == ']') {
				continue
			}
			if last != ' ' && last != '{' && last != 0 {
				out = append(out, ' ')
				last = ' '
			}
			continue
		}

		// skip trailing commas before } or ]
		if c == ',' {
			j := i + 1
			for j < len(b) && (b[j] == ' ' || b[j] == '\n' || b[j] == '\r' || b[j] == '\t') {
				j++
			}
			if j < len(b) && (b[j] == '}' || b[j] == ']') {
				continue
			}
		}

		out = append(out, c)
		last = c
	}

	return out, true
}

func insertNewline(fset *token.FileSet, pos1, pos2 token.Pos) (analysis.TextEdit, bool) {
	if fset.Position(pos1).Line == fset.Position(pos2).Line {
		return analysis.TextEdit{Pos: pos2, End: pos2, NewText: []byte("\n")}, true
	}
	return analysis.TextEdit{}, false
}
