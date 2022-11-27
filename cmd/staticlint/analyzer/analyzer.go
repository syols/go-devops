package analyzer

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const Doc = `own parser that forbids using a direct call to os.Exit in the main function of the main package`

var Analyzer = &analysis.Analyzer{
	Name: "own_analyzer",
	Doc:  Doc,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() == "main" {
		for _, file := range pass.Files {
			ast.Inspect(file, inspectFunc(pass))
		}
	}
	return nil, nil
}

func inspectFunc(pass *analysis.Pass) func(node ast.Node) bool {
	return func(node ast.Node) bool {
		funcType, isFunc := node.(*ast.FuncType)
		if !isFunc {
			return false
		}

		info, isOk := pass.TypesInfo.Types[funcType]
		if !isOk {
			return false
		}

		if info.Type.String() == "os.Exit" {
			fmt.Printf("os.Exit in main")
			return false
		}
		return true
	}
}
