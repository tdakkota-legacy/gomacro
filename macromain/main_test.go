package macromain

import (
	"bytes"
	"github.com/tdakkota/gomacro"
	"go/ast"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func CreateMacro(value string) macro.HandlerFunc {
	return func(cursor macro.Context, node ast.Node) error {
		if callExpr, ok := node.(*ast.CallExpr); ok {
			if f, ok := callExpr.Fun.(*ast.Ident); ok && f.Name == "eval" {
				for i := range callExpr.Args {
					switch v := callExpr.Args[i].(type) {
					case *ast.BasicLit:
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

func withTempFile(prefix string, cb func(file *os.File) error) error {
	f, err := ioutil.TempFile("", prefix+".*.go")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	defer f.Close()

	return cb(f)
}

var gotool = filepath.Join(runtime.GOROOT(), "bin", "go")

func TestRewriter(t *testing.T) {
	err := withTempFile("test-run", func(f *os.File) error {
		generatedValue := time.Now().String()
		err := macro.NewReWriter("../testdata/eval.go", "", macro.Macros{
			"eval": CreateMacro(generatedValue),
		}, macro.DefaultPrinter()).RewriteTo(f)

		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(nil)
		cmd := exec.Command(gotool, "run", f.Name())
		cmd.Stdout = buf
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}

		if generatedValue != buf.String() {
			t.Errorf("expected '%v', got '%v'", generatedValue, buf.String())
		}

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}
}
