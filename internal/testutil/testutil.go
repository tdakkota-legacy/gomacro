package testutil

import (
	"bytes"
	"go/ast"
	"go/token"
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

func GoToolRun(filename string) (string, error) {
	buf := bytes.NewBuffer(nil)
	cmd := exec.Command(GoTool, "run", filename)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

var GoTool = filepath.Join(runtime.GOROOT(), "bin", "go")
