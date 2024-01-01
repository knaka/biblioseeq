package common

import (
	"net/url"
	"os"
	"strings"
)

// PGPASSWORD= ./.gobin/psqldef --file=db/schema.sql --host=127.0.0.1 --port=5432 --user=knaka --dry-run yuber_dev
// PGPASSWORD= ./.gobin/psqldef --file=db/schema.sql --host=127.0.0.1 --port=5432 --user=knaka yuber_dev

// Converge converges a PostgreSQL database.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Converge() error {
	//todo: どうしよう
	//env.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return nil
	}
	u, err := url.Parse(dbUrl)
	if err != nil {
		return err
	}
	args := []string{}
	pass, _ := u.User.Password()
	dbName := strings.ReplaceAll(u.Path, "/", "")
	err = ExecWith(
		map[string]string{
			"PGPASSWORD": pass,
		},
		"psqldef",
		append([]string{
			"--file", "db/schema.sql",
			"--host", u.Hostname(),
			"--port", u.Port(),
			"--user", u.User.Username(),
			dbName,
		}, args...)...,
	)
	return err
}

// Dump dumps a database.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Dump() error {
	//todo: どうしよう
	//env.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return nil
	}
	u, err := url.Parse(dbUrl)
	if err != nil {
		return err
	}
	args := []string{}
	dbName := strings.ReplaceAll(u.Path, "/", "")
	err = ExecWith(
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
	return err
}
