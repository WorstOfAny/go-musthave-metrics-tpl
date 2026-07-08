// Анализатор, который проверяет, есть ли в коде вызов os.Exit
package osexitcheckanalyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for os.Exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	expr := func(x *ast.ExprStmt) {
		if call, ok := x.X.(*ast.CallExpr); ok {
			if isOsExit(pass, call) {
				pass.Reportf(x.Pos(), "expression is os.Exit")
			}
		}
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if x, ok := node.(*ast.ExprStmt); ok {
				expr(x)
			}

			return true
		})
	}
	return nil, nil
}

func isOsExit(pass *analysis.Pass, call *ast.CallExpr) bool {
	if sst, ok := call.Fun.(*ast.SelectorExpr); ok {
		if ident, ok := sst.X.(*ast.Ident); ok {
			if ident.Name == "os" {
				if sst.Sel.Name == "Exit" {
					return true
				}
			}
		}
	}

	return false
}
