package rewriter

import (
	"errors"
	"go/ast"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/internal/loader"

	"github.com/tdakkota/gomacro/pragma"
)

type ReWriter struct {
	source, output string
	appendMode     bool
	macros         macro.Macros
	printer        Printer

	loaded loader.Loaded
}

func NewReWriter(source, output string, macros macro.Macros, printer Printer) ReWriter {
	return ReWriter{
		source:  source,
		output:  output,
		macros:  macros,
		printer: printer,
	}
}

func (r ReWriter) Source() string {
	return r.source
}

func (r ReWriter) Output() string {
	return r.output
}

func (r ReWriter) Rewrite() error {
	fi, err := os.Stat(r.source)
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
	ctx, err := loader.LoadOne(r.source)
	if err != nil {
		return err
	}

	return r.runMacro(w, ctx)
}

func (r ReWriter) rewriteDir() error {
	return loader.LoadWalk(r.source, func(l loader.Loaded, ctx macro.Context) error {
		r.loaded = l
		file := ctx.FileSet.File(ctx.File.Pos()).Name()

		outputFile, err := r.prepareOutputFile(file)
		if err != nil {
			return err
		}

		return r.rewriteOneFile(outputFile, ctx)
	})
}

func (r ReWriter) prepareOutputFile(file string) (string, error) {
	if r.appendMode {
		return prepareGenFile(file)
	}

	return prepareOutputFile(r.source, r.output, file)
}

func (r ReWriter) rewriteFile() error {
	ctx, err := loader.LoadOne(r.source)
	if err != nil {
		return err
	}

	return r.rewriteOneFile(r.output, ctx)
}

func (r ReWriter) rewriteOneFile(output string, ctx macro.Context) error {
	err := os.MkdirAll(filepath.Dir(output), os.ModePerm)
	if err != nil {
		return err
	}

	var w io.Writer = os.Stdout
	if output != "" {
		w, err = os.Create(output)
		if err != nil {
			return err
		}
	}

	return r.runMacro(w, ctx)
}

var errFailed = errors.New("failed to generate")

func copyDecls(a []ast.Decl) []ast.Decl {
	copyDecl := make([]ast.Decl, len(a))
	copy(copyDecl, a)
	return copyDecl
}

func (r ReWriter) runMacro(w io.Writer, context macro.Context) error {
	macroRunner := NewRunner(context.FileSet)
	globalPragmas := pragma.ParsePragmas(context.File.Doc)
	globalMacros := r.macros.Get(globalPragmas.Use()...)

	rewrites := len(globalMacros)
	copyDecl := copyDecls(context.File.Decls)
	for _, decl := range copyDecl {
		pragmas := pragma.ParsePragmas(loader.LoadComments(decl))

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
		err := r.fixImports(context)
		if err != nil {
			return err
		}
	}

	return r.printer.PrintFile(w, context.FileSet, context.File)
}

func (r ReWriter) fixImports(context macro.Context) error {
	specs := astutil.Imports(context.FileSet, context.File)
	absPath, err := filepath.Abs(r.source)
	if err != nil {
		return err
	}

	for _, group := range specs {
		for _, imprt := range group {
			importPath, _ := strconv.Unquote(imprt.Path.Value)
			if r.loaded.Packages.Has(importPath) {
				fsPath, err := filepath.Abs(r.loaded.Packages[importPath])
				if err != nil {
					return err
				}

				rel, err := filepath.Rel(absPath, fsPath)
				if err != nil {
					return err
				}
				rel = strings.ReplaceAll(rel, "\\", "/")

				newPath := strings.TrimSuffix(importPath, rel) // delete subpackage path
				newPath = strings.TrimSuffix(newPath, "/")
				newPath = strings.TrimSuffix(newPath, path.Base(newPath))  // delete target path
				newPath = path.Join(newPath, filepath.Base(r.output), rel) // replace target path

				astutil.RewriteImport(context.FileSet, context.File, importPath, newPath)
			}
		}
	}

	return nil
}
