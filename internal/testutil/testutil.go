package testutil

import (
	"go/ast"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	builders "github.com/tdakkota/astbuilders"
	macro "github.com/tdakkota/gomacro"
)

func call(node ast.Node, name string, cb func(callExpr *ast.CallExpr) error) error {
	if callExpr, ok := node.(*ast.CallExpr); ok {
		if ident, ok := callExpr.Fun.(*ast.Ident); ok && ident.Name == name {
			return cb(callExpr)
		}
	}

	return nil
}

// CreateTestMacroEval creates simple eval macro.
func CreateTestMacroEval(value string) macro.HandlerFunc {
	return func(cursor macro.Context, node ast.Node) error {
		return call(node, "eval", func(callExpr *ast.CallExpr) error {
			for i := range callExpr.Args {
				if v, ok := callExpr.Args[i].(*ast.BasicLit); ok {
					lit := builders.StringLit(value)
					lit.ValuePos = v.Pos()
					cursor.Replace(lit)
				}
			}

			return nil
		})
	}
}

// WithTempFile creates temporary file and calls cb.
// Then closes file and deletes it.
func WithTempFile(prefix string, cb func(file *os.File) error) error {
	f, err := ioutil.TempFile("", prefix+".*.go")
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	return cb(f)
}

// WithTempDir creates temporary directory and calls cb.
// Then deletes it.
func WithTempDir(prefix string, cb func(path string) error) error {
	dirPath, err := ioutil.TempDir("", prefix)
	if err != nil {
		return err
	}
	defer func() {
		_ = os.RemoveAll(dirPath)
	}()

	return cb(dirPath)
}

// CopyFile copies file.
func CopyFile(to, from string) error {
	dst, err := os.Create(to)
	if err != nil {
		return err
	}
	defer func() {
		_ = dst.Close()
	}()

	src, err := os.Open(filepath.Clean(from))
	if err != nil {
		return err
	}
	defer func() {
		_ = src.Close()
	}()

	_, err = io.Copy(dst, src)
	return err
}

// GoTool is a path to go tool.
var GoTool = filepath.Join(runtime.GOROOT(), "bin", "go") //nolint:gochecknoglobals

// RunGoTool runs go command with given args.
func RunGoTool(output io.Writer, args ...string) error {
	cmd := exec.Command(GoTool, args...) // #nosec G204
	cmd.Stdin = os.Stdin
	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GoRun runs given filename using GoTool and returns stdout output.
func GoRun(filename string) (string, error) {
	buf := strings.Builder{}

	if err := RunGoTool(&buf, "run", filename); err != nil {
		return "", err
	}

	return buf.String(), nil
}
