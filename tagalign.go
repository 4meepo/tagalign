package tagalign

import (
	"fmt"
	"go/ast"
	"strconv"
	"strings"

	"github.com/fatih/structtag"

	"golang.org/x/tools/go/analysis"
)

func NewAnalyzerWithIssuesReporter() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "tagalign",
		Doc:  "check that struct tags are well aligned",
		Run: func(p *analysis.Pass) (interface{}, error) {
			var err error
			run(p)
			return nil, err
		},
	}
}

func run(pass *analysis.Pass) error {
	for _, f := range pass.Files {
		a := analyzer{}
		ast.Inspect(f, func(n ast.Node) bool {
			a.find(pass, n)
			return true
		})
		a.align(pass)
	}

	return nil
}

type analyzer struct {
	unalignedFieldsGroups [][]*ast.Field // fields in this group, must be consecutive in struct.
}

func (a *analyzer) find(pass *analysis.Pass, n ast.Node) {
	v, ok := n.(*ast.StructType)
	if !ok {
		return
	}

	fields := v.Fields.List
	if len(fields) == 0 {
		return
	}

	var fs []*ast.Field
	split := func() {
		if len(fs) > 1 {
			a.unalignedFieldsGroups = append(a.unalignedFieldsGroups, fs)
		}
		fs = nil
	}

	for i, field := range fields {
		if field.Tag == nil {
			// field without tags
			split()
			continue
		}

		if i > 0 {
			preLineNum := pass.Fset.Position(fields[i-1].Tag.Pos()).Line
			lineNum := pass.Fset.Position(field.Tag.Pos()).Line
			if lineNum-preLineNum > 1 {
				// fields with tags are not consecutive, including two case:
				// 1. splited by lines
				// 2. splited by a struct
				split()
				continue
			}
		}

		fs = append(fs, field)

	}

	split()
	return
}

func (a *analyzer) align(pass *analysis.Pass) {
	for _, fields := range a.unalignedFieldsGroups {
		// offsets := make([]int, len(fields))

		var maxTagNum int
		var tagsGroup [][]*structtag.Tag
		for _, field := range fields {
			// offsets[i] = pass.Fset.Position(field.Tag.Pos()).Column
			tag, err := strconv.Unquote(field.Tag.Value)
			if err != nil {
				break
			}

			tags, err := structtag.Parse(tag)
			if err != nil {
				break
			}

			maxTagNum = max(maxTagNum, tags.Len())

			tagsGroup = append(tagsGroup, tags.Tags())
		}

		// 记录每列 tag的最大长度
		tagMaxLens := make([]int, maxTagNum)

		for j := 0; j < maxTagNum; j++ {
			var maxLength int
			for i := 0; i < len(tagsGroup); i++ {
				if len(tagsGroup[i]) <= j {
					// in case of index out of range
					continue
				}
				maxLength = max(maxLength, len(tagsGroup[i][j].String()))
			}
			tagMaxLens[j] = maxLength
		}

		for i, field := range fields {
			tags := tagsGroup[i]

			newTagBuilder := strings.Builder{}
			for i, tag := range tags {
				format := alignFormat(tagMaxLens[i] + 1) // with an extra space
				newTagBuilder.WriteString(fmt.Sprintf(format, tag.String()))
			}

			newTagValue := fmt.Sprintf("`%s`", newTagBuilder.String())
			if field.Tag.Value == newTagValue {
				// nothing changed
				continue
			}

			pass.Report(analysis.Diagnostic{
				Pos:     field.Tag.Pos(),
				End:     field.Tag.End(),
				Message: "tag is not aligned, should be: " + newTagValue, SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "align tag",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     field.Tag.Pos(),
								End:     field.Tag.End(),
								NewText: []byte(newTagValue),
							},
						},
					},
				},
			})
		}
	}
}

func alignFormat(length int) string {
	return "%" + fmt.Sprintf(fmt.Sprintf("-%ds", length))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
