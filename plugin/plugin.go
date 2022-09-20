package main

import (
	"golang.org/x/tools/go/analysis"

	"github.com/timonwong/todolint"
)

// AnalyzerPlugin provides analyzers as a plugin.
// It follows golangci-lint style plugin.
var AnalyzerPlugin analyzerPlugin

type analyzerPlugin struct{}

func (analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		todolint.NewAnalyzer(),
	}
}
