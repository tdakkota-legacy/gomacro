package macro

import (
	"bytes"
	"go/ast"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/tdakkota/gomacro/pragma"
)

type Printer interface {
	PrintFile(w io.Writer, fset *token.FileSet, file *ast.File) error
}

type ReWriter struct {
	path, output string
	macros       Macros
	printer      Printer
}

func NewReWriter(path string, output string, macros Macros, printer Printer) ReWriter {
	return ReWriter{path: path, output: output, macros: macros, printer: printer}
}

func (r ReWriter) Rewrite() error {
	fi, err := os.Stat(r.path)
	if err != nil {
		return err
	}

	switch mode := fi.Mode(); {
	case mode.IsDir():
		return r.rewriteDir()
	case mode.IsRegular():
		return r.rewriteFile()
	}

	return nil
}

func (r ReWriter) RewriteTo(w io.Writer) error {
	ctx, err := loadOne(r.path)
	if err != nil {
		return err
	}

	return r.runMacro(w, ctx)
}

func loadDelayed(pkgs []*packages.Package) Delayed {
	delayed := Delayed{}
	for _, pkg := range pkgs {
		delayed.Add(pkg)
	}

	return delayed
}

func (r ReWriter) rewriteDir() error {
	loadPath := r.path
	if !strings.HasSuffix(r.path, "/...") {
		loadPath = loadPath + "/..."
	}

	pkgs, err := load(loadPath)
	if err != nil {
		return err
	}

	delayed := loadDelayed(pkgs)

	for _, pkg := range pkgs {
		ctx := createContext(delayed, pkg)

		for _, file := range pkg.Syntax {
			ctx.File = file

			outputFile, err := createDir(r.path, r.output, ctx.FileSet.File(file.Pos()).Name())
			if err != nil {
				return err
			}

			err = r.rewriteOneFile(outputFile, ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r ReWriter) rewriteFile() error {
	ctx, err := loadOne(r.path)
	if err != nil {
		return err
	}

	return r.rewriteOneFile(r.output, ctx)
}

func (r ReWriter) rewriteOneFile(output string, ctx Context) error {
	buf := new(bytes.Buffer)
	err := r.runMacro(buf, ctx)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(output), os.ModePerm)
	if err != nil {
		return err
	}

	if buf.Len() > 0 {
		var w io.Writer = os.Stdout
		if output != "" {
			w, err = os.Create(output)
			if err != nil {
				return err
			}
		}

		_, err = io.Copy(w, buf)
		return err
	}

	return nil
}

func loadComments(decl ast.Decl, imports **ast.GenDecl) (comments *ast.CommentGroup) {
	switch v := decl.(type) {
	case *ast.GenDecl:
		if v.Tok == token.IMPORT {
			*imports = v
		}
		comments = v.Doc
	case *ast.FuncDecl:
		comments = v.Doc
	}

	return
}

func (r ReWriter) runMacro(w io.Writer, context Context) error {
	var imports *ast.GenDecl

	globalPragmas := pragma.ParsePragmas(context.File.Doc)
	globalMacros := r.macros.Get(globalPragmas.Use()...)

	macroRunner := NewRunner(context.FileSet)

	rewrites := len(globalMacros)
	for _, decl := range context.File.Decls {
		pragmas := pragma.ParsePragmas(loadComments(decl, &imports))

		localMacros := r.macros.Get(pragmas.Use()...)
		rewrites += len(localMacros)
		for _, macro := range localMacros {
			context.Pragmas = pragmas
			macroRunner.Run(macro, context, decl)
		}

		for name, macro := range globalMacros {
			if _, ok := localMacros[name]; ok {
				continue
			}
			context.Pragmas = pragmas
			macroRunner.Run(macro, context, decl)
		}
	}

	if len(context.File.Imports) > 0 {
		fixImports(imports, context.File)
	}

	if rewrites > 0 {
		err := r.printer.PrintFile(w, context.FileSet, context.File)
		if err != nil {
			return err
		}
	}

	return nil
}
