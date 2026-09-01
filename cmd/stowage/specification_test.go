package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func errRoundTrip(t *testing.T, field string) {
	t.Errorf("JSON roundtrip error: %s didn't roundtrip", field)
}

func TestSpec(t *testing.T) {
	spec := Specification{
		Name:        "test1",
		Environment: []string{"TEST"},
	}

	assert.Equal(t, spec.Name, "test1", "Name should match on spec")

	json := spec.toJSON()
	spec2 := Specification{}
	spec2.fromJSON(json)

	assert.Equal(t, spec.Name, spec2.Name, "Name did not round-trip")
	assert.Equal(t, spec.Environment[0], spec2.Environment[0], "Environment did not round-trip")
}

func TestSpecOptionsAndMountsRoundTrip(t *testing.T) {
	spec := Specification{
		Name:  "test2",
		Image: "example/image",
		Options: runtimeOptions{
			Tty:         true,
			Interactive: true,
			Privileged:  true,
			Readonly:    true,
		},
		Mounts: []runtimeMount{
			{Host: "/host/path", Guest: "/guest/path"},
			{Guest: "/root", Cwd: true},
		},
	}

	spec2 := Specification{}
	spec2.fromJSON(spec.toJSON())

	assert.Equal(t, spec.Options, spec2.Options, "Options did not round-trip")
	assert.Equal(t, spec.Mounts, spec2.Mounts, "Mounts did not round-trip")
}

func TestSpecRun(t *testing.T) {
	spec := Specification{
		Name:        "test1",
		Environment: []string{"TEST"},
	}

	cmd := spec.runCommandSlice()
	assert.Equal(t, len(cmd), 6, "Docker run cmd slice wrong size")
}

func TestSpecRunOptionFlags(t *testing.T) {
	spec := Specification{
		Name:  "test-flags",
		Image: "example/image",
		Options: runtimeOptions{
			Tty:         true,
			Interactive: true,
			Privileged:  true,
			Readonly:    true,
		},
	}

	cmd := spec.runCommandSlice()

	assert.Contains(t, cmd, "-t")
	assert.Contains(t, cmd, "-i")
	assert.Contains(t, cmd, "--privileged")
	assert.Contains(t, cmd, "--read-only=true")
}

func TestSpecRunHostMount(t *testing.T) {
	spec := Specification{
		Name:  "test-mount",
		Image: "example/image",
		Mounts: []runtimeMount{
			{Host: "/data", Guest: "/mnt/data"},
		},
	}

	cmd := spec.runCommandSlice()

	assert.Contains(t, cmd, "-v")
	assert.Contains(t, cmd, "/data:/mnt/data")
}

func TestSpecRunCwdMount(t *testing.T) {
	wd, err := os.Getwd()
	assert.NoError(t, err)

	spec := Specification{
		Name:  "test-cwd",
		Image: "example/image",
		Mounts: []runtimeMount{
			{Guest: "/root", Cwd: true},
		},
	}

	cmd := spec.runCommandSlice()

	assert.Contains(t, cmd, "-w=/root")
	assert.Contains(t, cmd, wd+":/root")
}
