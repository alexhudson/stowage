package main

import (
	"path/filepath"
)

// Binary is a type
type Binary struct {
	name string
	spec *Specification
}

func (b *Binary) getWrapperPath() string {
	return filepath.Join(installPrefix(), "bin", b.name)
}
