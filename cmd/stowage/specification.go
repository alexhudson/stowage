package main

import (
	"encoding/json"
	"strings"
)

type runtimeOptions struct {
	Tty         bool
	Interactive bool
	Privileged  bool
	Readonly    bool
}

// TODO: options for mount (selinux, etc..)
type runtimeMount struct {
	Host  string // path on the docker host
	Guest string // destination mount on the container
	Cwd   bool   // whether we're asking to mount current working directory
}

// Specification is a type
type Specification struct {
	Name        string
	Image       string
	Command     string
	Environment []string

	Options runtimeOptions
	Mounts  []runtimeMount
}

func (s *Specification) fromJSON(byt []byte) error {
	if err := json.Unmarshal(byt, &s); err != nil {
		panic(err)
	}
	return nil
}

func (s *Specification) toJSON() []byte {
	ret, _ := json.Marshal(s)
	return ret
}

func (s *Specification) getImage() string {
	return s.Image
}

func (s *Specification) getCommand() string {
	return s.Command
}

func (s *Specification) runCommandSlice() []string {
	cmd := []string{"docker", "run", "--rm"}

	if s.Options.Tty {
		cmd = append(cmd, "-t")
	}
	if s.Options.Interactive {
		cmd = append(cmd, "-i")
	}
	if s.Options.Privileged {
		cmd = append(cmd, "--privileged")
	}
	if s.Options.Readonly {
		cmd = append(cmd, "--read-only=true")
	}

	if s.Mounts != nil {
		for _, mount := range s.Mounts {
			source := mount.Host
			if mount.Cwd {
				source = "`pwd`"
				cmd = append(cmd, "-w="+mount.Guest)
			}
			cmd = append(cmd, "-v", source+":"+mount.Guest)
		}
	}

	for _, variable := range s.Environment {
		cmd = append(cmd, "-e", variable)
	}

	cmd = append(cmd, s.getImage())
	if s.getCommand() != "" {
		cmd = append(cmd, s.getCommand())
	}
	// TODO: args
	return cmd
}

func (s *Specification) runCommand() string {
	return strings.Join(s.runCommandSlice(), " ")
}
