package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
)

var Analyzer = &analysis.Analyzer{
	Name:     "goprintffuncname",
	Doc:      "Checks that printf-like functions are named with `f` at the end.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	// pass.ResultOf[inspect.Analyzer] will be set if we've added inspect.Analyzer to Requires.
	myInspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{ // filter needed nodes: visit only them
		(*ast.FuncDecl)(nil),
	}

	myInspector.Preorder(nodeFilter, func(node ast.Node) {
		funcDecl := node.(*ast.FuncDecl)

		params := funcDecl.Type.Params.List
		if len(params) != 2 { // [0] must be format (string), [1] must be args (...interface{})
			return
		}

		firstParamType, ok := params[0].Type.(*ast.Ident)
		if !ok { // first param type isn't identificator so it can't be of type "string"
			return
		}

		if firstParamType.Name != "string" { // first param (format) type is not string
			return
		}

		secondParamType, ok := params[1].Type.(*ast.Ellipsis)
		if !ok { // args are not ellipsis (...args)
			return
		}

		elementType, ok := secondParamType.Elt.(*ast.InterfaceType)
		if !ok { // args are not of interface type, but we need interface{}
			return
		}

		if elementType.Methods != nil && len(elementType.Methods.List) != 0 {
			return // has >= 1 method in interface, but we need an empty interface "interface{}"
		}

		if strings.HasSuffix(funcDecl.Name.Name, "f") {
			return
		}

		pass.Reportf(node.Pos(), "printf-like formatting function '%s' should be named '%sf'",
			funcDecl.Name.Name, funcDecl.Name.Name)
	})

	return nil, nil
}
