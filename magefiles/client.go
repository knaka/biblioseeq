package main

import (
	"github.com/magefile/mage/mg"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/knaka/go-utils"
)

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

// Build builds web client.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Client) Build() (err error) {
	return runClientTask("build")
}
