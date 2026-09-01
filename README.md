
[![Release](https://img.shields.io/github/release/alexhudson/stowage.svg?style=flat-square)](https://github.com/alexhudson/stowage/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexhudson/stowage?style=flat-square)](https://goreportcard.com/report/github.com/alexhudson/stowage)
[![SayThanks.io](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg?style=flat-square)](https://saythanks.io/to/alexhudson)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

# Synopsis

**stowage** is a very simple package-manager-alike for Docker containers that wrap cli tools. The idea is not to replace package management; this is largely a development-environment convenience to bring together tools in a relatively simple way.

My primary use case is development teams and CI pipelines. By making development tools available as containers, they can be easily shared and re-used by others developers and CI alike - stowage makes this process simpler, which in turn encourages further such tools and sharing to occur.

This is really not a production environment tool; the use case is not so strong there anyway, but the fundamental use of the root account and the various little conveniences are largely inappropriate and could lead to security issues. However, I don't think there's anything wrong with creating tools to query/manipulate the production environment using stowage.

This project adheres to the Contributor Covenant [code of conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. We appreciate your contribution. Please refer to our [contributing guidelines](CONTRIBUTING.md).

# Installation

Linux, MacOS and Windows are all supported.

## General (works on all platforms)

You can download a copy of the binary from the [stowage releases](https://github.com/alexhudson/stowage/releases/latest) page, and copy that somewhere convenient (i.e., `/usr/local/bin`). This is usually enough.

## MacOS 

MacOS users can use homebrew to install `stowage`:

```
$ brew tap alexhudson/stowage
$ brew install stowage
```

This is the more automatic version of the same process, you will end up with the latest release binary in the right place, and indeed it will update as releases happen.

## Linux

You can download stowage as a container and use it to self-install:

```
$ sudo docker pull ealexhudson/stowage
$ sudo docker run ealexhudson/stowage get-started | sudo sh
```

(If you don't like the idea of piping unknown stuff to a sudo shell, good for you! Just examine the output for the actual docker bootstrapping command.)

## Windows

Download `stowage.exe` from the [stowage releases](https://github.com/alexhudson/stowage/releases/latest) page and put it somewhere convenient, e.g. `%LOCALAPPDATA%\stowage\bin`. Add that directory to your `PATH` (once, via System Properties or `setx PATH "%PATH%;%LOCALAPPDATA%\stowage\bin"`).

Wrapped commands work a bit differently here than on Linux/MacOS: instead of generating a shell script per command, `stowage install` creates an NTFS hard link of `stowage.exe` under the wrapped command's name (e.g. `hello-world.exe`) in `%LOCALAPPDATA%\stowage\bin`. When that link is run, it recognizes it wasn't invoked as `stowage` and dispatches straight to `docker` with the right arguments instead of showing the `stowage` CLI. No admin rights are required - everything installs per-user under `%LOCALAPPDATA%\stowage`.

Docker Desktop (or another Docker install with `docker` on `PATH`) is assumed to already be present.

# Getting started

Once you've installed `stowage` and make it available, you can use it to list all the stowed containers you already have. Usually, this list will be empty, although if you used a container to self-install stowage then it will list itself!

```
$ stowage list
stowage
```

On Linux and MacOS, `install`/`uninstall` don't need `sudo` by default: stowage installs into `$HOME/.local/bin` (and keeps its spec/repo state under `$HOME/.local/share/stowage`) for the current user, so just make sure `$HOME/.local/bin` is on your `PATH`. If you run stowage as root (e.g. via `sudo`), it installs system-wide into `/usr/local` instead, matching the traditional behaviour. Either way, you'll still need to run the wrapped commands under `sudo` if your user does not otherwise have permission to talk to the Docker daemon.

Let's try installing the generic Docker `hello-world` image, and see how this works:

```
$ stowage install hello-world
$ stowage list
stowage
hello-world
$ hello-world

Hello from Docker!
This message shows that your installation appears to be working correctly.

[.. etc ..]
```

The install command can take a Docker image name (in which case it tries to choose some simple defaults), or URL - so long as Docker underneath will understand it. So, you can point to custom repositories, etc.

stowage also has a "specfile" system. While the defaults for containers are as good as possible, it's not possible to work out how all containers should be invoked. These specfiles can also be used to install the command with the `install-spec` command. Developers working on containers can embed the specfile within the image; if this is done, stowage will pick it out and use it.

# Motivation

This is primarily a tool for development environments: the point is to be able to quickly and easily distribute tools to fellow workers that will effectively auto-update, and be able to use them within build pipelines (which are probably container builds).

stowage owes a lot of inspiration to GNU Stow, but has a different model as it attempts to be semi-self hosted - you can run stowage through a container, and then use stowage to manage stowage.

Some better examples of using `stowage`:

1. Installing CLI tools for an environment

To pick an example at random, we can do:

```
$ sudo stowage install --command azure microsoft/azure-cli
$ sudo azure 
info:             _    _____   _ ___ ___
info:            /_\  |_  / | | | _ \ __|
info:      _ ___/ _ \__/ /| |_| |   / _|___ _ _
info:    (___  /_/ \_\/___|\___/|_|_\___| _____)
info:       (_______ _ _)         _ ______ _)_ _ 
info:              (______________ _ )   (___ _ _)
info:    
info:    Microsoft Azure: Microsoft's Cloud Platform
info:    
info:    Tool version 0.10.9
[ .. etc .. ]
```

We now have a "local binary" `azure` that can be used quite conveniently. A similar pattern exists for other API-using CLI tools, such as AWS, OpenStack, to name a few. 

# Making a CLI tool installable via stowage

If you have a CLI tool wrapped in a container, you're already most of the way there - stowage will try to pick some simple defaults. 

### Tips for your Dockerfile

1. Use ENTRYPOINT and CMD together

Point `ENTRYPOINT` at your wrapped executable, and give `CMD` some reasonable default - for `stowage`, I picked `-h` so that when you run it without arguments it gives you the help text. This is a pretty reasonable convention for containers that want to behave like statically linked binaries.

# Building stowage from source

There are two versions of stowage; a proof-of-concept script written in Python that served as the basis of the project, and the new utility written in Go. There is no build required on the Python side, only the Go.

Check out the source anywhere you like - stowage uses Go modules, so it does not need to live under `$GOPATH`. You can then run `make` to see the various targets available.

stowage requires `golang` to be installed. Fetch the required dependencies by running:

```
$ make setup
```

You can then build stowage simply by running:

```
$ make build
```

It is created in the top-level directory; test your new utility by running `./stowage -h` (or `.\stowage.exe -h` on Windows) and verify that the help is output.

## Installation after you built from source

NOTE this is not the general installation instruction; this is just if you have modified stowage - please see above for the regular installation instructions!

There are three different options for installation:

### Install the stowage utility into $GOPATH/bin

This is achieved simply by running:

```
$ make install
```

`stowage` will then be available for your user (assuming you have `$GOPATH/bin` in your `$PATH`) and you will be able to use it as normal.

### Install the stowage utility into a system-wide location.

There is no Makefile command to do this, but you can simply run:

```
$ sudo cp ./stowage /usr/local/bin
```

This will make the `stowage` utility available for all the users on your system without having to mess about with `$PATH`.

### Create the container and stow it locally

Since stowage can manage utilities that are packaged into containers, and stowage itself is such a utility, stowage can self-install itself. To do this, you first need to build a local version of the container image (otherwise, the self-install will pull in the published image, which will not contain the version of the utility you have already build). Then, simply self-install it as per the normal instructions.

One caveat: the container is designed to run in a Linux/amd64 container. If you have built the utility on another platform, such as MacOS, you will get an error about compatibility if you try to run the container, like `standard_init_linux.go:178: exec user process caused "exec format error"`.

To solve this, you need to do a cross-compiled build first. That's what the `make build` below does - it's fine to skip that step if you are already on Linux/amd64.

```
$ GOOS=linux GOARCH=amd64 make build
$ make container
$ docker run ealexhudson/stowage get-started | sudo sh
```

# Miscelleous matters

## Security

There's no clear security model for this tool yet, though it no longer relies on root access by default: stowage installs per-user into `$HOME/.local` (or `%LOCALAPPDATA%\stowage` on Windows) unless it's run as root, in which case it falls back to the traditional system-wide `/usr/local`. Particularly on systems where access to docker is based on group membership or SELinux labels, there's no massive reason for this to be root.

## Tests

Run `make test` (or `go test ./cmd/stowage`) to run the unit tests.

## Contributors

Please, feel free to open pull requests - I welcome all courteous developers, and any help (raising issues, suggesting ideas, etc.) makes me happy.

## License

This is licensed under the MIT license.
