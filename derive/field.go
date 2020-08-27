package derive

import (
	"go/ast"
	"go/types"
	"reflect"
	"strings"
)

type Tag string

func (t Tag) Lookup(s string) (string, bool) {
	split := strings.Split(string(t), ",")

	for i := range split {
		if v, ok := reflect.StructTag(split[i]).Lookup(s); ok {
			return v, true
		}
	}

	return "", false
}

type Field struct {
	Named    bool
	TypeName *types.TypeName
	Tag      Tag
	Selector ast.Expr
}
