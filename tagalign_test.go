package tagalign

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	a, _ := NewAnalyzerWithIssuesReporter()
	analysistest.Run(t, analysistest.TestData(), a)
}
