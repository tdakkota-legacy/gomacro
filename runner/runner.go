package runner

import (
	"errors"
	"io"

	macro "github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/internal/rewriter"
)

// Runner is a macro runner.
type Runner struct {
	Source, Output string
}

// Run runs given macros using path and writes result to output.
// If path is file, output should be file too.
// If path is dir, output should be dir too.
// If output dir does not exist, dir will be created.
func (r Runner) Run(macros macro.Macros) error {
	return r.run(macros, func(r rewriter.ReWriter) error {
		return r.Rewrite()
	})
}

// Print runs given macros using path and writes result to writer.
func (r Runner) Print(w io.Writer, macros macro.Macros) error {
	return r.run(macros, func(r rewriter.ReWriter) error {
		return r.RewriteTo(w)
	})
}

// Reader runs given macros using src reader and writes result to writer.
// This method would work in playground.
func (r Runner) Reader(name string, src io.Reader, w io.Writer, macros macro.Macros) error {
	return r.run(macros, func(r rewriter.ReWriter) error {
		return r.RewriteReader(name, src, w)
	})
}

func (r Runner) run(macros macro.Macros, f func(rewriter.ReWriter) error) error {
	path := r.Source
	if path == "" {
		path = "./"
	}

	writer := rewriter.NewReWriter(path, r.Output, macros, rewriter.DefaultPrinter())

	err := f(writer)
	if err != nil {
		if errors.Is(err, macro.ErrStop) {
			return nil
		}

		return err
	}

	return nil
}
