package common

import (
	"github.com/magefile/mage/mg"
	"os"
)

// noinspection GoUnusedExportedType, GoUnnecessarilyExportedIdentifiers
type Client mg.Namespace

type bar struct {
	dir string
}

var clientInfo []bar

func InitClient(dirs ...string) {
	for _, dir := range dirs {
		clientInfo = append(clientInfo, bar{dir: dir})
	}
}

// Build builds web client application.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Client) Build() error {
	// Docker build does not build clients.
	if os.Getenv("NO_CLIENT_BUILD") != "" {
		return nil
	}
	var err error
	for _, clt := range clientInfo {
		err = (func() error {
			_ = os.Chdir(clt.dir)
			defer (func() { _ = os.Chdir("..") })()
			return ExecWith(nil, "npm", "run", "build:development")
		})()
		if err != nil {
			return err
		}
	}
	return nil
}
