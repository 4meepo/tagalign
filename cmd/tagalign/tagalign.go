package main

import (
	"flag"
	"os"
	"strings"

	"github.com/4meepo/tagalign"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	var autoSort bool
	var fixedOrder string

	// Just for declase
	flag.BoolVar(&autoSort, "auto-sort", false, "enable auto sort tags")
	flag.StringVar(&fixedOrder, "fixed-order", "", "specify the fixed order of tags, the other tags will be sorted by name")

	// read from os.Args
	args := os.Args
	for i, arg := range args {
		if arg == "-auto-sort" {
			autoSort = true
		}
		if arg == "-fixed-order" {
			fixedOrder = args[i+1]
		}
	}

	var options []tagalign.Option
	if autoSort {
		options = append(options, tagalign.WithAutoSort(strings.Split(fixedOrder, ",")...))
	}

	singlechecker.Main(tagalign.NewAnalyzer(options...))
}
