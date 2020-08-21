package loader

import "golang.org/x/tools/go/packages"

type Loaded struct {
	Packages LoadedPackages
	Module   *packages.Module
}
