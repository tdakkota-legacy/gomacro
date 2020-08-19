package rewriter

import (
	"bytes"
	"errors"
	"github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/internal/loader"
	"go/ast"
	"io"
	"os"
	"path/filepath"

	"github.com/tdakkota/gomacro/pragma"
)

type ReWriter struct {
	path, output string
	macros       macro.Macros
	printer      Printer
}

func NewReWriter(path, output string, macros macro.Macros, printer Printer) ReWriter {
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
	ctx, err := loader.LoadOne(r.path)
	if err != nil {
		return err
	}

	return r.runMacro(w, ctx)
}

func (r ReWriter) rewriteDir() error {
	return loader.LoadWalk(r.path, func(ctx macro.Context) error {
		file := ctx.FileSet.File(ctx.File.Pos()).Name()

		outputFile, err := prepareOutputFile(r.path, r.output, file)
		if err != nil {
			return err
		}

		return r.rewriteOneFile(outputFile, ctx)
	})
}

func (r ReWriter) rewriteFile() error {
	ctx, err := loader.LoadOne(r.path)
	if err != nil {
		return err
	}

	return r.rewriteOneFile(r.output, ctx)
}

func (r ReWriter) rewriteOneFile(output string, ctx macro.Context) error {
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

var errFailed = errors.New("failed to generate")

func (r ReWriter) runMacro(w io.Writer, context macro.Context) error {
	var imports *ast.GenDecl

	globalPragmas := pragma.ParsePragmas(context.File.Doc)
	globalMacros := r.macros.Get(globalPragmas.Use()...)

	macroRunner := NewRunner(context.FileSet)

	rewrites := len(globalMacros)
	for _, decl := range context.File.Decls {
		pragmas := pragma.ParsePragmas(loader.LoadComments(decl, &imports))

		localMacros := r.macros.Get(pragmas.Use()...)
		rewrites += len(localMacros)
		for _, handler := range localMacros {
			context.Pragmas = pragmas

			macroRunner.Run(handler, context, decl)
			if macroRunner.failed {
				return errFailed
			}
		}

		for name, handler := range globalMacros {
			if _, ok := localMacros[name]; ok {
				continue
			}
			context.Pragmas = pragmas

			macroRunner.Run(handler, context, decl)
			if macroRunner.failed {
				return errFailed
			}
		}
	}

	if len(context.File.Imports) > 0 {
		fixImports(imports, context.File)
	}

	return r.printer.PrintFile(w, context.FileSet, context.File)
}
