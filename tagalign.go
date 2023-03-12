package tagalign

import "golang.org/x/tools/go/analysis"

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "tagalign",

		Doc: "check that struct tags are aligned",
		Run: run,
	}
}

func run(pass *analysis.Pass) (any, error) {
	panic("unimplement..")
}
