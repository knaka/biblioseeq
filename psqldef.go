package common

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	. "github.com/knaka/go-utils"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
var BeforeQueryFiles = []string{
	filepath.Join("db", "schema_before.sql"),
}

//goland:noinspection GoUnnecessarilyExportedIdentifiers
var AfterQueryFiles = []string{
	filepath.Join("db", "schema_after.sql"),
}

// Converge converges a PostgreSQL database.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Converge() (err error) {
	defer Catch(&err)
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return errors.New("DB_URL is not set")
	}
	u, err := url.Parse(dbUrl)
	if err != nil {
		return err
	}
	// read db/schema_before.sql
	for _, f := range BeforeQueryFiles {
		V0(execDbQuery(string(V(os.ReadFile(f)))))
	}
	var args []string
	pass, _ := u.User.Password()
	dbName := strings.ReplaceAll(u.Path, "/", "")
	V0(RunWith(
		"",
		map[string]string{
			"PGPASSWORD": pass,
		},
		"psqldef",
		append([]string{
			"--file", "db/schema.sql",
			"--file", "db/schema_info.sql",
			"--host", u.Hostname(),
			"--port", u.Port(),
			"--user", u.User.Username(),
			"--skip-view", // view management is not reliable
			dbName,
		}, args...)...,
	))
	for _, f := range AfterQueryFiles {
		V0(execDbQuery(string(V(os.ReadFile(f)))))
	}
	return
}

// Dump dumps a database.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Dump() error {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return errors.New("DB_URL is not set")
	}
	u, err := url.Parse(dbUrl)
	if err != nil {
		return err
	}
	var args []string
	dbName := strings.ReplaceAll(u.Path, "/", "")
	return RunWith(
		"",
		nil,
		"psqldef",
		append([]string{
			"--file", "db/schema.sql",
			"--export",
			"--host", u.Hostname(),
			"--port", u.Port(),
			"--user", u.User.Username(),
			dbName,
		}, args...)...,
	)
}

// Gen generates bindings and helpers for a database.
func (Db) Gen() (err error) {
	if err = Gobin("go-generate-fast"); err != nil {
		return
	}
	return RunWith("db", nil, "go-generate-fast", ".")
}
