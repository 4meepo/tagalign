package tagalign

import (
	"path/filepath"
	"testing"

	"github.com/fatih/structtag"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	// unsort example
	a := NewAnalyzer()
	unsort, err := filepath.Abs("testdata/unsort")
	assert.NoError(t, err)
	analysistest.Run(t, unsort, a)

}
func TestAnalyzerWithOrder(t *testing.T) {
	// sort with fixed order
	a := NewAnalyzer(WithAutoSort("json", "yaml", "xml"))
	sort, err := filepath.Abs("testdata/sort")
	assert.NoError(t, err)
	analysistest.Run(t, sort, a)
}
func TestSprintf(t *testing.T) {
	format := alignFormat(20)
	assert.Equal(t, "%-20s", format)
}

func Test_sortByFixedOrder(t *testing.T) {
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
