package main

import (
	"log"

	macro "github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/derive/cursor"
	"github.com/tdakkota/gomacro/runner"
)

func main() {
	m, err := cursor.Create()
	if err != nil {
		log.Fatal(err)
	}

	runner.Main(macro.Macros{
		"derive_binary": m,
	})
}

func run(path, output string) error {
	m, err := cursor.Create()
	if err != nil {
		log.Fatal(err)
	}

	return runner.Run(path, output, macro.Macros{
		"derive_binary": m,
	})
}
