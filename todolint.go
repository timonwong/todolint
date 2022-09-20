package todolint

import (
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/timonwong/todolint/internal/list"
)

const Doc = `Requires TODO comments to be in the form of "TODO(author) ...`

func NewAnalyzer() *analysis.Analyzer {
	t := newTodoLint()
	a := &analysis.Analyzer{
		Name: "todolint",
		Doc:  Doc,
		Run:  t.run,
	}
	t.bindFlags(&a.Flags)
	return a
}

type todolint struct {
	keywords list.String

	re *regexp.Regexp
}

func newTodoLint() *todolint {
	return &todolint{
		keywords: list.NewString("TODO", "FIXME"),
	}
}

func (t *todolint) bindFlags(fs *flag.FlagSet) {
	fs.Var(&t.keywords, "keywords", "comma-separated list of case-insensitive keywords to check for")
}

func (t *todolint) processConfig() error {
	if len(t.keywords) == 0 {
		return errors.New("at least one keyword must be specified")
	}

	keywords := make([]string, 0, len(t.keywords))
	for _, keyword := range t.keywords {
		if !isValidKeyword(keyword) {
			return fmt.Errorf("invalid keyword %q", keyword)
		}
		keywords = append(keywords, keyword)
	}

	expr := fmt.Sprintf(`(\W)?((?i:%s)\b)(\()?`, strings.Join(keywords, "|"))
	t.re = regexp.MustCompile(expr)
	return nil
}

func (t *todolint) run(pass *analysis.Pass) (interface{}, error) {
	err := t.processConfig()
	if err != nil {
		return nil, err
	}

	for _, f := range pass.Files {
		for _, commentGroup := range f.Comments {
			for _, comment := range commentGroup.List {
				issues := t.checkComment(comment)
				for _, issue := range issues {
					pass.Reportf(issue.Pos, issue.Message)
				}
			}
		}
	}
	return nil, nil
}

func (t *todolint) checkComment(comment *ast.Comment) (issues []analysis.Diagnostic) {
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

		r := t.matchTodoComment(cl[skip:])
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

		issues = append(issues, analysis.Diagnostic{
			Pos:      pos,
			Category: "comment",
			Message:  fmt.Sprintf("TODO comment should be in the form %s(author)", expectTodo),
		})
	}

	return issues
}

func (t *todolint) matchTodoComment(s string) *matchResult {
	if s == "" {
		return nil
	}

	r := t.re.FindStringSubmatchIndex(s)
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

// keyword can only contain letters and numbers
func isValidKeyword(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !isAlphaNum(c) {
			return false
		}
	}
	return true
}

func isAlphaNum(c int32) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9')
}

func isWhitespace(ch byte) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' }

func isEmptyOrWhitespaceLeading(s string) bool {
	if s == "" {
		return true
	}
	return isWhitespace(s[0])
}
