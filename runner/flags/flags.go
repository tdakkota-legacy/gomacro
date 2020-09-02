package flags

type Flags uint

const (
	AppendMode Flags = 1 << iota
	AddGeneratedComment
)

func (b Flags) Set(flag Flags) Flags { return b | flag }
func (b Flags) Has(flag Flags) bool  { return b&flag != 0 }
