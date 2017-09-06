package main

import (
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

func TestSpecRun(t *testing.T) {
	spec := Specification{
		Name:        "test1",
		Environment: []string{"TEST"},
	}

	cmd := spec.runCommandSlice()
	assert.Equal(t, len(cmd), 6, "Docker run cmd slice wrong size")
}
