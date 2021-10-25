package main

import (
	"testing"

	"github.com/tdakkota/gomacro/internal/testutil"
)

func Test_run(t *testing.T) {
	t.Skip("FIXME: Imports fails for some reason")

	err := testutil.WithTempDir("gomacrotest", func(path string) error {
		return run(`./testdata/proto`, path)
	})

	if err != nil {
		t.Fatal(err)
	}
}
