package tagalign

import (
	"testing"

	"github.com/alfatraining/structtag"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testCases := []struct {
		desc string
		dir  string
		opts []Option
	}{
		{
			desc: "only align",
			dir:  "align_only",
		},
		{
			desc: "sort only",
			dir:  "sort_only",
			opts: []Option{WithAlign(false), WithSort(nil...)},
		},
		{
			desc: "sort with order",
			dir:  "sortorder",
			opts: []Option{WithAlign(false), WithSort("xml", "json", "yaml")},
		},
		{
			desc: "align and sort with fixed order",
			dir:  "alignsortorder",
			opts: []Option{WithSort("json", "yaml", "xml")},
		},
		{
			desc: "strict style",
			dir:  "strict",
			opts: []Option{WithSort("json", "yaml", "xml"), WithStrictStyle()},
		},
		{
			desc: "align single field",
			dir:  "single_field",
		},
		{
			desc: "bad syntax tag",
			dir:  "bad_syntax_tag",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			a := NewAnalyzer(test.opts...)

			analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), a, test.dir)
		})
	}
}

func TestAnalyzer_cgo(t *testing.T) {
	a := NewAnalyzer()

	analysistest.Run(t, analysistest.TestData(), a, "cgo")
}

func Test_alignFormat(t *testing.T) {
	format := alignFormat(20)
	assert.Equal(t, "%-20s", format)
}

func Test_sortTags(t *testing.T) {
	tags, err := structtag.Parse(`zip:"foo" json:"foo,omitempty" yaml:"bar" binding:"required" xml:"baz" gorm:"column:foo"`)
	assert.NoError(t, err)

	sortTags([]string{"json", "yaml", "xml"}, tags)
	assert.Equal(t, "json", tags.Tags()[0].Key)
	assert.Equal(t, "yaml", tags.Tags()[1].Key)
	assert.Equal(t, "xml", tags.Tags()[2].Key)
	assert.Equal(t, "binding", tags.Tags()[3].Key)
	assert.Equal(t, "gorm", tags.Tags()[4].Key)
	assert.Equal(t, "zip", tags.Tags()[5].Key)
}
