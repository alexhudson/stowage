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
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "Install a wrapper for a container",
			Action:  cmdInstall,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "command",
					Value: "",
					Usage: "name for the command",
				},
				cli.StringFlag{
					Name:  "entrypoint",
					Value: "",
					Usage: "image entrypoint for the command",
				},
			},
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
			Name:     "search",
			Aliases:  []string{"s"},
			Category: "Repositories",
			Usage:    "Search for a wrapper via the installed repositories",
			Action:   cmdSearch,
		},
		{
			Name:     "repo-add",
			Aliases:  []string{"ra"},
			Category: "Repositories",
			Usage:    "Add a repository to the system by URL",
			Action:   cmdRepoAdd,
		},
		{
			Name:     "repo-list",
			Aliases:  []string{"rl"},
			Category: "Repositories",
			Usage:    "List known repositories",
			Action:   cmdRepoList,
		},
		{
			Name:     "repo-scan",
			Aliases:  []string{"rs"},
			Category: "Repositories",
			Usage:    "Scan a directory and create a repository file",
			Action:   cmdRepoScan,
		},
		{
			Name:    "install-spec",
			Aliases: []string{"is"},
			Usage:   "Install a wrapper for a container from a defined spec file",
			Action:  cmdInstallSpec,
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
