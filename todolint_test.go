package todolint

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestLinter(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, NewAnalyzer(), "a")
}

type AnalysisPassStub struct {
	positions []token.Pos
}

func (a *AnalysisPassStub) Reportf(pos token.Pos, format string, args ...interface{}) {
	a.positions = append(a.positions, pos)
}

func getLineAtPos(s string, n int) string {
	for i := n; i < len(s); i++ {
		if s[i] == '\n' {
			return s[n:i]
		}
	}
	return s[n:]
}

func getLinesFromPositions(s string, positions []token.Pos) []string {
	var lines []string
	for _, pos := range positions {
		lines = append(lines, getLineAtPos(s, int(pos)))
	}
	return lines
}

func TestCheckComment(t *testing.T) {
	pass := &AnalysisPassStub{}

	testCases := []struct {
		name    string
		comment string
		want    []string
	}{
		{
			name: "block-multiple-todos",
			comment: `/*

    * TODO(author): This is OK
    * FIXME: Missing author
    * @TODO(author): Invalid leading "@"
    * FixMe(author): Not all-caps

    * FIXME(author): This is OK
*/`,
			want: []string{
				"FIXME: Missing author",
				"@TODO(author): Invalid leading \"@\"",
				"FixMe(author): Not all-caps",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			checkComment(pass, &ast.Comment{Text: tc.comment})
			result := getLinesFromPositions(tc.comment, pass.positions)
			if !cmp.Equal(tc.want, result) {
				t.Errorf("got: %v, want: %v, diff: %s", result, tc.want, cmp.Diff(tc.want, result))
			}
		})
	}
}
