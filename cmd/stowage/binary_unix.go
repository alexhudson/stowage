//go:build !windows

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// tryDispatch is a no-op on Unix: wrapped commands are installed as their
// own standalone shell scripts, so the stowage binary is never invoked
// under a wrapped command's name.
func tryDispatch() {}

func (b *Binary) install() error {
	wrapperFilePath := b.getWrapperPath()

	// remove any existing wrapper
	b.uninstall()

	if err := os.MkdirAll(filepath.Dir(wrapperFilePath), 0755); err != nil {
		panic(err)
	}

	meta := versionMeta()
	command := b.spec.runCommand()
	content := fmt.Sprintf(commandWrapper, meta, command)

	// install new wrapper
	err := ioutil.WriteFile(wrapperFilePath, []byte(content), 0755)
	if err != nil {
		panic(err)
	}

	return nil
}

func (b *Binary) uninstall() error {
	wrapperFilePath := b.getWrapperPath()

	if err := os.Remove(wrapperFilePath); err != nil && !os.IsNotExist(err) {
		fmt.Println("Couldn't remove wrapper!")
	}

	return nil
}

const commandWrapper = `#!/bin/sh
## %s ##
%s "$@"
`
