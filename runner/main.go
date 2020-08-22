package runner

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/internal/rewriter"
)

// Main parses source and output path from flags and calls Run function.
func Main(macros macro.Macros) {
	flag.Parse()
	if err := Run(flag.Arg(0), flag.Arg(1), macros); err != nil {
		fmt.Println(err)
		return
	}
}

// Run runs given macros using path and writes result to output.
// If path is file, output should be file too.
// If path is dir, output should be dir too.
// If output dir does not exist, dir will be created.
func Run(path, output string, macros macro.Macros) error {
	return run(path, output, macros, func(r rewriter.ReWriter) error {
		return r.Rewrite()
	})
}

// Run runs given macros using path and writes result to writer.
func Print(path string, w io.Writer, macros macro.Macros) error {
	return run(path, "", macros, func(r rewriter.ReWriter) error {
		return r.RewriteTo(w)
	})
}

func run(path, output string, macros macro.Macros, f func(rewriter.ReWriter) error) error {
	if path == "" {
		path = "./"
	}

	err := f(rewriter.NewReWriter(path, output, macros, rewriter.DefaultPrinter()))
	if err != nil {
		if errors.Is(err, macro.ErrStop) {
			return nil
		}

		return err
	}

	return nil
}
