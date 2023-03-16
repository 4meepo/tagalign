package main

import (
	"github.com/4meepo/tagalign"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	a := tagalign.NewAnalyzerWithIssuesReporter()
	singlechecker.Main(a)
}
