# gomacro

[![Go](https://github.com/tdakkota/gomacro/workflows/Go/badge.svg)](https://github.com/tdakkota/gomacro/actions)
[![Documentation](https://godoc.org/github.com/tdakkota/gomacro?status.svg)](https://pkg.go.dev/github.com/tdakkota/gomacro)
[![codecov](https://codecov.io/gh/tdakkota/gomacro/branch/master/graph/badge.svg)](https://codecov.io/gh/tdakkota/gomacro)
[![license](https://img.shields.io/github/license/tdakkota/gomacro.svg)](https://github.com/tdakkota/gomacro/blob/master/LICENSE)

Procedural macro toolbox for Golang.

## Install
```
go get github.com/tdakkota/gomacro
```

## Example
```go
package main

import (
	"go/ast"
	"io/ioutil"
	"os"

	builders "github.com/tdakkota/astbuilders"
	macro "github.com/tdakkota/gomacro"
	"github.com/tdakkota/gomacro/runner"
)

// hello() -> println("Hello, ", "World")
// hello(arg) -> println("Hello, ", arg)
// hello(args...) -> println("Hello, ", args...)
func Macro() macro.Handler {
	return macro.OnlyFunction("hello", func(ctx macro.Context, call *ast.CallExpr) error {
		args := append([]ast.Expr{builders.StringLit("hello, ")}, call.Args...)
		if len(call.Args) == 0 {
			args = append(args, builders.StringLit("world"))
		}

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

func main() {
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
}
```

Output:
```go
package main        
                            
//procm:use=runme           
func main() {               
	println("hello, ", "world")          
}                           
```