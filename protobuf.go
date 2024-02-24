package common

import (
	. "github.com/knaka/go-utils"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/target"
)

// Bufgen generates protocol buffer binding code.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Bufgen() error {
	source := V(target.NewestModTime("proto"))
	dest := V(target.NewestModTime("pbgen"))
	if dest.Compare(source) > 0 {
		return nil
	}
	return RunWith("", nil, "buf", "generate", "proto/")
}

// Pbgen is an alias of 'bufgen'.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Pbgen() { mg.Deps(Bufgen) }
