package rewriter

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tdakkota/gomacro"

	"github.com/tdakkota/gomacro/internal/testutil"
)

func testDataPath(path string) string {
	return filepath.Join("./testdata/", path)
}

func createRewriter(path, output string) (string, ReWriter) {
	generatedValue := time.Now().String()
	return generatedValue, NewReWriter(path, output, macro.Macros{
		"eval": testutil.CreateMacro(generatedValue),
	}, DefaultPrinter())
}

func runRewriteTest(path, outputPath, runPath string, cb func(ReWriter) error) error {
	generatedValue, reWriter := createRewriter(path, outputPath)
	err := cb(reWriter)
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
		return runRewriteTest(testDataPath("src/eval.go"), f.Name(), f.Name(), func(writer ReWriter) error {
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

		return runRewriteTest(testDataPath("src/eval.go"), outputFile, outputFile, func(writer ReWriter) error {
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

		return runRewriteTest("./testdata/src", path, outputFile, func(writer ReWriter) error {
			return writer.Rewrite()
		})
	})

	if err != nil {
		t.Fatal(err)
	}
}
