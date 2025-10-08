// Package tasks contains common tasks defined by all builds. Build scripts should
// use //mage:import to include these and define any other tasks they may need.
package tasks

import (
	"path/filepath"
	"runtime"

	"cookforyou.com/linebot-liff-template/common/mage"

	"github.com/magefile/mage/sh"
)

var RepositoryRoot = func() string {
	_, curFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(curFile), "..", "..", "..")
}()

// Fmt formats source code.
func Fmt() error {
	if err := sh.RunV("go", "run", "mvdan.cc/gofumpt@"+mage.VerGoFumpt, "-l", "-w", "."); err != nil {
		return err
	}
	if err := sh.RunV("go", "run", "github.com/golangci/golangci-lint/cmd/golangci-lint@"+mage.VerGolangCILint, "run", "--fix", "--timeout=10m"); err != nil {
		return err
	}
	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}

	return nil
}

// Lint runs code lint checks.
func Lint() error {
	return sh.RunV("go", "run", "github.com/golangci/golangci-lint/cmd/golangci-lint@"+mage.VerGolangCILint, "run", "--timeout=10m")
}

// Update updates dependencies.
func Update() error {
	if err := sh.RunV("go", "get", "-u", "./..."); err != nil {
		return err
	}

	if err := sh.RunV("go", "mod", "tidy"); err != nil {
		return err
	}

	return nil
}
