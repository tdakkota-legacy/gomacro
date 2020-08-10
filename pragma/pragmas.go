package pragma

import "strings"

const (
	UsePragma   = "use"
	MacroPragma = "macro"
)

type Pragmas map[string]string

func (p Pragmas) Use() []string {
	r := strings.Split(p[UsePragma], ",")
	if len(r) == 1 && r[0] == "" {
		return nil
	}
	return r
}

func (p Pragmas) Macro() (string, bool) {
	v, ok := p[MacroPragma]

	return v, ok
}
