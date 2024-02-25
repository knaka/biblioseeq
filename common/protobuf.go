package common

import (
	. "github.com/knaka/go-utils"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/target"
)

type Pb mg.Namespace

// Gen generates protocol buffer binding code.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Pb) Gen() (err error) {
	source := V(target.NewestModTime("proto"))
	dest := V(target.NewestModTime("pbgen"))
	if dest.Compare(source) > 0 {
		return nil
	}
	return RunWith("", nil, "buf", "generate", "proto/")
}

func init() {
	AddGenFn(Pb.Gen)
}
