package macro

import (
	"errors"
	"fmt"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/tdakkota/gomacro/pragma"
)

// ErrStop is special error to interrupt AST traversal.
var ErrStop = errors.New("macro exit")

// Context is a macro context.
type Context struct {
	*astutil.Cursor
	Pre     bool
	Delayed Delayed
	Report  func(Report)
	// AST objects.
	ASTInfo
	// Type checker objects.
	TypeInfo
	// Parsed magic comments.
	Pragmas pragma.Pragmas
}

// Report represents macro error report.
type Report struct {
	Pos     token.Pos
	Message string
}

// Reportf reports macro error.
func (c Context) Reportf(pos token.Pos, format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	c.Report(Report{Pos: pos, Message: msg})
}
