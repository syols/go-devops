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
			ast.Inspect(file, func(node ast.Node) bool {
				if funcType, ok := node.(*ast.FuncType); ok {
					if info, isOk := pass.TypesInfo.Types[funcType]; isOk {
						if info.Type.String() == "os.Exit" {
							fmt.Printf("os.Exit in main")
							return false
						}
					}
				}
				return true
			})
		}
	}

	return nil, nil
}
