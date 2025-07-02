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
			if n == nil || pass.Fset.Position(n.Pos()).Line == pass.Fset.Position(n.End()).Line {
				return false // Skip processing nodes that are already on one line.
			}

			switch node := n.(type) {
			case *ast.CompositeLit:
				return handleCompositeLit(pass, f, node, maxLen)
			case *ast.FuncLit:
				return handleFuncLit(pass, f, node, maxLen)
			case *ast.FuncDecl:
				return handleFuncDecl(pass, f, node, maxLen)
			case *ast.CallExpr:
				return handleCallExpr(pass, f, node, maxLen)
			default:
				return true
			}
		})
	}

	return nil, nil //nolint:nilnil
}

func hasLineComment[T interface {
	ast.Node
	comparable
}](file *ast.File, nodes ...T) bool {
	var zero T
	for _, node := range nodes {
		if node == zero {
			continue
		}
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
	}
	return false
}

func condenseNode(fset *token.FileSet, f *ast.File, node ast.Node, maxLen int) (analysis.TextEdit, bool) {
	if node == nil || hasLineComment(f, node) {
		return analysis.TextEdit{}, false // Skip nodes with line comments
	}

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		return analysis.TextEdit{}, false
	}

	line, ok := condense(buf.Bytes())
	if !ok {
		return analysis.TextEdit{}, false
	}

	if len(line) > maxLen {
		return analysis.TextEdit{}, false
	}

	return analysis.TextEdit{Pos: node.Pos(), End: node.End(), NewText: line}, true
}

// condenseFieldList condenses a field list (parameters or returns) onto a single line.
func condenseFieldList(fset *token.FileSet, f *ast.File, fields *ast.FieldList) (analysis.TextEdit, bool) {
	if fields == nil || len(fields.List) == 0 || hasLineComment(f, fields) {
		return analysis.TextEdit{}, false
	}

	funcType := &ast.FuncType{Func: token.NoPos, TypeParams: fields}

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, funcType); err != nil {
		return analysis.TextEdit{}, false
	}

	original := buf.Bytes()
	fieldListText := original[bytes.IndexByte(original, '[')+1 : bytes.LastIndexByte(original, ']')]

	line, ok := condense(fieldListText)
	if !ok {
		return analysis.TextEdit{}, false
	}

	return analysis.TextEdit{Pos: fields.Pos() + 1, End: fields.End() - 1, NewText: bytes.Trim(line, " \n\r\t,")}, true
}

//nolint:cyclop,gocognit
func condense(b []byte) ([]byte, bool) {
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
