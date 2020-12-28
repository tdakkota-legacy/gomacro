package pragma

import (
	"fmt"
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
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
			key, value, ok := parsePragma(tt.pair)
			assert.Equal(t, tt.wantKey, key)
			assert.Equal(t, tt.wantValue, value)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestParsePragmas(t *testing.T) {
	comment := &ast.CommentGroup{List: []*ast.Comment{
		{Text: `//procm:macro=value`},
		{Text: `//procm:key=value`},
	}}

	pragmas := ParsePragmas(comment)
	assert.Len(t, pragmas, 2)
	v, ok := pragmas.Macro()
	assert.Equal(t, "value", v)
	assert.True(t, ok)
	assert.Equal(t, "value", pragmas["macro"])
	assert.Equal(t, "value", pragmas["key"])
	assert.Equal(t, "", pragmas["key2"])
}
