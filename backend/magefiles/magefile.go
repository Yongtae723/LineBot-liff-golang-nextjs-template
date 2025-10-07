//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/sh"

	//mage:import
	_ "cookforyou.com/linebot-liff-template/common/mage/tasks"
)

// Run runs the backend server
func Run() error {
	fmt.Println("Starting backend server...")
	os.Setenv("ENV", "local")
	return sh.RunV("go", "run", "cmd/main.go")
}
