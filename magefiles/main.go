// Mage is a make/rake-like build tool using Go. You write plain-old go functions, and Mage automatically uses them as Makefile-like runnable targets.
package main

import (
	. "github.com/knaka/go-utils"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/target"
)

type Pb mg.Namespace

func Gobin(name string) string {
}

// Gen generates protocol buffer binding code.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Pb) Gen() (err error) {
	V0(Gobin("protoc-gen-connect-go"))
	V0(Gobin("protoc-gen-go"))
	source := V(target.NewestModTime("proto"))
	dest := V(target.NewestModTime("pbgen"))
	if dest.Compare(source) > 0 {
		return nil
	}
	return RunWith("", nil, "buf", "generate", "proto/")
}
