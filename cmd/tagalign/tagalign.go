package main

import (
	"flag"
	"os"
	"strconv"
	"strings"

	"github.com/4meepo/tagalign"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	var noalign bool
	var sort bool
	var order string
	var strict bool
	var stopAlignThreshold int

	// just for declaration.
	flag.BoolVar(&noalign, "noalign", false, "Whether disable tags align. Default is false.")
	flag.BoolVar(&sort, "sort", false, "Whether enable tags sort. Default is false.")
	flag.BoolVar(&strict, "strict", false, "Whether enable strict style. Default is false. Note: strict must be used with align and sort together.")
	flag.StringVar(&order, "order", "", "Specify the order of tags, the other tags will be sorted by name.")
	flag.IntVar(&stopAlignThreshold, "threshold", 0, "Specifies the maximum allowable length difference between struct tags in the same column before alignment stops. When the difference between the longest and shortest tags exceeds this threshold, the alignment for subsequent tags in that column will be disabled, while the preceding tags will remain aligned. This helps maintain readability without creating large gaps. Because if the IDE wrap the line with too much gaps, it may be difficult to read. ( Default to 0, means disabled.)")

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
		if arg == "-threshold" {
			stopThreshold, err := strconv.Atoi(args[i+1])
			if err != nil {
				panic("invalid threshold value")
			}
			stopAlignThreshold = stopThreshold
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
	if stopAlignThreshold > 0 {
		options = append(options, tagalign.WithStopAlignThreshold(stopAlignThreshold))
	}

	singlechecker.Main(tagalign.NewAnalyzer(options...))
}
