package main

import (
	"os"
	"path/filepath"
	"runtime"
)

// installPrefix returns the base directory stowage installs into: wrapper
// binaries go in "bin/" under it, and spec/repo state goes under "share/stowage/".
func installPrefix() string {
	if runningInSelfInstallContainer() {
		return "/stowage/install-tree"
	}

	if prefix, ok := windowsInstallPrefix(); ok {
		return prefix
	}

	return unixInstallPrefix()
}

// globalBinDir returns the directory wrapper binaries are installed into
// for a regular (non-venv) install.
func globalBinDir() string {
	return filepath.Join(installPrefix(), "bin")
}

func runningInSelfInstallContainer() bool {
	_, err := os.Stat("/stowage/install-tree")
	return !os.IsNotExist(err)
}

// unixInstallPrefix defaults to a per-user install under $HOME/.local so
// stowage works without root/sudo, matching the XDG convention of
// $HOME/.local/bin. Running as root (e.g. via sudo) keeps the traditional
// system-wide /usr/local, preserving existing system-install workflows.
func unixInstallPrefix() string {
	if os.Geteuid() == 0 {
		return "/usr/local"
	}

	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".local")
	}

	return "/usr/local"
}

func windowsInstallPrefix() (string, bool) {
	if runtime.GOOS != "windows" {
		return "", false
	}
	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		return filepath.Join(localAppData, "stowage"), true
	}
	return "", false
}
