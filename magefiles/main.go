// Mage is a make/rake-like build tool using Go. You write plain-old go functions, and Mage automatically uses them as Makefile-like runnable targets.
package main

import (
	. "github.com/knaka/go-utils"
	gobin "github.com/knaka/gobin/lib"
	"github.com/magefile/mage/mg"
	"os"
)

func mgRunWith(dir string, args ...string) error {
	return gobin.RunEx(args,
		gobin.WithDir(dir),
		gobin.WithStdout(Ternary(mg.Verbose(), os.Stdout, nil)),
		gobin.WithStderr(os.Stderr),
	)
}

// Gen generates code.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Gen() (err error) {
	mg.Deps(
		PB.Gen,
		DB.Gen,
	)
	return mgRunWith(".", "go-generate-fast", ".")
}
