package rewriter

import (
	"go/ast"
	"go/format"
	"go/token"
	"io"
)

// Printer is an abstraction for Go source printer.
type Printer interface {
	PrintFile(w io.Writer, fset *token.FileSet, file *ast.File) error
}

type pp struct{}

func (p pp) PrintFile(w io.Writer, fset *token.FileSet, file *ast.File) error {
	ast.SortImports(fset, file)

	return format.Node(w, fset, file)
}

// DefaultPrinter returns default Printer implementation
func DefaultPrinter() Printer {
	return pp{}
}
