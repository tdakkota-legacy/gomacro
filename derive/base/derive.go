package base

import (
	"go/types"

	builders "github.com/tdakkota/astbuilders"
)

type Dispatcher interface {
	Dispatch(field Field, typ types.Type, s builders.StatementBuilder) (builders.StatementBuilder, error)
}
