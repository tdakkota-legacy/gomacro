package rewriter

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
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

// prepareOutputFile creates output file directory along with any necessary parents
// and returns absolute path to output file.
// Output file path is output+filepath.Rel(source, filePath)
func prepareOutputFile(source, output, filePath string) (string, error) {
	absPath, err := filepath.Abs(source)
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(absPath, filePath)
	if err != nil {
		return "", err
	}

	outputFile := filepath.Join(output, rel)
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
