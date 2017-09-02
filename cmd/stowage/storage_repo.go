package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func (s *storage) loadRepository(path string) (Repository, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Repository{}, err
	}

	repo := Repository{}
	repo.fromJSON(data)

	return repo, nil
}

func (s *storage) loadRepositoryByName(name string) (Repository, error) {
	repoFilePath := s.getRepofilePathByName(name)
	return s.loadRepository(repoFilePath)
}

func (s *storage) saveRepositoryByName(repo *Repository, name string) error {
	// ensure we have a directory to install specs into
	if _, err := os.Stat(s.repoDir); os.IsNotExist(err) {
		err := os.MkdirAll(s.repoDir, 0755)
		if err != nil {
			fmt.Printf("ERROR: Repodir %s doesn't exist and can't be created. Run as root / sudo?\n", s.repoDir)
			return nil
		}
	}

	repoFilePath := s.getRepofilePath(repo)

	return s.saveRepository(repo, repoFilePath)
}

func (s *storage) saveRepository(repo *Repository, path string) error {
	// remove any existing repo
	if _, err := os.Stat(path); os.IsExist(err) {
		os.Remove(path)
	}

	// install new specfile
	err := ioutil.WriteFile(path, []byte(repo.toJSON()), 0644)
	if err != nil {
		panic(err)
	}

	return nil
}

func (s *storage) listRepositories() []string {
	files, _ := ioutil.ReadDir(s.repoDir)
	result := make([]string, 0)
	for _, f := range files {
		result = append(result, f.Name())
	}
	return result
}
