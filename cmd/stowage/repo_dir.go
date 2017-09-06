package main

import (
	"io/ioutil"
	"path/filepath"
)

// RepositoryDir represents a local directory we want to create a repo from
type RepositoryDir struct {
	Path string
	Repo Repository
}

func (rd *RepositoryDir) getRepositoryFilePath() string {
	return filepath.Join(rd.Path, "stowage.json")
}

func (rd *RepositoryDir) getRepository() Repository {
	store := createStorage()
	repo, _ := store.loadRepository(rd.getRepositoryFilePath())

	rd.Repo = repo

	return repo
}

func (rd *RepositoryDir) scan() {
	store := createStorage()

	files, _ := ioutil.ReadDir(rd.Path)

	for _, f := range files {
		if f.Name() != "stowage.json" {
			spec, err := store.loadSpecification(filepath.Join(rd.Path, f.Name()))
			if err == nil {
				rd.Repo.addSpecification(&spec)
			}
		}
	}
}

func (rd *RepositoryDir) save() {
	store := createStorage()

	store.saveRepository(&rd.Repo, rd.getRepositoryFilePath())
}
