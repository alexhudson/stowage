package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

// Installer is used to figure out what is being requested and how to install
type Installer struct {
	Request     string
	RequestSpec string
	Spec        *Specification
}

/* Requests can be a number of things:
a URL, e.g. https://example.com/stowage/
a reference to a repo, e.g. myrepo\some-binary
a local file, e.g. ./somedir/some-binary.json
a docker hub reference, e.g. ealexhudson/stowage

If it doesn't look like a URL or repo reference, we try it as a local
file. If that doesn't work, we assume it's a Docker hub reference.

*/

func (i *Installer) setup() bool {
	if i.RequestSpec != "" {
		return i.setupSpec()
	}
	return i.setupImage()
}

func (i *Installer) setupSpec() bool {
	specname := i.RequestSpec

	if strings.Index(specname, "::") > -1 {
		// this is a URL
		return i.loadSpecFromURL(specname)
	}

	_ = i.loadSpecFromFile(specname)

	return i.Spec != nil
}

func (i *Installer) setupImage() bool {
	name := i.Request
	image := i.Request

	// check this isn't in a repo first
	repoSep := strings.Index(name, "::")
	if repoSep > -1 {
		repo := name[0:repoSep]
		specname := name[repoSep+2:]

		return i.loadSpecFromRepo(repo, specname)
	}

	// handle tags or digests - we want to remove this from
	// the name, but keep the reference on the image
	tagSep := strings.Index(name, ":")
	if tagSep > -1 {
		name = name[0:tagSep]
	}
	digestSep := strings.Index(name, "@")
	if digestSep > -1 {
		name = name[0:digestSep]
	}

	// try fetching image; if this fails later things may not work but
	// that's not necessarily fatal
	fetchCmd := exec.Command("docker", "image", "pull", name)
	fetchCmd.Run()

	spec := Specification{
		Name:    name,
		Image:   image,
		Command: "",
	}

	// check if we have a custom label with our specfile.
	specCmd := exec.Command("docker", "inspect", "--format",
		"{{ index .Config.Labels \"org.stowage.spec\" }}",
		name,
	)
	imgSpec, err := specCmd.Output()
	if err == nil && len(imgSpec) > 4 {
		// if a spec file was provided via the image, let's load that up
		spec.fromJSON(imgSpec)
	}

	i.Spec = &spec

	cliName := strings.LastIndex(name, "/")
	if (i.Spec != nil) && (cliName > -1) && (i.Spec.Name == name) {
		// this is a reference to an image in a registry - we want to use
		// the project name as cli name. Only do this if we haven't already
		// learned a better name somewhere (e.g. Docker label)
		i.Spec.Name = name[cliName+1:]
	}

	return true
}

func (i *Installer) loadSpecFromFile(path string) bool {
	store := createStorage()
	spec, err := store.loadSpecification(path)
	if err != nil {
		return false
	}

	i.Spec = &spec
	return true
}

func (i *Installer) loadSpecFromURL(url string) bool {
	var spec Specification

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Specfile missing from repository!")
		return false
	}
	buf, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(buf, &spec)

	i.Spec = &spec
	return true
}

func (i *Installer) loadSpecFromRepo(repoName string, name string) bool {
	store := createStorage()

	repo, err := store.loadRepositoryByName(repoName)
	if err != nil {
		fmt.Println("No such repository.")
		return false
	}

	urlForSpec := repo.getURLForSpec(name)

	return i.loadSpecFromURL(urlForSpec)
}

func (i *Installer) setCommandName(name string) {
	i.Spec.Name = name
}

func (i *Installer) run() bool {
	store := createStorage()
	store.saveSpecification(i.Spec)

	binary := Binary{name: i.Spec.Name, spec: i.Spec}
	binary.install()

	fmt.Printf("%s installed\n", i.Spec.Name)
	return true
}
