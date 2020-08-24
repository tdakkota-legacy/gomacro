package rewriter

import (
	"errors"
	"fmt"
	"go/token"
	"os"
)

var errFailed = errors.New("failed to generate")

func NewErrorCallback(fset *token.FileSet) func(pos token.Pos, err error) {
	return func(pos token.Pos, err error) {
		_, _ = fmt.Fprintln(os.Stderr, fset.Position(pos).String(), err.Error())
	}
}

type errorPrinter struct {
	Error func(pos token.Pos, err error)
}

func (r errorPrinter) err(pos token.Pos, err error) {
	r.Error(pos, err)
}

func (r errorPrinter) Report(pos token.Pos, args ...interface{}) {
	r.err(pos, fmt.Errorf(fmt.Sprint(args...)))
}

func (r errorPrinter) Reportf(pos token.Pos, f string, args ...interface{}) {
	r.err(pos, fmt.Errorf(f, args...))
}
