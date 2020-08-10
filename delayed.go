package macro

import (
	"go/ast"
	"go/types"

	"github.com/tdakkota/gomacro/pragma"
	"golang.org/x/tools/go/packages"
)

type Delayed map[string]DelayedMacro

type DelayedMacro map[string]struct{}

func (d Delayed) Find(macroName string, typeName *types.TypeName) bool {
	if m, ok := d[macroName]; ok {
		return m.Find(typeName)
	}
	return false
}

func (d DelayedMacro) Find(typeName *types.TypeName) bool {
	if _, ok := d[typeName.Name()]; ok {
		return true
	}

	return false
}

func (d Delayed) AddTypeName(macroName string, typeName *types.TypeName) {
	if d[macroName] == nil {
		d[macroName] = map[string]struct{}{}
	}

	d[macroName][typeName.Name()] = struct{}{}
}

func (d DelayedMacro) AddTypeName(typeName *types.TypeName) {
	d[typeName.Name()] = struct{}{}
}

func (d Delayed) Add(pkg *packages.Package) {
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				m := pragma.ParsePragmas(genDecl.Doc).Use()

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					for _, macroName := range m {
						if d[macroName] == nil {
							d[macroName] = map[string]struct{}{}
						}

						obj := pkg.TypesInfo.ObjectOf(typeSpec.Name)
						d[macroName][obj.Name()] = struct{}{}
					}
				}
			}
		}
	}
}
