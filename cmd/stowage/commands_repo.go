package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/urfave/cli"
)

func cmdRepoAdd(c *cli.Context) error {
	var repo Repository

	uri := c.Args().Get(0)
	response, _ := http.Get(uri + "stowage.json")
	buf, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(buf, &repo)
	repo.URI = uri

	if repo.Name == "" {
		fmt.Println("ERROR: Repository is misconfigured; 'name' mising")
		return nil
	}

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

func cmdSearch(c *cli.Context) error {
	store := createStorage()

	repos := store.listRepositories()
	term := c.Args().Get(0)

	var result []RepositoryEntry
	for _, repoName := range repos {
		repo, _ := store.loadRepositoryByName(repoName)
		result = repo.search(term)

		for _, hit := range result {
			fmt.Println(repo.Name + "::" + hit.Name + "\t" + hit.Description)
		}
	}
	return nil
}
