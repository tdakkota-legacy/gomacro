package macromain

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/testutil"
)

func createRewriter(path, output string) (string, macro.ReWriter) {
	generatedValue := time.Now().String()
	return generatedValue, macro.NewReWriter(path, output, macro.Macros{
		"eval": testutil.CreateMacro(generatedValue),
	}, macro.DefaultPrinter())
}

func runRewriteTest(path, outputPath, runPath string, cb func(macro.ReWriter) error) error {
	generatedValue, rewriter := createRewriter(path, outputPath)
	err := cb(rewriter)
	if err != nil {
		return err
	}

	output, err := testutil.GoToolRun(runPath)
	if err != nil {
		return err
	}

	if generatedValue != output {
		return fmt.Errorf("expected '%v', got '%v'", generatedValue, output)
	}

	return nil
}

func TestRewriteTo(t *testing.T) {
	err := testutil.WithTempFile("test-run", func(f *os.File) error {
		return runRewriteTest("../testdata/src/eval.go", f.Name(), f.Name(), func(writer macro.ReWriter) error {
			return writer.RewriteTo(f)
		})
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestRewriteFile(t *testing.T) {
	err := testutil.WithTempDir("gomacrotest", func(path string) error {
		outputFile := filepath.Join(path, "eval.go")

		return runRewriteTest("../testdata/src/eval.go", outputFile, outputFile, func(writer macro.ReWriter) error {
			return writer.Rewrite()
		})
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestRewriteDir(t *testing.T) {
	err := testutil.WithTempDir("gomacrotest", func(path string) error {
		outputFile := filepath.Join(path, "eval.go")

		return runRewriteTest("../testdata/src", path, outputFile, func(writer macro.ReWriter) error {
			return writer.Rewrite()
		})
	})

	if err != nil {
		t.Fatal(err)
	}
}
