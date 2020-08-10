package macro

import (
	"github.com/tdakkota/gomacro/macroctx"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

type Runner struct {
	errorPrinter
}

func NewRunner(fset *token.FileSet) Runner {
	return Runner{
		errorPrinter: errorPrinter{NewErrorCallback(fset)},
	}
}

func (r Runner) Run(handler Handler, context macroctx.Context, node ast.Node) ast.Node {
	context.Report = func(report macroctx.Report) {
		r.Reportf(report.Pos, report.Message)
	}

	return astutil.Apply(node, func(cursor *astutil.Cursor) bool {
		context.Pre = true
		context.Cursor = cursor

		if v := cursor.Node(); v != nil {
			err := handler.Handle(context, v)
			if err != nil {
				return false
			}
		}
		return true
	}, r.post(handler, context))
}

func (r Runner) post(handler Handler, context macroctx.Context) astutil.ApplyFunc {
	context.Report = func(report macroctx.Report) {
		r.Reportf(report.Pos, report.Message)
	}

	return func(cursor *astutil.Cursor) bool {
		if v := cursor.Node(); v != nil {
			context.Cursor = cursor

			err := handler.Handle(context, v)
			if err != nil {
				r.err(v.Pos(), err)
				return false
			}
		}

		return true
	}
}
