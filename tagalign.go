package tagalign

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "tagalign",
		Doc:  "check that struct tags are aligned",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		ast.Inspect(f, checkStruct)
	}

	return nil, nil
}

func checkStruct(n ast.Node) bool {
	v, ok := n.(*ast.StructType)
	if !ok {
		return true
	}

	for _, field := range v.Fields.List {
		tag := field.Tag
		fmt.Println(tag.Value)
	}

	return true
}
