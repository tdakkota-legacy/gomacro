package base

import (
	"go/ast"
	"go/types"
	"reflect"
)

type Field struct {
	TypeName *types.TypeName
	Tag      reflect.StructTag
	Selector ast.Expr
}
