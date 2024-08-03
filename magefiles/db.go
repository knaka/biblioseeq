package main

import "github.com/magefile/mage/mg"

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type DB mg.Namespace

// Gen generates database binding code.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (DB) Gen() (err error) {
	return mgRunWith("db", "go-generate-fast", ".")
}
