// Mage is a make/rake-like build tool using Go. You write plain-old go functions, and Mage automatically uses them as Makefile-like runnable targets.
package main

import (
	. "github.com/knaka/go-utils"
	"github.com/knaka/go-utils/fs"
	"github.com/knaka/gobin"
	"github.com/magefile/mage/mg"
	"os"
	"os/exec"
)

func gobinRun(dir string, args ...string) (err error) {
	_, err = gobin.RunEx(args,
		gobin.WithDir(dir),
		gobin.WithStdin(os.Stdin),
		gobin.WithStdout(Ternary(mg.Verbose(), os.Stdout, nil)),
		gobin.WithStderr(os.Stderr),
	)
	return
}

func run(dir string, args ...string) (err error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = V(fs.CanonPath(dir))
	cmd.Stdin = os.Stdin
	cmd.Stdout = Ternary(mg.Verbose(), os.Stdout, nil)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Gen generates code.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Gen() (err error) {
	mg.Deps(
		PB.Gen,
		DB.Gen,
	)
	return gobinRun(".", "go-generate-fast", ".")
}
