package main

import (
	macro "github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/derive/cursor"
	"github.com/tdakkota/gomacro/runner"
)

func main() {
	runner.Main(macro.Macros{
		"derive_binary": macro.HandlerFunc(cursor.DeriveBinary),
	})
}
