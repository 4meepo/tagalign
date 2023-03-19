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

type Mode int

const (
	StandaloneMode Mode = iota
	GolangciLintMode
)

func NewAnalyzer(mode Mode) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "tagalign",
		Doc:  "check that struct tags are well aligned",
		Run: func(p *analysis.Pass) (any, error) {
			RunTagAlign(p, mode)
			return nil, nil
		},
	}
}

func RunTagAlign(pass *analysis.Pass, mode Mode) []Issue {
	var issues []Issue
	for _, f := range pass.Files {
		a := Helper{mode: mode}
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
	mode         Mode
	fieldsGroups [][]*ast.Field // fields in this group, must be consecutive in struct.
	issues       []Issue
}

// Issue is used to integrate with golangci-lint's inline auto fix.
type Issue struct {
	Pos       token.Position
	Message   string
	InlineFix InlineFix
}
type InlineFix struct {
	StartCol  int //zero-based
	Length    int
	NewString string
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
			w.fieldsGroups = append(w.fieldsGroups, fs)
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
			if fields[i-1].Tag == nil {
				// if previous filed do not have a tag
				fs = append(fs, field)
				continue
			}
			preLineNum := pass.Fset.Position(fields[i-1].Tag.Pos()).Line
			lineNum := pass.Fset.Position(field.Tag.Pos()).Line
			if lineNum-preLineNum > 1 {
				// fields with tags are not consecutive, including two case:
				// 1. splited by lines
				// 2. splited by a struct
				split()

				// check if the field is a struct
				if _, ok := field.Type.(*ast.StructType); ok {
					continue
				}
			}
		}

		fs = append(fs, field)

	}

	split()
	return
}

func (w *Helper) align(pass *analysis.Pass) {
	for _, fields := range w.fieldsGroups {
		offsets := make([]int, len(fields))

		var maxTagNum int
		var tagsGroup [][]*structtag.Tag
		for i, field := range fields {
			offsets[i] = pass.Fset.Position(field.Tag.Pos()).Column
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

		// record the max length of each column tag
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

			unquoteTag := strings.TrimSpace(newTagBuilder.String())
			newTagValue := fmt.Sprintf("`%s`", unquoteTag)
			if field.Tag.Value == newTagValue {
				// nothing changed
				continue
			}

			msg := "tag is not aligned, should be: " + unquoteTag

			if w.mode == GolangciLintMode {
				iss := Issue{
					Pos:     pass.Fset.Position(field.Tag.Pos()),
					Message: msg,
					InlineFix: InlineFix{
						StartCol:  offsets[i] - 1,
						Length:    len(field.Tag.Value),
						NewString: newTagValue,
					},
				}
				w.issues = append(w.issues, iss)
			}

			if w.mode == StandaloneMode {
				pass.Report(analysis.Diagnostic{
					Pos:     field.Tag.Pos(),
					End:     field.Tag.End(),
					Message: msg,
					SuggestedFixes: []analysis.SuggestedFix{
						{
							Message: msg,
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
	if w.mode != GolangciLintMode {
		panic("Issues() should only be called in golangci-lint mode")
	}
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
