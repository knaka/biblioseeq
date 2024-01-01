package common

import (
	"os"
)

// Psql executes psql(1) command on the database defined in .env* file
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Psql() error {
	var err error
	for _, dbUrl := range []string{
		os.Getenv("DB_URL"),
		os.Getenv("ADMIN_DB_URL"),
	} {
		err = RunWith(nil, "psql", dbUrl)
		if err == nil {
			return nil
		}
	}
	return err
}
