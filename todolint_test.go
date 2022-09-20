package todolint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLinter(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, NewAnalyzer(), "a")
}
