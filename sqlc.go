package common

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

type Sqlc mg.Namespace

// Gen generates sqlc queries code.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Sqlc) Gen() error {
	source, _ := target.NewestModTime(
		filepath.Join("db", "migrations"),
		filepath.Join("db", "queries"),
		filepath.Join("db", "schema.sql"),
	)
	destPath := filepath.Join("db", "sqlcgen")
	dest, _ := target.NewestModTime(destPath)
	if !dest.IsZero() && dest.Compare(source) > 0 {
		return nil
	}
	err := sh.Rm(destPath)
	if err != nil {
		return err
	}
	return ExecWith(nil, "sqlc", "generate")
}

// Vet runs queries through a set of lint rules.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Sqlc) Vet() error {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Panicf("panic d5e4c70")
	}
	u, err := url.ParseRequestURI(dbUrl)
	if err != nil {
		log.Panicf("panic 5d83c26 (%v)", err)
	}
	values := u.Query()
	values.Set("options", "-c enable_seqscan=off")
	u.RawQuery = values.Encode()
	err = os.Setenv("DB_URL", u.String())
	if err != nil {
		log.Panicf("panic 8b901e4 (%v)", err)
	}
	log.Println("1436ea9", u.String())
	// Linting queries — sqlc 1.23.0 documentation https://docs.sqlc.dev/en/stable/howto/vet.html
	return ExecWith(map[string]string{
		//"SQLCDEBUG": "dumpvetenv=1,dumpexplain=1", // Environment variables — sqlc 1.23.0 documentation https://docs.sqlc.dev/en/stable/reference/environment-variables.html#sqlcdebug
	}, "sqlc", "vet")
}
