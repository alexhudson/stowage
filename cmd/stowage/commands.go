package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"syscall"

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
	store := createStorage()
	name := c.Args().First()

	spec, err := store.loadSpecification(name)
	if err != nil {
		spec = Specification{
			Name:    name,
			Image:   name,
			Command: "",
		}
	}
	store.saveSpecification(&spec)

	binary := Binary{name: spec.Name, spec: spec}
	binary.install()

	fmt.Printf("%s installed\n", spec.Name)
	return nil
}

func cmdSelfInstall(c *cli.Context) error {
	store := createStorage()

	store.saveSpecification(&selfSpec)

	binary := Binary{name: "stowage", spec: selfSpec}
	binary.install()

	return nil
}

func cmdRemove(c *cli.Context) error {
	store := createStorage()
	name := c.Args().First()

	binary := Binary{name: name}
	binary.uninstall()

	spec := Specification{Name: name}
	store.removeSpecification(&spec)

	return nil
}

func cmdList(c *cli.Context) error {
	store := createStorage()

	specs := store.listSpecifications()

	for _, spec := range specs {
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

	binary, lookErr := exec.LookPath("docker")
	if lookErr != nil {
		panic(lookErr)
	}

	// TODO - really should clean this!
	env := os.Environ()

	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}

	return nil
}

func cmdRepoAdd(c *cli.Context) error {
	var repo Repository

	uri := c.Args().Get(0)
	response, _ := http.Get(uri + "_stowage.json")
	buf, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(buf, &repo)
	repo.URI = uri

	store := createStorage()
	store.saveRepositoryByName(&repo, repo.Name)

	return nil
}

func cmdRepoScan(c *cli.Context) error {
	repoDir := RepositoryDir{
		Path: c.Args().Get(0),
	}

	repo := repoDir.getRepository()
	if repo.Name == "" {
		repo.Name = c.Args().Get(1)
	}

	repoDir.scan()
	repoDir.save()

	return nil
}

func cmdRepoList(c *cli.Context) error {
	store := createStorage()

	repos := store.listRepositories()

	for _, repo := range repos {
		fmt.Println(repo)
	}

	return nil
}
