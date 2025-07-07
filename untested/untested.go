// Package untested checks that exported functions and methods have tests.
package untested

import (
	"go/ast"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"
)

var (
	internalFlag  = false
	generatedFlag = false
)

// NewAnalyzer returns an analyzer that reports exported functions and methods
// without corresponding tests.
func NewAnalyzer() *analysis.Analyzer {
	analyzer := &analysis.Analyzer{
		Name: "untested",
		Doc:  "check that exported functions and methods have tests",
		Run:  run,
	}

	analyzer.Flags.BoolVar(&internalFlag, "internal", false, "check functions in internal packages")
	analyzer.Flags.BoolVar(&generatedFlag, "generated", false, "check functions in generated files")

	return analyzer
}

// run is the main analyzer function that finds exported functions without tests.
// It builds a call graph from test files and checks which exported functions
// are not referenced directly or transitively from any test.
func run(pass *analysis.Pass) (any, error) {
	if !internalFlag && isInternalPackage(pass.Pkg.Path()) {
		return nil, nil
	}

	var exportedFunctions []*ast.FuncDecl

	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename

		if strings.HasSuffix(filename, "_test.go") || !generatedFlag && ast.IsGenerated(file) {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			if funcDecl, ok := n.(*ast.FuncDecl); ok {
				if funcDecl.Name.IsExported() {
					exportedFunctions = append(exportedFunctions, funcDecl)
				}
			}
			return true
		})
	}

	// If no exported functions, nothing to check
	if len(exportedFunctions) == 0 {
		return nil, nil
	}

	// Load packages with tests to find test references
	cfg := &packages.Config{
		Mode:  packages.LoadSyntax,
		Dir:   filepath.Dir(pass.Fset.Position(pass.Files[0].Pos()).Filename),
		Tests: true,
	}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return nil, err
	}

	testReferences := make(map[string]bool)

	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			filename := pkg.Fset.Position(file.Pos()).Filename
			if strings.HasSuffix(filename, "_test.go") {
				collectTestReferences(file, pkg, pass.Pkg, testReferences)
			}
		}
	}

	// Check each exported function for tests
	for _, funcDecl := range exportedFunctions {
		if key := getFuncDeclName(funcDecl); !testReferences[key] {
			pass.ReportRangef(funcDecl, "exported %s %q has no test", getFuncType(funcDecl), key)
		}
	}

	return nil, nil
}

// collectTestReferences builds a call graph from the given test file and all
// non-test files in the package, then propagates test coverage transitively.
// It marks exported functions as tested if they are called directly from tests
// or indirectly through helper functions.
func collectTestReferences(
	file *ast.File,
	pkg *packages.Package,
	targetPkg *types.Package,
	testReferences map[string]bool,
) {
	// Build a simple call graph using AST analysis
	callGraph := make(map[string][]string)
	ast.Inspect(file, func(n ast.Node) bool {
		if node, ok := n.(*ast.FuncDecl); ok && node.Name != nil {
			funcName := node.Name.Name
			collectCallsFromFunction(node, pkg, targetPkg, callGraph, testReferences, funcName)
		}
		return true
	})

	// Also collect calls from target package functions (for transitive calls)
	// Need to check all files in the target package, not just the test file
	for _, targetFile := range pkg.Syntax {
		targetFilename := pkg.Fset.Position(targetFile.Pos()).Filename
		if !strings.HasSuffix(targetFilename, "_test.go") {
			ast.Inspect(targetFile, func(n ast.Node) bool {
				if funcDecl, ok := n.(*ast.FuncDecl); ok {
					if funcDecl.Name != nil {
						funcName := funcDecl.Name.Name
						collectCallsFromFunction(funcDecl, pkg, targetPkg, callGraph, testReferences, funcName)
					}
				}
				return true
			})
		}
	}

	propagateTestCoverage(callGraph, testReferences)
}

// propagateTestCoverage performs transitive closure on the call graph to mark
// all functions reachable from test functions as tested. It iteratively
// propagates test coverage: if a function is called by a test or an already
// tested function, it gets marked as tested too.
func propagateTestCoverage(callGraph map[string][]string, testReferences map[string]bool) {
	changed := true
	for changed {
		changed = false

		for caller, callees := range callGraph {
			if isTestFunction(caller) || testReferences[caller] {
				for _, callee := range callees {
					if !testReferences[callee] {
						testReferences[callee] = true
						changed = true
					}
				}
			}
		}
	}
}

// collectCallsFromFunction analyzes a function declaration to find all function
// calls within it, building the call graph and marking exported functions as tested.
func collectCallsFromFunction(
	funcDecl *ast.FuncDecl,
	pkg *packages.Package,
	targetPkg *types.Package,
	callGraph map[string][]string,
	testReferences map[string]bool,
	funcName string,
) {
	ast.Inspect(funcDecl, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		var ident *ast.Ident
		switch e := call.Fun.(type) {
		case *ast.Ident:
			ident = e
		case *ast.SelectorExpr:
			ident = e.Sel
		default:
			return true
		}

		obj := pkg.TypesInfo.ObjectOf(ident)
		if obj == nil {
			return true
		}

		fn, ok := obj.(*types.Func)
		if !ok {
			return true
		}

		switch {
		case fn.Pkg() != nil && fn.Pkg().Name() == targetPkg.Name():
			if fn.Exported() {
				key := getFuncTypeName(fn)
				testReferences[key] = true
				if funcName != "" {
					callGraph[funcName] = append(callGraph[funcName], key)
				}
			} else if funcName != "" {
				callGraph[funcName] = append(callGraph[funcName], fn.Name())
			}
		case fn.Pkg() == nil || fn.Pkg().Name() == pkg.Types.Name():
			if funcName != "" {
				callGraph[funcName] = append(callGraph[funcName], fn.Name())
			}
		}

		return true
	})
}

// isTestFunction determines if a function name represents a test, benchmark, or example function.
func isTestFunction(name string) bool {
	return strings.HasPrefix(name, "Test") || strings.HasPrefix(name, "Benchmark") || strings.HasPrefix(name, "Example")
}

// getFuncTypeName returns the qualified name of a function from types info,
// including the receiver type for methods (e.g., "Type.Method" or "Function").
func getFuncTypeName(fn *types.Func) string {
	sig := fn.Type().(*types.Signature)

	recv := sig.Recv()
	if recv == nil {
		return fn.Name()
	}

	recvType := recv.Type()
	if ptr, ok := recvType.(*types.Pointer); ok {
		recvType = ptr.Elem()
	}

	if named, ok := recvType.(*types.Named); ok {
		return named.Obj().Name() + "." + fn.Name()
	}

	return recvType.String() + "." + fn.Name()
}

// getFuncDeclName returns the qualified name of a function from AST declaration,
// including the receiver type for methods (e.g., "Type.Method" or "Function").
func getFuncDeclName(fn *ast.FuncDecl) string {
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		recv := fn.Recv.List[0].Type
		if star, ok := recv.(*ast.StarExpr); ok {
			recv = star.X
		}

		if ident, ok := recv.(*ast.Ident); ok {
			return ident.Name + "." + fn.Name.Name
		}
	}

	return fn.Name.Name
}

// getFuncType returns "method" for methods or "function" for functions,
// used for error message formatting.
func getFuncType(fn *ast.FuncDecl) string {
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		return "method"
	}

	return "function"
}

// isInternalPackage determines if a package path contains an "internal" component.
func isInternalPackage(pkgPath string) bool {
	for part := range strings.SplitSeq(pkgPath, "/") {
		if part == "internal" {
			return true
		}
	}

	return false
}
