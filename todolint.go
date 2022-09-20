package todolint

import (
	"fmt"
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
		fmt.Printf("%v\n", f.Comments)
	}
	return nil, nil
}

func checkComment(pass *analysis.Pass, comment *ast.Comment) {
	const (
		groupLeading  = 1
		groupTodoText = 2
		groupTrailing = 3

		commentLeading = 2
	)

	commentText := comment.Text
	skip := commentLeading // skip "//" or "/*"

	// skip spaces
	for ; skip < len(commentText); skip++ {
		if !isWhitespace(commentText[skip]) {
			break
		}
	}

	r := matchTodoComment(commentText[skip:])
	if r == nil {
		return
	}

	leading := r.Group(groupLeading)
	todo := r.Group(groupTodoText)
	expectTodo := strings.ToUpper(todo)
	trailing := r.Group(groupTrailing)

	isLeadingOk := isEmptyOrWhitespaceLeading(leading)
	if isLeadingOk && expectTodo == todo && trailing == "(" {
		return
	}

	pos := comment.Pos() + token.Pos(skip)
	if isLeadingOk {
		pos += token.Pos(r.GroupPos(groupTodoText))
	} else {
		pos += token.Pos(r.GroupPos(groupLeading))
	}
	pass.Reportf(pos, "TODO comment should be in the form %s(author)", expectTodo)
}

var todoRE = regexp.MustCompile(`(\W)?((?i:TODO|FIXME))(\()?`)

func matchTodoComment(s string) *matchResult {
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
