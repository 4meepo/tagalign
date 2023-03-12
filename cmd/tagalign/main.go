package main

import (
	"github.com/4meepo/tagalign"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(tagalign.NewAnalyzer())
}
