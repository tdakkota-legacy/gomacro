package base

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTag_Lookup(t *testing.T) {
	table := []struct {
		tag, key string
		expected string
	}{
		{`if:"$m.predicate >= 10"`, "if", `$m.predicate >= 10`},
		{`if:"$m.predicate >= 10", empty:"1"`, "if", `$m.predicate >= 10`},
		{`if:"$m.predicate >= 10", empty:"1"`, "empty", `1`},
	}

	a := require.New(t)
	for _, test := range table {
		tag, _ := Tag(test.tag).Lookup(test.key)
		a.Equal(test.expected, tag)
	}
}
