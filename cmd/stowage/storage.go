package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type storage struct {
	specDir string
	repoDir string
}

func createStorage() storage {
	prefix := "/usr/local"
	if _, err := os.Stat("/stowage/install-tree"); !os.IsNotExist(err) {
		prefix = "/stowage/install-tree"
	}

	s := storage{
		specDir: filepath.Join(prefix, "share/stowage/spec/"),
		repoDir: filepath.Join(prefix, "share/stowage/repo/"),
	}
	return s
}

func (s *storage) getRepofilePath(repo *Repository) string {
	return s.getRepofilePathByName(repo.Name)
}

func (s *storage) getRepofilePathByName(name string) string {
	return filepath.Join(s.repoDir, name)
}

func (s *storage) getSpecfilePath(spec *Specification) string {
	return s.getSpecfilePathByName(spec.Name)
}

func (s *storage) getSpecfilePathByName(name string) string {
	return filepath.Join(s.specDir, name)
}

func (s *storage) saveSpecification(spec *Specification) error {
	// ensure we have a directory to install specs into
	if _, err := os.Stat(s.specDir); os.IsNotExist(err) {
		err := os.MkdirAll(s.specDir, 0755)
		if err != nil {
			fmt.Printf("ERROR: Specdir doesn't exist and can't be created. Run as root / sudo?\n")
			return nil
		}
	}

	specFilePath := s.getSpecfilePath(spec)
	// remove any existing spec
	if _, err := os.Stat(specFilePath); os.IsExist(err) {
		os.Remove(specFilePath)
	}

	// install new specfile
	err := ioutil.WriteFile(specFilePath, []byte(spec.toJSON()), 0644)
	if err != nil {
		panic(err)
	}

	return nil
}

func (s *storage) loadSpecification(path string) (Specification, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return Specification{}, err
	}
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return Specification{}, err
	}

	spec := Specification{}
	spec.fromJSON(data)

	return spec, nil
}

func (s *storage) loadSpecificationByName(name string) (Specification, error) {
	specFilePath := s.getSpecfilePathByName(name)
	return s.loadSpecification(specFilePath)
}

func (s *storage) removeSpecification(spec *Specification) error {
	specFilePath := s.getSpecfilePath(spec)

	if err := os.Remove(specFilePath); os.IsExist(err) {
		fmt.Println("Couldn't remove specfile!")
	}

	return nil
}

func (s *storage) listSpecifications() []string {
	files, _ := ioutil.ReadDir(s.specDir)
	result := make([]string, 0)
	for _, f := range files {
		result = append(result, f.Name())
	}
	return result
}
