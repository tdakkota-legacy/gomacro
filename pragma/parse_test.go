package pragma

import (
	"fmt"
	"go/ast"
	"testing"

	"github.com/stretchr/testify/require"
)

//nolint:scopelint
func Test_parsePragma(t *testing.T) {
	tests := []struct {
		pair      string
		wantKey   string
		wantValue string
		wantOk    bool
	}{
		{
			pair:      "key=value",
			wantKey:   "key",
			wantValue: "value",
			wantOk:    true,
		},
		{
			pair:      "key=value===!",
			wantKey:   "key",
			wantValue: "value===!",
			wantOk:    true,
		},
		{
			pair:      "key",
			wantKey:   "key",
			wantValue: "",
			wantOk:    true,
		},
		{
			pair:      "",
			wantKey:   "",
			wantValue: "",
			wantOk:    false,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s,%t", tt.pair, tt.wantOk), func(t *testing.T) {
			a := require.New(t)
			key, value, ok := parsePragma(tt.pair)
			a.Equal(tt.wantKey, key)
			a.Equal(tt.wantValue, value)
			a.Equal(tt.wantOk, ok)
		})
	}
}

func TestParsePragmas(t *testing.T) {
	comment := &ast.CommentGroup{List: []*ast.Comment{
		{Text: `//procm:macro=value`},
		{Text: `//procm:key=value`},
	}}

	a := require.New(t)
	pragmas := ParsePragmas(comment)
	a.Len(pragmas, 2)
	a.Equal("value", pragmas["macro"])
	a.Equal("value", pragmas["key"])
	a.Equal("", pragmas["key2"])
}
