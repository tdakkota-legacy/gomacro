package rewriter

import (
	"errors"
	"go/ast"
	"os"
	"path"
	"path/filepath"
	"strings"

	builders "github.com/tdakkota/astbuilders"
)

var ErrNotRelative = errors.New("paths is not relative")

func urlRel(basePath, targetPath string) (string, error) {
	base, target := path.Clean(basePath), path.Clean(targetPath)
	if base == target {
		return base, nil
	}

	var b string
	for {
		target, b = path.Split(target)
		target = path.Clean(target)

		if base == target {
			return path.Clean(strings.TrimPrefix(targetPath, base)), nil
		} else if target == b {
			return "", ErrNotRelative
		}
	}
}

// Output path is to+filepath.Rel(from, rel)
func changePathBase(from, to, rel string) (string, error) {
	relative, err := filepath.Rel(from, rel)
	if err != nil {
		return "", err
	}

	return filepath.Join(to, relative), nil
}

// prepareOutputFile creates output file directory along with any necessary parents
// and returns absolute path to output file.
// Output file path is output+filepath.Rel(source, filePath)
func prepareOutputFile(source, output, filePath string) (string, error) {
	absPath, err := filepath.Abs(source)
	if err != nil {
		return "", err
	}

	outputFile, err := changePathBase(absPath, output, filePath)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Dir(outputFile), os.ModePerm)
	if err != nil {
		return "", err
	}

	return outputFile, nil
}

func prepareGenFile(filePath string) (string, error) {
	ext := filepath.Ext(filePath)
	name := strings.TrimSuffix(filePath, ext)
	return name + ".gen" + ext, nil
}

func copyDecls(a []ast.Decl) []ast.Decl {
	copyDecl := make([]ast.Decl, len(a))
	copy(copyDecl, a)
	return copyDecl
}

func copyFile(file *ast.File) *ast.File {
	return builders.NewFileBuilder(file.Name.Name).
		AddImports(file.Imports...).
		Complete()
}
