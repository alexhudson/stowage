## Synopsis

**stowage** is a very simple package-manager-alike for Docker containers that wrap cli tools. The idea is not to replace package management; this is largely a development-environment convenience to bring together tools in a relatively simple way.

This is absolutely not a production environment tool; the use case is not so strong there anyway, but the fundamental use of the root account and the various little conveniences are largely inappropriate and could lead to security issues.

## Example

You can download stowage as a container and use it to self-install:

  $ sudo docker pull ealexhudson/stowage
  $ sudo docker run --rm -ti --privileged --read-only=true -v /:/stowage-managed-system ealexhudson/stowage self-install

After that, stowage should be available on the system - but nothing will yet be installed:

  $ sudo stowage list
  $ sudo stowage install hello-world
  $ sudo hello-world

  Hello from Docker!
  This message shows that your installation appears to be working correctly.

  [.. etc ..]

The install command can take a Docker image name (in which case it tries to choose some simple defaults), or a reference to a local JSON file (for more complex commands - there are some examples in the examples folder). In the future I also want to add an ability to refer to JSON file by URL.

## Motivation

This is primarily a tool for development environments: the point is to be able to quickly and easily distribute tools to fellow workers that will effectively auto-update, and be able to use them within build pipelines (which are probably container builds).

stowage owes a lot of inspiration to GNU Stow, but has a different model as it attempts to be semi-self hosted - you can run stowage through a container, and then use stowage to manage stowage.

## Making a CLI tool installable via stowage

If you have a CLI tool wrapped in a container, you're already most of the way there - stowage will try to pick some simple defaults. 

## Developing

stowage can be run locally as a single script rather than through the container; this is the easiest method for development.

In the future I'd likely convert this to another language which produces a single static binary (e.g. golang) - but for now, it's remaining in Python because it's easy to test a few different ideas rapidly.

## Security

There's no clear security model for this tool yet - it relies a lot on the premise of root access. In the future, I'd like to significantly limit this as far as possible. Particularly on systems where access to docker is based on group membership or SELinux labels, there's no massive reason for this to be root.

Also, the need for stowage to access the root filesystem is mainly for it to install the various docker wrapper scripts that it creates - the more security conscious are probably willing to adjust $PATH instead, and that should be a supported option.


## Tests

Need to write some :-)

## Contributors

Please, feel free to open pull requests - I welcome all courteous developers, and any help (raising issues, suggesting ideas, etc.) makes me happy.

## License

This is licensed under the MIT license.