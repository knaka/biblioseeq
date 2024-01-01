package common

import (
	. "app/internal/utils"
	"github.com/magefile/mage/mg"
	"os"
)

// noinspection GoUnusedExportedType, GoUnnecessarilyExportedIdentifiers
type Client mg.Namespace

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
			wd := Ensure(os.Getwd())
			Assert(os.Chdir(dir))
			defer (func() { Ignore(os.Chdir(wd)) })()
			return RunWith(nil, "npm", "run", "build:development")
		})()
		if err != nil {
			return err
		}
	}
	return nil
}
