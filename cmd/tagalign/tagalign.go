package main

import (
	"flag"
	"os"
	"strings"

	"github.com/4meepo/tagalign"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	var align bool
	var sort bool
	var order string

	// Just for declare flags.
	flag.BoolVar(&align, "align", false, "Whether enable tags align. Default is true.")
	flag.BoolVar(&sort, "sort", false, "Whether enable tags sort. Default is false.")
	flag.StringVar(&order, "order", "", "Specify the order of tags, the other tags will be sorted by name.")

	// read from os.Args
	args := os.Args
	for i, arg := range args {
		if arg == "-align" {
			align = true
		}
		if arg == "-sort" {
			sort = true
		}
		if arg == "-order" {
			order = args[i+1]
		}
	}

	var options []tagalign.Option
	if sort {
		options = append(options, tagalign.WithAlign(align), tagalign.WithSort(strings.Split(order, ",")...))
	}

	singlechecker.Main(tagalign.NewAnalyzer(options...))
}
