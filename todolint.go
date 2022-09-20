package todolint

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

const Doc = `Requires TODO comments to be in the form of "TODO(author) ...`

func NewAnalyzer() *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name: "todolint",
		Doc:  Doc,
		Run:  run,
	}
	return a
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, f := range pass.Files {
		for _, commentGroup := range f.Comments {
			for _, comment := range commentGroup.List {
				checkComment(pass, comment)
			}
		}
	}
	return nil, nil
}

type AnalysisPass interface {
	Reportf(pos token.Pos, format string, args ...interface{})
}

func checkComment(pass AnalysisPass, comment *ast.Comment) {
	const (
		groupLeading  = 1
		groupTodoText = 2
		groupTrailing = 3

		commentLeading = 2
	)

	c := comment.Text

	var clines []string
	if c[1] == '*' {
		/*-style comment, remove trailing comment markers */
		c = c[commentLeading : len(c)-commentLeading]
		clines = strings.Split(c, "\n")
	} else {
		//-style comment, no line breaks
		clines = []string{c[commentLeading:]}
	}

	start := commentLeading // skip leading "//" or "/*"
	var prevCommentLen int
	for _, cl := range clines {
		start += prevCommentLen
		prevCommentLen = len(cl) + 1 // 1 for \n

		// skip whitespaces
		var skip int
		for ; skip < len(cl); skip++ {
			if !isWhitespace(cl[skip]) {
				break
			}
		}

		r := matchTodoComment(cl[skip:])
		if r == nil {
			continue
		}

		leading := r.Group(groupLeading)
		todo := r.Group(groupTodoText)
		expectTodo := strings.ToUpper(todo)
		trailing := r.Group(groupTrailing)

		isLeadingOk := isEmptyOrWhitespaceLeading(leading)
		if isLeadingOk && expectTodo == todo && trailing == "(" {
			continue
		}

		pos := comment.Pos() + token.Pos(start+skip)
		if isLeadingOk {
			pos += token.Pos(r.GroupPos(groupTodoText))
		} else {
			pos += token.Pos(r.GroupPos(groupLeading))
		}
		pass.Reportf(pos, "TODO comment should be in the form %s(author)", expectTodo)
	}
}

var todoRE = regexp.MustCompile(`(\W)?((?i:TODO|FIXME))(\()?`)

func matchTodoComment(s string) *matchResult {
	if s == "" {
		return nil
	}

	r := todoRE.FindStringSubmatchIndex(s)
	if len(r) == 0 {
		return nil
	}

	return &matchResult{s, r}
}

type matchResult struct {
	s       string
	indices []int
}

func (m *matchResult) Group(n int) string {
	start, end := m.indices[2*n], m.indices[2*n+1]
	if start == -1 || end == -1 {
		return ""
	}
	return m.s[start:end]
}

func (m *matchResult) GroupPos(n int) int {
	return m.indices[2*n]
}

func isWhitespace(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

func isEmptyOrWhitespaceLeading(s string) bool {
	if s == "" {
		return true
	}
	return isWhitespace(s[0])
}
