package tagalign

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
	"golang.org/x/tools/go/analysis"
)

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "tagalign",
		Doc:  "check that struct tags are aligned",
		Run:  run,
	}
}

func run(pass *analysis.Pass) (any, error) {
	var groups []group
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			return checkStruct(pass, n, &groups)
		})
		process(&groups)
	}
	// goanalysis.NewIssue

	return nil, nil
}

func checkStruct(pass *analysis.Pass, n ast.Node, groups *[]group) bool {
	v, ok := n.(*ast.StructType)
	if !ok || len(v.Fields.List) == 0 {
		return true
	}

	preProcessStruct(pass.Fset, v, groups)

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

func preProcessStruct(fset *token.FileSet, st *ast.StructType, groups *[]group, inline ...bool) {
	lastLineNum := fset.Position(st.Fields.List[0].Pos()).Line
	grp := group{}
	l := len(st.Fields.List)
	for idx, field := range st.Fields.List {
		if field.Tag == nil {
			continue
		}

		tag, err := strconv.Unquote(field.Tag.Value)
		if err != nil {
			continue
		}

		tag = strings.TrimLeft(tag, " ")
		tag = strings.TrimRight(tag, " ")

		tags, err := structtag.Parse(tag)
		if err != nil {
			continue
		}

		if _, ok := field.Type.(*ast.StructType); ok {
			if idx+1 < l {
				lastLineNum = fset.Position(st.Fields.List[idx+1].Pos()).Line
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
			lastLineNum = lineNum
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

func process(groups *[]group) {
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
