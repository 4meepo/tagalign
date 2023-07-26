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
	var strict bool

	// just for declaration.
	flag.BoolVar(&noalign, "noalign", false, "Whether disable tags align. Default is false.")
	flag.BoolVar(&sort, "sort", false, "Whether enable tags sort. Default is false.")
	flag.BoolVar(&strict, "strict", false, "Whether enable strict style. Default is false. Note: strict must be used with align and sort together.")
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
		if arg == "-strict" {
			strict = true
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
	if strict {
		if noalign {
			// cannot use noalign and strict together.
			panic("cannot use `-noalign` and `-strict` together.")
		}
		if !sort {
			// cannot use strict without sort.
			panic("cannot use `-strict` without `-sort`.")
		}
		options = append(options, tagalign.WithStrictStyle())
	}

	singlechecker.Main(tagalign.NewAnalyzer(options...))
}
