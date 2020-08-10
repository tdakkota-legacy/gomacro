package macroctx

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/tdakkota/gomacro/pragma"
	"golang.org/x/tools/go/ast/astutil"
)

type Context struct {
	*astutil.Cursor
	Pre     bool
	Delayed Delayed
	Report  func(Report)

	// AST
	File    *ast.File
	FileSet *token.FileSet

	// Types
	Package    *types.Package
	TypesInfo  *types.Info
	TypesSizes types.Sizes

	Pragmas pragma.Pragmas
}

func (c Context) AddDecls(decls ...ast.Decl) {
	c.File.Decls = append(c.File.Decls, decls...)
}

type Report struct {
	Pos     token.Pos
	Message string
}

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

func (c Context) AddImports(importSpec ...*ast.ImportSpec) {
	for _, spec := range importSpec {
		contains := false
		for _, imprt := range c.File.Imports {
			if spec.Path.Value == imprt.Path.Value && importName(spec) == importName(imprt) {
				contains = true
				break
			}
		}

		if !contains {
			c.File.Imports = append(c.File.Imports, spec)
		}
	}
}
