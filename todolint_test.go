package todolint

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

type dummyTestingErrorf struct {
	*testing.T
}

func (t dummyTestingErrorf) Errorf(format string, args ...interface{}) {}

func TestLinter(t *testing.T) {
	testCases := []struct {
		name      string
		flags     []string
		wantError string
	}{
		{
			name: "all",
		},
		{
			name:      "error-empty-keywords",
			flags:     []string{"-keywords="},
			wantError: "at least one keyword must be specified",
		},
		{
			name:      "error-invalid-keyword-1",
			flags:     []string{"-keywords=,abc"},
			wantError: "invalid keyword \"\"",
		},
		{
			name:      "error-invalid-keyword-2",
			flags:     []string{"-keywords=abc,⏰"},
			wantError: "invalid keyword \"⏰\"",
		},
	}

	testdata := analysistest.TestData()
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			tl := NewAnalyzer()
			err := tl.Flags.Parse(tc.flags)
			require.NoError(t, err)

			var result []*analysistest.Result
			if tc.wantError != "" {
				result = analysistest.Run(&dummyTestingErrorf{t}, testdata, tl, "a")
			} else {
				result = analysistest.Run(t, testdata, tl, "a")
			}
			require.Len(t, result, 1)

			if tc.wantError != "" {
				assert.Error(t, result[0].Err)
				assert.ErrorContains(t, result[0].Err, tc.wantError)
			}
		})
	}
}

func getLineAtPos(s string, n int) string {
	for i := n; i < len(s); i++ {
		if s[i] == '\n' {
			return s[n:i]
		}
	}
	return s[n:]
}

func getLinesFromIssues(s string, issues []analysis.Diagnostic) []string {
	var lines []string
	for _, issue := range issues {
		lines = append(lines, getLineAtPos(s, int(issue.Pos)))
	}
	return lines
}

func TestTodoLint_CheckComment(t *testing.T) {
	testCases := []struct {
		name    string
		comment string
		want    []string
	}{
		{
			name:    "simple comment",
			comment: "// This is a simple comment without todos",
		},
		{
			name: "block comments mixed todos",
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
			t.Parallel()

			tl := newTodoLint()
			err := tl.processConfig()
			require.NoError(t, err)

			issues := tl.checkComment(&ast.Comment{Text: tc.comment})
			result := getLinesFromIssues(tc.comment, issues)
			assert.Equal(t, tc.want, result)
		})
	}
}
