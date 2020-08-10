package macro

import (
	"go/ast"
	"go/printer"
	"go/token"
	"io"
)

type pp struct {
	config *printer.Config
}

func (p pp) PrintFile(w io.Writer, fset *token.FileSet, file *ast.File) error {
	ast.SortImports(fset, file)
	return p.config.Fprint(w, fset, file)
}

func DefaultPrinter() Printer {
	config := &printer.Config{
		Mode:     printer.TabIndent,
		Tabwidth: 8,
	}

	return pp{config}
}
