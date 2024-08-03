// Mage is a make/rake-like build tool using Go. You write plain-old go functions, and Mage automatically uses them as Makefile-like runnable targets.
package main

import (
	gobin "github.com/knaka/gobin/lib"
	"github.com/magefile/mage/mg"
)

//goland:noinspection GoUnusedExportedType, GoUnnecessarilyExportedIdentifiers
type Pb mg.Namespace

// Gen generates protocol buffer binding code.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Pb) Gen() (err error) {
	//V0(gobin.Install("protoc-gen-connect-go", "protoc-gen-go"))
	//if V(target.NewestModTime("pbgen")).Compare(
	//	V(target.NewestModTime("proto"))) > 0 {
	//	return
	//}
	//return gobin.Run("buf", "generate", "proto")
	return gobin.RunWith([]string{"go-generate-fast"}, gobin.WithDir("pb"))
}

// Gen generates code.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Gen() (err error) {
	return gobin.Run("go-generate-fast")
}
