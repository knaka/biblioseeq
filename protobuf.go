package common

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/target"
)

// Bufgen generates protocol buffer binding code.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Bufgen() error {
	source, _ := target.NewestModTime("proto")
	dest, _ := target.NewestModTime("pbgen")
	if dest.Compare(source) > 0 {
		return nil
	}
	return ExecWith(nil, "buf", "generate", "proto/")
}

// Pbgen is an alias of 'bufgen'.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Pbgen() {
	mg.Deps(Bufgen)
}
