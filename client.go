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

var ClientDirs []string

// Build builds web client application.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Client) Build() error {
	// Docker build does not build clients.
	if os.Getenv("NO_CLIENT_BUILD") != "" {
		return nil
	}
	var err error
	for _, dir := range ClientDirs {
		err = (func() error {
			_ = os.Chdir(dir)
			defer (func() { _ = os.Chdir("..") })()
			return ExecWith(nil, "npm", "run", "build:development")
		})()
		if err != nil {
			return err
		}
	}
	return nil
}
