package common

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

// Pwhash (password string) prints the bcrypt hash of the specified argument.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Pwhash(password string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	_, _ = fmt.Fprintln(os.Stdout, string(hash))
}
