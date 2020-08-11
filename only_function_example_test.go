package macro_test

import (
	"go/ast"
	"io/ioutil"
	"os"

	builders "github.com/tdakkota/astbuilders"
	macro "github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/runner"
)

func Macro() macro.Handler {
	return macro.OnlyFunction("hello", func(ctx macro.Context, call *ast.CallExpr) error {
		args := append([]ast.Expr{builders.StringLit("hello ")}, call.Args...)
		toReplace := builders.CallName("println", args...)
		ctx.Replace(toReplace)
		return nil
	})
}

const src = `
package main

//procm:use=runme
func main() {
	hello()
}
`

func writeTempFile(src string) (string, error) {
	f, err := ioutil.TempFile("", "src.*.go")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.WriteString(src)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func Example_OnlyFunction() {
	srcPath, err := writeTempFile(src)
	if err != nil {
		panic(err)
	}

	err = runner.Print(srcPath, os.Stdout, macro.Macros{
		"runme": Macro(),
	})
	if err != nil {
		panic(err)
	}
	// Output: package main
	//
	// //procm:use=runme
	// func main() {
	// 	println("hello ")
	//
	// }
}
