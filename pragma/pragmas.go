package pragma

import "strings"

const (
	// UsePragma is a name of pragma which denotes to run macro.
	UsePragma = "use"
)

// Pragmas contains Go AST annotations.
type Pragmas map[string]string

// Use gets slice of macro names to run.
func (p Pragmas) Use() []string {
	r := strings.Split(p[UsePragma], ",")
	if len(r) == 1 && r[0] == "" {
		return nil
	}
	return r
}
