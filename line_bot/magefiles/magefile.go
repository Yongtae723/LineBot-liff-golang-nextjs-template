//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"
)

// Run runs the LINE Bot server
func Run() error {
	fmt.Println("Starting LINE Bot server...")
	os.Setenv("ENV", "local")
	return sh.RunV("go", "run", "cmd/main.go")
}
