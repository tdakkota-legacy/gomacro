package runner

import (
	"flag"
	"fmt"
	"io"

	"github.com/tdakkota/gomacro"
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
	return Runner{path, output}.Run(macros)
}

// Run runs given macros using path and writes result to writer.
func Print(path string, w io.Writer, macros macro.Macros) error {
	return Runner{path, ""}.Print(w, macros)
}
