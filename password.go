package common

import (
	. "app/internal/utils"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

// Pwhash (password string) prints the bcrypt hash of the specified argument.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Pwhash(password string) {
	hash := Ensure(bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost))
	Assert(fmt.Fprintln(os.Stdout, string(hash)))
}
