package main

import (
	"fmt"

	"github.com/urfave/cli"
)

// installIntoEnvironment resolves entry's request the same way `stowage
// install`/`install-spec` do, applies any recorded name/entrypoint
// override, and installs the result into env.
func installIntoEnvironment(env venv, entry venvManifestEntry) (*Specification, error) {
	installer := Installer{Request: entry.Request, RequestSpec: entry.SpecRequest}
	if !installer.setup() {
		req := entry.Request
		if req == "" {
			req = entry.SpecRequest
		}
		return nil, fmt.Errorf("don't know how to install %s", req)
	}

	if entry.Name != "" {
		installer.setCommandName(entry.Name)
	}
	if entry.Entrypoint != "" {
		installer.setEntrypoint(entry.Entrypoint)
	}

	if err := env.install(installer.Spec); err != nil {
		return nil, err
	}

	return installer.Spec, nil
}

func cmdVenvCreate(c *cli.Context) error {
	dir := c.Args().First()
	if dir == "" {
		fmt.Println("ERROR: specify a directory for the environment, e.g. stowage venv .venv")
		return nil
	}

	env, err := newVenv(dir)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return nil
	}

	if err := env.ensureDirs(); err != nil {
		fmt.Printf("ERROR: couldn't create environment: %s\n", err)
		return nil
	}

	manifestPath := c.String("manifest")
	if manifestPath == "" {
		manifestPath = defaultVenvManifestName
	}

	if manifest, err := loadVenvManifest(manifestPath); err == nil {
		for _, entry := range manifest.Entries {
			if _, err := installIntoEnvironment(env, entry); err != nil {
				fmt.Printf("ERROR: couldn't install %s from %s: %s\n", entry.Name, manifestPath, err)
				continue
			}
			fmt.Printf("%s installed\n", entry.Name)
		}
	}

	if err := env.writeActivationScripts(); err != nil {
		fmt.Printf("ERROR: couldn't write activation scripts: %s\n", err)
		return nil
	}

	fmt.Printf("Environment created at %s\n", env.Dir)
	fmt.Println("Activate it with one of:")
	fmt.Printf("  source %s/bin/activate    (bash/zsh)\n", dir)
	fmt.Printf("  %s\\bin\\Activate.ps1        (PowerShell)\n", dir)
	fmt.Printf("  %s\\bin\\activate.bat        (cmd.exe)\n", dir)
	fmt.Println("Once activated, `stowage install`/`uninstall`/`list` act on the environment.")

	return nil
}
