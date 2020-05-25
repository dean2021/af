package system

import (
	"testing"
)

func TestKernelVersion(t *testing.T) {
	version, err := KernelVersion()
	if err != nil {
		t.Error(err)
	}
	t.Log(version)
}

func TestGetCurrentDirectory(t *testing.T) {
	directory := GetCurrentDirectory()
	t.Log(directory)
}
