package derive

import "go/types"

func checkType(typ types.Type, typeName *types.TypeName) bool {
	if v, ok := typ.(*types.Named); ok {
		return v.Obj().Id() == typeName.Id()
	}

	return false
}

type container interface {
	Elem() types.Type
}
