package tagalign

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/fatih/structtag"

	"golang.org/x/tools/go/analysis"
)

func NewAnalyzer(report bool) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "tagalign",
		Doc:  "check that struct tags are well aligned",
		Run: func(p *analysis.Pass) (any, error) {
			RunTagAlign(p, report)
			return nil, nil
		},
	}
}

func RunTagAlign(pass *analysis.Pass, report bool) []Issue {
	var issues []Issue
	for _, f := range pass.Files {
		a := Helper{report: report}
		ast.Inspect(f, func(n ast.Node) bool {
			a.find(pass, n)
			return true
		})
		a.align(pass)
		issues = append(issues, a.issues...)
	}
	return issues
}

type Helper struct {
	report                bool
	unalignedFieldsGroups [][]*ast.Field // fields in this group, must be consecutive in struct.
	issues                []Issue
}

type Issue struct {
	Pos               token.Position
	Message           string
	InlineReplacement string
}

func (w *Helper) find(pass *analysis.Pass, n ast.Node) {
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
			w.unalignedFieldsGroups = append(w.unalignedFieldsGroups, fs)
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

func (w *Helper) align(pass *analysis.Pass) {
	for _, fields := range w.unalignedFieldsGroups {
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

			// new new builder
			newTagBuilder := strings.Builder{}
			for i, tag := range tags {
				format := alignFormat(tagMaxLens[i] + 1) // with an extra space
				newTagBuilder.WriteString(fmt.Sprintf(format, tag.String()))
			}

			newTagValue := strings.TrimSpace(fmt.Sprintf("`%s`", newTagBuilder.String()))
			if field.Tag.Value == newTagValue {
				// nothing changed
				continue
			}

			msg := "tag is not aligned, should be: " + newTagValue
			// for integrate with golangci-lint
			iss := Issue{
				Pos:               pass.Fset.Position(field.Tag.Pos()),
				Message:           msg,
				InlineReplacement: newTagValue,
			}
			w.issues = append(w.issues, iss)

			if w.report {
				pass.Report(analysis.Diagnostic{
					Pos:     field.Tag.Pos(),
					End:     field.Tag.End(),
					Message: msg,
					SuggestedFixes: []analysis.SuggestedFix{
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
}

// Issues returns all issues found by the analyzer.
// It is used to integrate with golangci-lint.
func (w *Helper) Issues() []Issue {
	return w.issues
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
