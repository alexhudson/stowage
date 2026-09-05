//go:build windows

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func wrapperExePath(b *Binary) string {
	return b.getWrapperPath() + ".exe"
}

func (b *Binary) install() error {
	wrapperFilePath := wrapperExePath(b)

	// remove any existing wrapper
	b.uninstall()

	if err := os.MkdirAll(filepath.Dir(wrapperFilePath), 0755); err != nil {
		panic(err)
	}

	selfPath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	// hard-link the stowage binary under the wrapped command's name; at
	// startup it recognizes it wasn't invoked as "stowage" and dispatches
	// to the matching spec instead. Fall back to a copy if the install
	// directory is on a different volume, since hard links can't cross those.
	if linkErr := os.Link(selfPath, wrapperFilePath); linkErr != nil {
		if copyErr := copyFile(selfPath, wrapperFilePath); copyErr != nil {
			panic(copyErr)
		}
	}

	return nil
}

func (b *Binary) uninstall() error {
	wrapperFilePath := wrapperExePath(b)

	if err := os.Remove(wrapperFilePath); err != nil && !os.IsNotExist(err) {
		fmt.Println("Couldn't remove wrapper!")
	}

	return nil
}

func copyFile(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// tryDispatch checks whether this process was invoked under a wrapped
// command's name (via its hard-linked .exe) rather than as "stowage" itself,
// and if so runs that command's spec against docker and exits - it never
// returns in that case.
func tryDispatch() {
	base := filepath.Base(os.Args[0])
	self := strings.TrimSuffix(base, filepath.Ext(base))
	if self == "" || self == "stowage" {
		return
	}

	spec, err := loadDispatchSpecification(self)
	if err != nil {
		return
	}

	args := spec.runCommandSlice()
	args = append(args, os.Args[1:]...)

	dockerPath, lookErr := exec.LookPath("docker")
	if lookErr != nil {
		fmt.Println("ERROR: docker not found in PATH")
		os.Exit(1)
	}

	os.Exit(runCommand(dockerPath, args[1:]))
}

// loadDispatchSpecification looks up a wrapped command's spec, preferring a
// venv-local spec store over the global one. The venv's spec store is
// derived structurally from this executable's own location (its wrapper
// lives at <venvDir>/bin/<name>.exe, spec store at
// <venvDir>/share/stowage/spec) rather than from the environment, so a
// hard-linked wrapper resolves its project-pinned version even if it's run
// directly rather than via an activated PATH. For a globally installed
// wrapper this resolves to the same directory as the global store, so the
// fallback lookup below is only needed if that structural guess fails.
func loadDispatchSpecification(name string) (Specification, error) {
	if selfPath, err := os.Executable(); err == nil {
		venvDir := filepath.Dir(filepath.Dir(selfPath))
		venvStore := newStorage(
			filepath.Join(venvDir, "share", "stowage", "spec"),
			filepath.Join(venvDir, "share", "stowage", "repo"),
		)
		if spec, err := venvStore.loadSpecificationByName(name); err == nil {
			return spec, nil
		}
	}

	globalStore := createStorage()
	return globalStore.loadSpecificationByName(name)
}
