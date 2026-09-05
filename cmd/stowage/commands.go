package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

var selfSpec = Specification{
	Name:    "stowage",
	Image:   "ealexhudson/stowage",
	Command: "",
	Options: runtimeOptions{
		Interactive: false,
		Tty:         true,
		Privileged:  true,
	},
	Mounts: []runtimeMount{
		{
			Host:  "/usr/local",
			Guest: "/stowage/install-tree",
		},
	},
}

func cmdGetStarted(c *cli.Context) error {

	fmt.Println(selfSpec.runCommand() + " self-install")
	return nil
}

func cmdInstall(c *cli.Context) error {
	installer := Installer{
		Request: c.Args().First(),
	}

	if !installer.setup() {
		fmt.Printf("ERROR: Don't know how to install %s\n", installer.Request)
		return nil
	}

	command := c.String("command")
	if command != "" {
		installer.setCommandName(command)
	}

	entrypoint := c.String("entrypoint")
	if entrypoint != "" {
		installer.setEntrypoint(entrypoint)
	}

	_ = installer.run()

	return nil
}

// install a command from a spec file
func cmdInstallSpec(c *cli.Context) error {
	installer := Installer{
		RequestSpec: c.Args().First(),
	}

	if !installer.setup() {
		fmt.Printf("ERROR: Don't know how to install %s\n", installer.RequestSpec)
		return nil
	}

	_ = installer.run()

	return nil
}

func cmdSelfInstall(c *cli.Context) error {
	store := createStorage()

	store.saveSpecification(&selfSpec)

	binary := Binary{name: "stowage", spec: &selfSpec, BinDir: globalBinDir()}
	binary.install()

	return nil
}

func cmdRemove(c *cli.Context) error {
	name := c.Args().First()

	if env, ok := activeVenv(); ok {
		if err := env.remove(name); err != nil {
			fmt.Printf("ERROR: %s\n", err)
			return nil
		}
		updateVenvManifest(func(m *venvManifest) { m.removeByName(name) })
		fmt.Printf("%s removed\n", name)
		return nil
	}

	store := createStorage()

	binary := Binary{name: name, BinDir: globalBinDir()}
	binary.uninstall()

	spec := Specification{Name: name}
	store.removeSpecification(&spec)

	return nil
}

func cmdList(c *cli.Context) error {
	var store storage
	if env, ok := activeVenv(); ok {
		store = env.storage()
	} else {
		store = createStorage()
	}

	for _, spec := range store.listSpecifications() {
		fmt.Println(spec)
	}

	return nil
}

func cmdRun(c *cli.Context) error {
	store := createStorage()

	spec, err := store.loadSpecificationByName(c.Args().First())
	if err != nil {
		fmt.Println("ERROR: no such command installed")
		return nil
	}
	args := spec.runCommandSlice()

	furtherArgs := c.Args()[1:]
	if len(furtherArgs) > 0 && furtherArgs[0] == "--" {
		furtherArgs = furtherArgs[1:]
	}
	for _, arg := range furtherArgs {
		args = append(args, arg)
	}

	binary, lookErr := exec.LookPath("docker")
	if lookErr != nil {
		panic(lookErr)
	}

	os.Exit(runCommand(binary, args[1:]))

	return nil
}

// runCommand runs name with args, connecting it to the parent's stdio, and
// returns its exit code.
func runCommand(name string, args []string) int {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err == nil {
		return 0
	}

	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}

	panic(err)
}
