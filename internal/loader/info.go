package loader

import "golang.org/x/tools/go/packages"

// Loaded contains loaded module information.
type Loaded struct {
	Packages LoadedPackages
	Module   *packages.Module
}
