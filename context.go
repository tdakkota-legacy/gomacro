package macro

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/tdakkota/gomacro/pragma"
	"golang.org/x/tools/go/ast/astutil"
)

// ErrStop is special error to interrupt AST traversal.
var ErrStop = errors.New("macro exit")

// Context is a macro context.
type Context struct {
	*astutil.Cursor
	Pre     bool
	Delayed Delayed
	Report  func(Report)

	// Current file.
	File *ast.File
	// Current token set.
	FileSet *token.FileSet

	// Current package info.
	Package *types.Package
	// Type checker info.
	TypesInfo  *types.Info
	TypesSizes types.Sizes

	// Parsed magic comments.
	Pragmas pragma.Pragmas
}

// AddDecls adds declarations to current file.
func (c Context) AddDecls(decls ...ast.Decl) {
	c.File.Decls = append(c.File.Decls, decls...)
}

// Report represents macro error report.
type Report struct {
	Pos     token.Pos
	Message string
}

// Reportf reports macro error.
func (c Context) Reportf(pos token.Pos, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	c.Report(Report{Pos: pos, Message: msg})
}

func importName(s *ast.ImportSpec) string {
	n := s.Name
	if n == nil {
		return ""
	}
	return n.Name
}

func importEqual(a, b *ast.ImportSpec) bool {
	return a.Path.Value == b.Path.Value && importName(a) == importName(b)
}

// AddImports adds new imports to file.
// If import already exists AddImports does nothing.
func (c Context) AddImports(importSpec ...*ast.ImportSpec) {
	for _, spec := range importSpec {
		contains := false
		for _, imprt := range c.File.Imports {
			if importEqual(spec, imprt) {
				contains = true
				break
			}
		}

		if !contains {
			c.File.Imports = append(c.File.Imports, spec)
		}
	}
}
