package macro

import (
	"go/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func testData() (*types.TypeName, Delayed) {
	typeName := types.NewTypeName(0, nil, "MyStruct", nil)
	return typeName, Delayed{
		"macro_test":  {typeName.Id(): {}},
		"macro_test2": {},
	}
}

func TestFind(t *testing.T) {
	a := require.New(t)
	typeName, d := testData()

	ok := d.Find("macro_test", typeName)
	a.True(ok)

	ok = d.Find("macro_test2", typeName)
	a.False(ok)
}

func TestAddTypeName(t *testing.T) {
	a := require.New(t)
	typeName, d := testData()

	d.AddTypeName("macro_test", typeName)
	a.Len(d, 2)
	a.Len(d["macro_test"], 1)

	d.AddTypeName("macro_test3", typeName)
	a.Len(d, 3)
	a.Len(d["macro_test3"], 1)
}
