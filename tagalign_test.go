package tagalign

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	a := NewAnalyzerWithIssuesReporter()
	analysistest.Run(t, analysistest.TestData(), a)
}

func TestSprintf(t *testing.T) {
	format := alignFormat(20)
	assert.Equal(t, "%-20s", format)
}
