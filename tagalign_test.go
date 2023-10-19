package tagalign

import (
	"path/filepath"
	"testing"

	"github.com/fatih/structtag"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_alignOnly(t *testing.T) {
	// only align
	a := NewAnalyzer()
	unsort, err := filepath.Abs("testdata/align")
	assert.NoError(t, err)
	analysistest.Run(t, unsort, a)
}

func Test_sortOnly(t *testing.T) {
	a := NewAnalyzer(WithAlign(false), WithSort(nil...))
	sort, err := filepath.Abs("testdata/sort")
	assert.NoError(t, err)
	analysistest.Run(t, sort, a)
}

func Test_sortWithOrder(t *testing.T) {
	// test disable align but enable sort
	a := NewAnalyzer(WithAlign(false), WithSort("xml", "json", "yaml"))
	sort, err := filepath.Abs("testdata/sortorder")
	assert.NoError(t, err)
	analysistest.Run(t, sort, a)
}

func Test_alignAndSortWithOrder(t *testing.T) {
	// align and sort with fixed order
	a := NewAnalyzer(WithSort("json", "yaml", "xml"))
	sort, err := filepath.Abs("testdata/alignsortorder")
	assert.NoError(t, err)
	analysistest.Run(t, sort, a)
}

func TestSprintf(t *testing.T) {
	format := alignFormat(20)
	assert.Equal(t, "%-20s", format)
}

func Test_sortBy(t *testing.T) {
	tags, err := structtag.Parse(`zip:"foo" json:"foo,omitempty" yaml:"bar" binding:"required" xml:"baz" gorm:"column:foo"`)
	assert.NoError(t, err)

	sortBy([]string{"json", "yaml", "xml"}, tags)
	assert.Equal(t, "json", tags.Tags()[0].Key)
	assert.Equal(t, "yaml", tags.Tags()[1].Key)
	assert.Equal(t, "xml", tags.Tags()[2].Key)
	assert.Equal(t, "binding", tags.Tags()[3].Key)
	assert.Equal(t, "gorm", tags.Tags()[4].Key)
	assert.Equal(t, "zip", tags.Tags()[5].Key)
}

func Test_strictStyle(t *testing.T) {
	// align and sort with fixed order
	a := NewAnalyzer(WithSort("json", "yaml", "xml"), WithStrictStyle())
	sort, err := filepath.Abs("testdata/strict")
	assert.NoError(t, err)
	analysistest.Run(t, sort, a)
}

func Test_alignSingleField(t *testing.T) {
	// only align
	a := NewAnalyzer()
	unsort, err := filepath.Abs("testdata/single_field")
	assert.NoError(t, err)
	analysistest.Run(t, unsort, a)
}

func Test_badSyntaxTag(t *testing.T) {
	// only align
	a := NewAnalyzer()
	unsort, err := filepath.Abs("testdata/bad_syntax_tag")
	assert.NoError(t, err)
	analysistest.Run(t, unsort, a)
}
