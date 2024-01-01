package common

import (
	"app/env"
	"os"
)

// Psql executes psql(1) command on the database defined in .env* file
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Psql() error {
	env.Load()
	dbUrl := os.Getenv("DB_URL")
	err := ExecWith(
		map[string]string{},
		"psql",
		dbUrl,
	)
	if err != nil {
		dbUrl = os.Getenv("ADMIN_DB_URL")
		err = ExecWith(
			map[string]string{},
			"psql",
			dbUrl,
		)
	}
	return err
}
