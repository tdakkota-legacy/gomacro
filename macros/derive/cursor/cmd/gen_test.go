package main

import (
	"testing"

	"github.com/tdakkota/gomacro/internal/testutil"
)

func Test_run(t *testing.T) {
	err := testutil.WithTempDir("gomacrotest", func(path string) error {
		return run(`./testdata/proto`, path)
	})

	if err != nil {
		t.Fatal(err)
	}
}
