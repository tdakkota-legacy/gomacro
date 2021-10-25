package rewriter

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/internal/loader"
	"github.com/tdakkota/gomacro/pragma"
)

// ReWriter is a Go source rewriter.
type ReWriter struct {
	source, output string
	macros         macro.Macros
	printer        Printer

	loaded loader.Loaded
}

// NewReWriter creates new ReWriter.
func NewReWriter(source, output string, macros macro.Macros, printer Printer) ReWriter {
	return ReWriter{
		source:  source,
		output:  output,
		macros:  macros,
		printer: printer,
	}
}

// Source returns source path.
func (r ReWriter) Source() string {
	return r.source
}

// Output returns output path.
func (r ReWriter) Output() string {
	return r.output
}

// Rewrite rewrites given source and output using macro list.
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

// RewriteTo rewrites one package using macros to given writer.
func (r ReWriter) RewriteTo(w io.Writer) error {
	ctx, err := loader.LoadOne(r.source)
	if err != nil {
		return err
	}

	return r.runMacro(w, ctx)
}

// RewriteReader loads sources from given reader and writes output to given writer.
func (r ReWriter) RewriteReader(name string, src io.Reader, w io.Writer) error {
	ctx, err := loader.LoadReader(src, name)
	if err != nil {
		return err
	}

	return r.runMacro(w, ctx)
}

func (r ReWriter) rewriteDir() error {
	return loader.LoadWalk(r.source, func(l loader.Loaded, ctx macro.Context) error {
		r.loaded = l
		file := ctx.FileSet.File(ctx.File.Pos()).Name()

		// skip generated files
		if strings.HasSuffix(file, ".gen.go") {
			return nil
		}

		outputFile, err := r.prepareOutputFile(file)
		if err != nil {
			return err
		}

		return r.rewriteOneFile(outputFile, ctx)
	})
}

func (r ReWriter) prepareOutputFile(file string) (string, error) {
	return prepareGenFile(file)
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

func (r ReWriter) runMacro(w io.Writer, context macro.Context) error {
	macroRunner := NewRunner(context.FileSet)
	globalPragmas := pragma.ParsePragmas(context.File.Doc)
	globalMacros := r.macros.Get(globalPragmas.Use()...)

	rewrites := len(globalMacros)
	copyDecl := copyDecls(context.File.Decls)
	context.File = copyFile(context.File)
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
		err := r.fixImports(true, context)
		if err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, "// Code generated by gomacro; DO NOT EDIT.\n"); err != nil {
		return err
	}

	return r.printer.PrintFile(w, context.FileSet, context.File)
}
