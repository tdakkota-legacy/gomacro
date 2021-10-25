package derive

import (
	"go/ast"
	"go/types"
	"reflect"
	"strings"
)

// Tag is a struct tag
type Tag string

// Lookup gets struct tag value by name.
func (t Tag) Lookup(s string) (string, bool) {
	split := strings.Split(string(t), ",")

	for i := range split {
		if v, ok := reflect.StructTag(split[i]).Lookup(s); ok {
			return v, true
		}
	}

	return "", false
}

// Field is a struct field representation.
type Field struct {
	Named    bool
	TypeName *types.TypeName
	Tag      Tag
	Selector ast.Expr
}
