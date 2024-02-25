package psqldb

import (
	"errors"
	common "github.com/knaka/magefiles-shared/common"
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
	u := V(url.Parse(dbUrl))
	for _, f := range BeforeQueryFiles {
		V0(execDbQuery(string(V(os.ReadFile(f)))))
	}
	var args []string
	pass, _ := u.User.Password()
	dbName := strings.ReplaceAll(u.Path, "/", "")
	V0(common.RunWith(
		"",
		map[string]string{
			"PGPASSWORD": pass,
		},
		"psqldef",
		append([]string{
			"--file", filepath.Join("db", "schema.sql"),
			"--file", filepath.Join("db", "schema_info.sql"),
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
	u := V(url.Parse(dbUrl))
	var args []string
	dbName := strings.ReplaceAll(u.Path, "/", "")
	return common.RunWith(
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
	if err = common.Gobin("go-generate-fast"); err != nil {
		return
	}
	return common.RunWith("db", nil, "go-generate-fast", ".")
}

func init() {
	common.AddGenFn(Db.Gen)
}
