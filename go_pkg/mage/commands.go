package mage

import (
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

// This file contains definitions of commands that are often needed in build scripts.

// RunMockery runs mockery, generating mocks based on a .mockery.yaml file.
func RunMockery() error {
	return sh.Run("go", "run", "github.com/vektra/mockery/v2@"+VerMockery)
}

// RunMockgen runs mockgen with the provided arguments.
func RunMockgen(pkg string, src string, dest string, opts ...string) error {
	if changes, err := target.Path(dest, src); err != nil {
		return err
	} else if !changes {
		return nil
	}

	args := []string{
		"run", "github.com/golang/mock/mockgen",
		"-package=" + pkg, "-source=" + src, "-destination=" + dest,
	}
	args = append(args, opts...)
	return sh.RunV("go", args...)
}
