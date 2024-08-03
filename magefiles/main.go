// Mage is a make/rake-like build tool using Go. You write plain-old go functions, and Mage automatically uses them as Makefile-like runnable targets.
package main

import (
	gobin "github.com/knaka/gobin/lib"
	"github.com/magefile/mage/mg"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/knaka/go-utils"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type PB mg.Namespace

// Gen generates protocol buffer binding code.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (PB) Gen() (err error) {
	return gobin.RunWith([]string{"go-generate-fast"}, gobin.WithDir("pb"))
}

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type DB mg.Namespace

func (DB) Gen() (err error) {
	return gobin.RunWith([]string{"go-generate-fast"}, gobin.WithDir("db"))
}

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Client mg.Namespace

func runClientTask(tasks ...string) (err error) {
	dir := filepath.Join(V(os.Getwd()), "client")
	cmd := exec.Command(filepath.Join(dir, "task"), tasks...)
	cmd.Dir = dir
	cmd.Stdout = Ternary(mg.Verbose(), os.Stdout, nil)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (Client) Build() (err error) {
	return runClientTask("build")
}

// Gen generates code.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Gen() (err error) {
	mg.Deps(
		PB.Gen,
		DB.Gen,
	)
	return gobin.Run("go-generate-fast")
}
