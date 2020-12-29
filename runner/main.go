package runner

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/runner/flags"
)

// Main parses source and output path from flags and calls Run function.
func Main(macros macro.Macros) {
	flag.Parse()
	input, output := flag.Arg(0), flag.Arg(1)

	switch {
	case input != "" && output != "":
		if err := Run(flag.Arg(0), flag.Arg(1), macros); err != nil {
			fmt.Println(err)
			return
		}
	case input != "":
		if err := Print(flag.Arg(0), os.Stdout, macros); err != nil {
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
		Flags:  flags.AppendMode | flags.AddGeneratedComment,
	}.Run(macros)
}

// Print runs given macros using path and writes result to writer.
func Print(path string, w io.Writer, macros macro.Macros) error {
	return Runner{Source: path, Flags: defaultFlags}.Print(w, macros)
}

// Reader runs given macros using src reader and writes result to writer.
func Reader(name string, src io.Reader, w io.Writer, macros macro.Macros) error {
	return Runner{Flags: defaultFlags}.Reader(name, src, w, macros)
}
