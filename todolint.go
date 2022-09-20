package todolint

import (
	"fmt"
	"go/token"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

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
	const (
		groupLeading  = 1
		groupTodoText = 2
		groupTrailing = 3
	)

	for _, f := range pass.Files {
		for _, commentGroup := range f.Comments {
			for _, comment := range commentGroup.List {
				commentText := comment.Text
				r := matchTodoComment(commentText)
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

				pos := comment.Pos()
				if isLeadingOk {
					pos += token.Pos(r.GroupPos(groupTodoText))
				} else {
					pos += token.Pos(r.GroupPos(groupLeading))
				}
				pass.Reportf(pos, "TODO comment should be in the form %s(author)", expectTodo)
			}
		}
		fmt.Printf("%v\n", f.Comments)
	}
	return nil, nil
}

var todoRE = regexp.MustCompile(`(?P<leading>\W)?(?P<todo>TODO|ToDo|todo|FIXME|FixMe|fixme|XXX|xxx)(?P<trailing>\()?`)

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
	pos := m.indices[2*n]
	return utf8.RuneCountInString(m.s[:pos])
}

func isEmptyOrWhitespaceLeading(s string) bool {
	if s == "" {
		return true
	}
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsSpace(r)
}
