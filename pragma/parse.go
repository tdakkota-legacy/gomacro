package pragma

import (
	"go/ast"
	"strings"
)

const linterPragmaPrefix = "//procm:"

type Parser struct {
	pragmaPrefix string
}

func NewParser(prefix string) Parser {
	return Parser{pragmaPrefix: prefix}
}

func (p Parser) isPragmaComment(s string) bool {
	return strings.HasPrefix(s, p.pragmaPrefix)
}

type Comment struct {
	Value string
}

func (p Parser) ParseCommentGroup(doc *ast.CommentGroup) []Comment {
	if doc == nil {
		return nil
	}

	var result []Comment
	for _, comment := range doc.List {
		if p.isPragmaComment(comment.Text) {
			result = append(result, Comment{
				Value: strings.TrimSpace(strings.TrimPrefix(comment.Text, linterPragmaPrefix)),
			})
		}
	}

	return result
}

func parsePragma(pair string) (string, string, bool) {
	if pair == "" {
		return "", "", false
	}

	if !strings.Contains(pair, "=") {
		return pair, "", true
	}

	split := strings.SplitN(pair, "=", 2)
	return split[0], split[1], true
}

func (p Parser) ParsePragmas(comments []Comment) Pragmas {
	result := Pragmas{}

	for _, pragma := range comments {
		pairs := strings.SplitN(pragma.Value, ",", 1)

		for _, pair := range pairs {
			if k, v, ok := parsePragma(pair); ok {
				result[k] = v
			}
		}
	}

	return result
}

func ParsePragmas(doc *ast.CommentGroup) Pragmas {
	parser := NewParser(linterPragmaPrefix)

	comments := parser.ParseCommentGroup(doc)
	return parser.ParsePragmas(comments)
}
