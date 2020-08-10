package base

import (
	"go/types"

	builders "github.com/tdakkota/astbuilders"
)

type Derive interface {
	Dispatch(field Field, typ types.Type, s builders.StatementBuilder) (builders.StatementBuilder, error)
	Basic(field Field, typ *types.Basic, s builders.StatementBuilder) (builders.StatementBuilder, error)
	Array(field Field, typ *types.Array, s builders.StatementBuilder) (builders.StatementBuilder, error)
	Slice(field Field, typ *types.Slice, s builders.StatementBuilder) (builders.StatementBuilder, error)
	Struct(field Field, typ *types.Struct, s builders.StatementBuilder) (builders.StatementBuilder, error)
	Pointer(field Field, typ *types.Pointer, s builders.StatementBuilder) (builders.StatementBuilder, error)
	Interface(field Field, typ *types.Interface, s builders.StatementBuilder) (builders.StatementBuilder, error)
	Map(field Field, typ *types.Map, s builders.StatementBuilder) (builders.StatementBuilder, error)
	Chan(field Field, typ *types.Chan, s builders.StatementBuilder) (builders.StatementBuilder, error)
}
