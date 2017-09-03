package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Installer is used to figure out what is being requested and how to install
type Installer struct {
	Request string
	Spec    *Specification
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
	name := i.Request

	if strings.Index(name, "://") > -1 {
		// this is a URL
		return i.loadSpecFromURL(name)
	}

	repoSep := strings.Index(name, ":")
	if repoSep > -1 {
		// this could be a repo reference
		repo := name[0:repoSep]
		name = name[repoSep+1:]

		if strings.Index(repo, "/") == -1 {
			// repo names cannot have slashes in them; this must be a
			// docker hub reference!
			return i.loadSpecFromRepo(repo, name)
		}
	}

	_ = i.loadSpecFromFie(name)

	if i.Spec == nil {
		spec := Specification{
			Name:    name,
			Image:   name,
			Command: "",
		}
		i.Spec = &spec
	}

	return true
}

func (i *Installer) loadSpecFromFie(path string) bool {
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
	fmt.Println(response.Body)
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

func (i *Installer) run() bool {
	store := createStorage()
	store.saveSpecification(i.Spec)

	binary := Binary{name: i.Spec.Name, spec: i.Spec}
	binary.install()

	fmt.Printf("%s installed\n", i.Spec.Name)
	return true
}
