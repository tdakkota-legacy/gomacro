package testutil

import (
	"bytes"
	"go/ast"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/tdakkota/gomacro"
)

func CreateMacro(value string) macro.HandlerFunc {
	return func(cursor macro.Context, node ast.Node) error {
		if callExpr, ok := node.(*ast.CallExpr); ok {
			if f, ok := callExpr.Fun.(*ast.Ident); ok && f.Name == "eval" {
				for i := range callExpr.Args {
					if v, ok := callExpr.Args[i].(*ast.BasicLit); ok {
						cursor.Replace(&ast.BasicLit{
							ValuePos: v.Pos(),
							Kind:     token.INT,
							Value:    strconv.Quote(value),
						})
					}
				}

			}
			return nil
		}

		return nil
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
	return buf.String(), RunGoTool(buf, "run", filename)
}
