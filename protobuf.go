package common

import (
	. "github.com/knaka/go-utils"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/target"
)

// Bufgen generates protocol buffer binding code.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Bufgen() error {
	source := Ensure(target.NewestModTime("proto"))
	dest := Ensure(target.NewestModTime("pbgen"))
	if dest.Compare(source) > 0 {
		return nil
	}
	return RunWith(nil, "buf", "generate", "proto/")
}

// Pbgen is an alias of 'bufgen'.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Pbgen() { mg.Deps(Bufgen) }
