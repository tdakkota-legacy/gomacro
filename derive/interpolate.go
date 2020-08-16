package derive

import (
	"go/ast"
	"go/parser"
	"strings"
)

type interpolator struct {
	replacer *strings.Replacer
}

func newInterpolator(a ...string) interpolator {
	return interpolator{replacer: strings.NewReplacer(a...)}
}

func (i interpolator) Interpolate(s string) string {
	return i.replacer.Replace(s)
}

func (i interpolator) Expr(s string) (ast.Expr, error) {
	return parser.ParseExpr(i.Interpolate(s))
}
