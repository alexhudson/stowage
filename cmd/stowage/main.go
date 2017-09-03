package main

import (
	"os"

	"github.com/urfave/cli"
)

var version = "dev"

func main() {
	app := cli.NewApp()
	app.Name = "stowage"
	app.Version = version
	app.Author = "Alex Hudson <ealexhudson@gmail.com>"
	app.Usage = "A simple package manager-alike for Docker 'binaries'"

	app.Commands = []cli.Command{
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search for a wrapper via the installed repositories",
			Action:  cmdSearch,
		},
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "Install a wrapper for a container",
			Action:  cmdInstall,
		},
		{
			Name:    "uninstall",
			Aliases: []string{"u"},
			Usage:   "Uninstall a wrapper for a container",
			Action:  cmdRemove,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all installed wrappers",
			Action:  cmdList,
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run a command directly",
			Action:  cmdRun,
		},
		{
			Name:    "repo-add",
			Aliases: []string{"ra"},
			Usage:   "Add a repository to the system by URL",
			Action:  cmdRepoAdd,
		},
		{
			Name:    "repo-list",
			Aliases: []string{"rl"},
			Usage:   "List known repositories",
			Action:  cmdRepoList,
		},
		{
			Name:    "repo-scan",
			Aliases: []string{"rs"},
			Usage:   "Scan a directory and create a repository file",
			Action:  cmdRepoScan,
		},
		{
			Name:    "get-started",
			Aliases: []string{"gi"},
			Usage:   "Output a shell script to self-install stowage",
			Action:  cmdGetStarted,
		},
		{
			Name:    "self-install",
			Aliases: []string{"si"},
			Usage:   "Install a wrapper for stowage itself",
			Action:  cmdSelfInstall,
		},
	}

	_ = app.Run(os.Args)
}
