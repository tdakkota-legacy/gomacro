package macromain

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/tdakkota/gomacro"
)

func Main(macros macro.Macros) {
	flag.Parse()
	if err := run(flag.Arg(0), flag.Arg(1), macros); err != nil {
		fmt.Println(err)
		os.Exit(-1)
		return
	}
}

func run(path, output string, macros macro.Macros) error {
	if path == "" {
		path = "./"
	}

	err := macro.NewReWriter(path, output, macros, macro.DefaultPrinter()).Rewrite()
	if err != nil {
		if errors.Is(err, macro.ErrStop) {
			return nil
		}

		return err
	}

	return nil
}
