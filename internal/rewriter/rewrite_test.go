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
	"golang.org/x/tools/go/ast/astutil"

	"github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/internal/loader"
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

	output, err := testutil.GoRun(runPath)
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
		"github.com/tdakkota/go-terra/proto":              "./proto",
		"github.com/tdakkota/go-terra/proto/structs/tile": "./proto/structs/tile",
	}

	fset := token.NewFileSet()
	f := &ast.File{}
	astutil.AddImport(fset, f, "github.com/tdakkota/go-terra/proto")
	astutil.AddImport(fset, f, "github.com/tdakkota/go-terra/proto/structs/tile")

	err := r.fixImports(false, macro.Context{
		ASTInfo: macro.ASTInfo{File: f, FileSet: fset},
	})
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, `"github.com/tdakkota/go-terra/proto_out"`, f.Imports[0].Path.Value)
	require.Equal(t, `"github.com/tdakkota/go-terra/proto_out/structs/tile"`, f.Imports[1].Path.Value)
}

func TestReWriterGetter(t *testing.T) {
	source, output := "s", "o"
	r := ReWriter{
		source: source,
		output: output,
	}
	require.Equal(t, source, r.Source())
	require.Equal(t, output, r.Output())
}

func TestReWriter_RewriteReader(t *testing.T) {
	err := testutil.WithTempFile("test-run", func(f *os.File) error {
		path := testDataPath("src/eval.go")
		src, err := os.Open(path)
		if err != nil {
			return err
		}

		return runRewriteTest(path, f.Name(), f.Name(), func(writer ReWriter) error {
			return writer.RewriteReader(f.Name(), src, f)
		})
	})

	if err != nil {
		t.Fatal(err)
	}
}
