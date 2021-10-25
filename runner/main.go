package runner

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tdakkota/gomacro"
)

// Main parses source and output path from flags and calls Runner's functions.
// In case if both arguments are empty, it would read source from stdin and write output to stdout.
// In case if one argument provided, it would read source from given path and write output to stdout.
// In case if both arguments are provided, it would read source from first argument path
// and write output to second argument path.
func Main(macros macro.Macros) {
	flag.Parse()
	input, output := flag.Arg(0), flag.Arg(1)

	switch {
	case input != "" && output != "":
		if err := Run(input, output, macros); err != nil {
			fmt.Println(err)
			return
		}
	case input != "":
		if err := Print(input, os.Stdout, macros); err != nil {
			fmt.Println(err)
			return
		}
	default:
		if err := Reader("stdin", os.Stdin, os.Stdout, macros); err != nil {
			fmt.Println(err)
			return
		}
	}
}

// Run runs given macros using path and writes result to output.
// If path is file, output should be file too.
// If path is dir, output should be dir too.
// If output dir does not exist, dir will be created.
func Run(path, output string, macros macro.Macros) error {
	return Runner{
		Source: path,
		Output: output,
	}.Run(macros)
}

// Print runs given macros using path and writes result to writer.
func Print(path string, w io.Writer, macros macro.Macros) error {
	return Runner{Source: path}.Print(w, macros)
}

// Reader runs given macros using src reader and writes result to writer.
func Reader(name string, src io.Reader, w io.Writer, macros macro.Macros) error {
	return Runner{}.Reader(name, src, w, macros)
}
