package common

import "github.com/magefile/mage/sh"

// Dockerfiles generates dockerfiles from the template.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Dockerfiles() error {
	return sh.RunWith(nil, "go", "generate", "./gendockerfiles")
}
