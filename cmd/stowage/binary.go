package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// Binary is a type
type Binary struct {
	name string
	spec *Specification
}

func (b *Binary) getWrapperPath() string {
	if runtime.GOOS == "windows" {
		userhome, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		prefix := filepath.Join(userhome, ".stowage")
		// if _, err := os.Stat(prefix); !os.IsNotExist(err) {
		if err := os.MkdirAll(prefix, os.ModePerm); err != nil {
			log.Fatal(err)
		}
		// }
		return filepath.Join(prefix, b.name+".bat")
	} else {
		prefix := "/usr/local"
		if _, err := os.Stat("/stowage/install-tree"); !os.IsNotExist(err) {
			prefix = "/stowage/install-tree"
		}

		return filepath.Join(prefix, "bin", b.name)
	}
}

func (b *Binary) install() error {
	wrapperFilePath := b.getWrapperPath()

	// remove any existing wrapper
	b.uninstall()

	meta := versionMeta()
	command := b.spec.runCommand()
	content := ""
	if runtime.GOOS == "windows" {
		content = fmt.Sprintf(commandWrapperWin, meta, command)
	} else {
		content = fmt.Sprintf(commandWrapper, meta, command)
	}

	// install new wrapper
	err := ioutil.WriteFile(wrapperFilePath, []byte(content), 0755)
	if err != nil {
		panic(err)
	}

	return nil
}

func (b *Binary) uninstall() error {
	wrapperFilePath := b.getWrapperPath()

	if err := os.Remove(wrapperFilePath); os.IsExist(err) {
		fmt.Println("Couldn't remove wrapper!")
	}

	return nil
}

const commandWrapper = `#!/bin/sh
## %s ##
%s "$@"
`

const commandWrapperWin = `@echo off
REM ## %s ##
%s %%*
endlocal
`
