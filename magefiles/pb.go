package main

import "github.com/magefile/mage/mg"

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type PB mg.Namespace

// Gen generates protocol buffer binding code.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (PB) Gen() (err error) {
	return mgRunWith("pb", "go-generate-fast", ".")
}
