package macro

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/tdakkota/gomacro/pragma"
)

// Delayed is mapping from macro name to DelayedTypes
type Delayed map[string]DelayedTypes

// DelayedTypes is a list of types will be changed by macros.
type DelayedTypes map[string]struct{}

func id(typeName types.Object) string {
	return typeName.Id()
}

// Find returns true if given typeName will be changed by macroName.
func (d Delayed) Find(macroName string, typeName *types.TypeName) bool {
	if m, ok := d[macroName]; ok {
		return m.Find(typeName)
	}
	return false
}

// Find returns true if given typeName will be changed by macro.
func (d DelayedTypes) Find(typeName *types.TypeName) bool {
	if _, ok := d[id(typeName)]; ok {
		return true
	}

	return false
}

// AddTypeName adds type name to list.
func (d Delayed) AddTypeName(macroName string, typeName *types.TypeName) {
	if d[macroName] == nil {
		d[macroName] = map[string]struct{}{}
	}

	d[macroName].AddTypeName(typeName)
}

// AddTypeName adds type name to list.
func (d DelayedTypes) AddTypeName(typeName *types.TypeName) {
	d[id(typeName)] = struct{}{}
}

func getTypeName(typesInfo *types.Info, spec ast.Spec) (*types.TypeName, bool) {
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return nil, ok
	}

	obj := typesInfo.ObjectOf(typeSpec.Name)
	v, ok := obj.(*types.TypeName)
	return v, ok
}

// AddDecls adds type names from declarations.
// If declaration have a magic comment, type name will be added.
func (d Delayed) AddDecls(typesInfo *types.Info, decls []ast.Decl) {
	for _, decl := range decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			m := pragma.ParsePragmas(genDecl.Doc).Use()

			for _, spec := range genDecl.Specs {
				v, ok := getTypeName(typesInfo, spec)
				if !ok {
					continue
				}

				for _, macroName := range m {
					if d[macroName] == nil {
						d[macroName] = DelayedTypes{}
					}

					d[macroName].AddTypeName(v)
				}
			}
		}
	}
}

// Add adds type names from package.
// If declaration have a magic comment, type name will be added.
func (d Delayed) Add(pkg *packages.Package) {
	for _, file := range pkg.Syntax {
		d.AddDecls(pkg.TypesInfo, file.Decls)
	}
}
