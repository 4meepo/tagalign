package main

import (
	"flag"
	"os"
	"strings"

	"github.com/4meepo/tagalign"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	var noalign bool
	var sort bool
	var order string

	// just for declaration.
	flag.BoolVar(&noalign, "noalign", false, "Whether disable tags align. Default is false.")
	flag.BoolVar(&sort, "sort", false, "Whether enable tags sort. Default is false.")
	flag.StringVar(&order, "order", "", "Specify the order of tags, the other tags will be sorted by name.")

	// read from os.Args
	args := os.Args
	for i, arg := range args {
		if arg == "-noalign" {
			noalign = true
		}
		if arg == "-sort" {
			sort = true
		}
		if arg == "-order" {
			order = args[i+1]
		}
	}

	var options []tagalign.Option
	if noalign {
		options = append(options, tagalign.WithAlign(false))
	}
	if sort {
		var orders []string
		if order != "" {
			orders = strings.Split(order, ",")
		}
		options = append(options, tagalign.WithSort(orders...))
	}

	singlechecker.Main(tagalign.NewAnalyzer(options...))
}
