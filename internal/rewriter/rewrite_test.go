package rewriter

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tdakkota/gomacro/internal/loader"
	"golang.org/x/tools/go/ast/astutil"

	"github.com/tdakkota/gomacro"

	"github.com/tdakkota/gomacro/internal/testutil"
)

func testDataPath(path string) string {
	return filepath.Join("testdata", path)
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

func Test_fixImports(t *testing.T) {
	_, r := createRewriter("./proto", "./proto_out")
	r.loaded.Packages = loader.LoadedPackages{
		"github.com/tdakkota/go-terra/proto": "./proto",
	}

	fset := token.NewFileSet()
	f := &ast.File{}
	astutil.AddImport(fset, f, "github.com/tdakkota/go-terra/proto")

	err := r.fixImports(macro.Context{
		ASTInfo: macro.ASTInfo{File: f, FileSet: fset},
	})
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, `"github.com/tdakkota/go-terra/proto_out"`, f.Imports[0].Path.Value)
}
