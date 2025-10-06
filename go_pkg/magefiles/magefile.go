//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/linebot-liff-template/go_pkg/mage"
	"github.com/magefile/mage/sh"
)

// Generate runs mockery to generate mocks
func Generate() error {
	fmt.Println("Generating mocks...")
	return mage.RunMockery()
}

// Test runs all tests
func Test() error {
	fmt.Println("Running tests...")
	return sh.RunV("go", "test", "./...", "-v", "-cover")
}

// Fmt formats the code
func Fmt() error {
	fmt.Println("Formatting code...")
	return sh.RunV("go", "fmt", "./...")
}

// Lint runs golangci-lint
func Lint() error {
	fmt.Println("Linting code...")
	return sh.RunV("golangci-lint", "run", "./...")
}

// Update updates dependencies
func Update() error {
	fmt.Println("Updating dependencies...")
	if err := sh.RunV("go", "get", "-u", "./..."); err != nil {
		return err
	}
	return sh.RunV("go", "mod", "tidy")
}

// Clean removes build artifacts
func Clean() error {
	fmt.Println("Cleaning...")
	return sh.Rm("build")
}
