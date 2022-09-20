package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/timonwong/todolint"
)

func main() {
	singlechecker.Main(todolint.NewAnalyzer())
}
