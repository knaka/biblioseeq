package common

import (
	"fmt"
	. "github.com/knaka/go-utils"
	"golang.org/x/crypto/bcrypt"
	"os"
)

// Pwhash (password string) prints the bcrypt hash of the specified argument.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Pwhash(password string) {
	hash := V(bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost))
	V0(fmt.Fprintln(os.Stdout, string(hash)))
}
