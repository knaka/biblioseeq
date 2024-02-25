package common

import (
	. "github.com/knaka/go-utils"
	"github.com/magefile/mage/mg"
	"os"
)

//goland:noinspection GoUnusedExportedType, GoUnnecessarilyExportedIdentifiers
type Client mg.Namespace

var clientDirBases []string

func SetClientDirBases(dirBases ...string) {
	clientDirBases = dirBases
}

// Build builds web client application.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Client) Build() error {
	// Docker build does not build clients.
	if os.Getenv("NO_CLIENT_BUILD") != "" {
		return nil
	}
	var err error
	for _, dirBase := range clientDirBases {
		err = (func() error {
			wd := V(os.Getwd())
			V0(os.Chdir(dirBase))
			defer (func() { V0(os.Chdir(wd)) })()
			return RunWith("", nil, "npm", "run", "build:development")
		})()
		if err != nil {
			return err
		}
	}
	return nil
}
