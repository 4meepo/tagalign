package tagalign

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
	"github.com/golangci/golangci-lint/pkg/golinters/goanalysis"
	"github.com/golangci/golangci-lint/pkg/lint/linter"

	"golang.org/x/tools/go/analysis"
)

func NewAnalyzerWithIssuesReporter() (*analysis.Analyzer, func(*linter.Context) []goanalysis.Issue) {
	var issues []goanalysis.Issue
	return &analysis.Analyzer{
			Name: "tagalign",
			Doc:  "check that struct tags are aligned",
			Run: func(p *analysis.Pass) (interface{}, error) {
				var err error
				issues, err = run(p)
				return nil, err
			},
		},
		func(ctx *linter.Context) []goanalysis.Issue {
			return issues
		}
}

func run(pass *analysis.Pass) ([]goanalysis.Issue, error) {
	for _, f := range pass.Files {
		var groups []group
		ast.Inspect(f, func(n ast.Node) bool {
			return findGroups(pass, n, &groups)
		})
		processGroups(&groups)
	}

	return nil, nil
}

// =======================

func findGroups(pass *analysis.Pass, n ast.Node, groups *[]group) bool {
	v, ok := n.(*ast.StructType)
	if !ok || len(v.Fields.List) == 0 {
		// no need to check non-struct or struct with 0 fields
		return true
	}

	findGroupInStruct(pass.Fset, v, groups)

	return true
}

type group struct {
	maxTagNum int
	lines     []*line
}
type line struct {
	field     *ast.Field
	tags      []string
	lens      []int
	spaceLens []int
	result    string
}

func findGroupInStruct(fset *token.FileSet, _struct *ast.StructType, groups *[]group, inline ...bool) {
	lastLineNum := fset.Position(_struct.Fields.List[0].Pos()).Line
	grp := group{}
	fieldsNum := len(_struct.Fields.List)
	for idx, field := range _struct.Fields.List {
		if field.Tag == nil {
			continue
		}

		tag, err := strconv.Unquote(field.Tag.Value)
		if err != nil {
			continue
		}

		tag = strings.TrimSpace(tag)

		tags, err := structtag.Parse(tag)
		if err != nil {
			continue
		}

		// in case the field is a struct type.
		if _, ok := field.Type.(*ast.StructType); ok {
			if idx+1 < fieldsNum {
				lastLineNum = fset.Position(_struct.Fields.List[idx+1].Pos()).Line // todo
			}

			*groups = append(*groups, grp)
			grp = group{}
			continue
		}

		if grp.maxTagNum < tags.Len() {
			grp.maxTagNum = tags.Len()
		}

		ln := &line{
			field: field,
		}

		lens := make([]int, 0, tags.Len())
		for _, key := range tags.Keys() {
			t, _ := tags.Get(key)
			lens = append(lens, length(t.String()))
			ln.tags = append(ln.tags, t.String())
		}

		ln.lens = lens

		lineNum := fset.Position(field.Pos()).Line
		if lineNum-lastLineNum >= 2 {
			*groups = append(*groups, grp)
			grp = group{
				maxTagNum: tags.Len(),
			}
		}

		lastLineNum = lineNum

		grp.lines = append(grp.lines, ln)
	}

	if len(grp.lines) > 0 {
		*groups = append(*groups, grp)
	}
}

func processGroups(groups *[]group) {
	for _, grp := range *groups {
		if len(grp.lines) <= 1 {
			continue
		}

		for i := 0; i < grp.maxTagNum; i++ {
			max := process0(grp.lines, i)
			updateResult(grp.lines, max, i)
		}

		for _, line := range grp.lines {
			line.result = "`" + line.result + "`"
		}
	}
}

func process0(lines []*line, idx int) int {
	max := 0
	for _, line := range lines {
		if len(line.lens) > idx {
			if line.lens[idx] > max {
				max = line.lens[idx]
			}
		}
	}

	return max
}

func updateResult(lines []*line, max, idx int) {
	for _, line := range lines {
		if len(line.tags) > idx {
			if l := len(line.lens); l > idx && idx < l-1 {
				line.result += line.tags[idx] + strings.Repeat(" ", max-line.lens[idx]+1)
			} else {
				line.result += line.tags[idx]
			}
		}
	}
}

func length(s string) int {
	return len([]rune(s))
}
