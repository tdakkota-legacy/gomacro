package testutil

import (
	"bytes"
	"go/ast"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	builders "github.com/tdakkota/astbuilders"
	"github.com/tdakkota/gomacro"
)

func call(node ast.Node, name string, cb func(callExpr *ast.CallExpr) error) error {
	if callExpr, ok := node.(*ast.CallExpr); ok {
		if ident, ok := callExpr.Fun.(*ast.Ident); ok && ident.Name == name {
			return cb(callExpr)
		}
	}

	return nil
}

func CreateMacro(value string) macro.HandlerFunc {
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

func WithTempFile(prefix string, cb func(file *os.File) error) error {
	f, err := ioutil.TempFile("", prefix+".*.go")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()

	return cb(f)
}

func WithTempDir(prefix string, cb func(path string) error) error {
	dirPath, err := ioutil.TempDir("", prefix)
	if err != nil {
		return err
	}
	defer os.RemoveAll(dirPath)

	return cb(dirPath)
}

var GoTool = filepath.Join(runtime.GOROOT(), "bin", "go")

func RunGoTool(output io.Writer, args ...string) error {
	cmd := exec.Command(GoTool, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = output
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GoRun(filename string) (string, error) {
	buf := bytes.NewBuffer(nil)

	if err := RunGoTool(buf, "run", filename); err != nil {
		return "", err
	}

	return buf.String(), nil
}
