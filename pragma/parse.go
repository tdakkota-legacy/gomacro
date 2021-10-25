package pragma

import (
	"go/ast"
	"strings"
)

const linterPragmaPrefix = "//procm:"

// Parser is a Go pragma parser.
type Parser struct {
	pragmaPrefix string
}

// NewParser creates new Parser.
func NewParser(prefix string) Parser {
	return Parser{pragmaPrefix: prefix}
}

func (p Parser) isPragmaComment(s string) bool {
	return strings.HasPrefix(s, p.pragmaPrefix)
}

// Comment represents Go comment.
type Comment struct {
	Value string
}

// ParseCommentGroup parses comments from given comment group.
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

func parsePragma(pair string) (k, v string, ok bool) {
	if pair == "" {
		return "", "", false
	}

	if !strings.Contains(pair, "=") {
		return pair, "", true
	}

	split := strings.SplitN(pair, "=", 2)
	return split[0], split[1], true
}

// ParsePragmas parses pragmas from given comments.
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

// ParsePragmas parses pragmas from given comment group.
func ParsePragmas(doc *ast.CommentGroup) Pragmas {
	parser := NewParser(linterPragmaPrefix)

	comments := parser.ParseCommentGroup(doc)
	return parser.ParsePragmas(comments)
}
