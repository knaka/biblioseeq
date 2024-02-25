package psql

import (
	"fmt"
	common "github.com/knaka/magefiles-common"
	"github.com/magefile/mage/sh"
	"log"
	"os"
	"os/exec"
)

func execGooseWoDump(args ...string) error {
	return execGooseSub(false, args...)
}

func execGoose(args ...string) error {
	return execGooseSub(true, args...)
}

var LoadEnv = func() {}

func execGooseSub(dump bool, args ...string) error {
	LoadEnv()
	dbUrl := os.Getenv("DB_URL")
	err := common.RunWith(
		"",
		map[string]string{
			"GOOSE_DRIVER":   "postgres",
			"GOOSE_DBSTRING": dbUrl,
		},
		"goose",
		append([]string{
			//"-table", "internal.goose_db_version",
			"-dir", "db/migrations/",
		}, args...)...,
	)
	if err != nil {
		return err
	}
	if dump {
		_, err = exec.LookPath("pg_dump")
		if err != nil {
			log.Panicf("panic 3950e7a (%v)", err)
		}
		output, err := os.Create("current_schema.sql")
		if err != nil {
			log.Panicf("panic 51fbd85 (%v)", err)
		}
		_, err = sh.Exec(nil, output, os.Stderr, "pg_dump", "--schema-only", "--schema", "public", "--schema", "blob", dbUrl)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to dump schema\n")
			//_, _ = fmt.Fprintf(os.Stderr, "err: %v\n", err)
		}
	}
	// todo: fallback to Docker Compose // docker compose exec db /bin/bash -c 'PGPASSWORD="$POSTGRES_PASSWORD" pg_dump --schema-only --schema public --schema blob -U $POSTGRES_USER hoge'
	return nil
}

// Reset rolls back all migrations.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Reset() error {
	return execGoose("reset")
}

// Up migrates the database up by 1.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Up() error {
	return execGoose("up-by-one")
}

// Down rolls back the database version by 1.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Down() error {
	return execGoose("down")
}

// New (name string) creates a new migration file with the current timestamp.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) New(name string) error {
	return execGooseWoDump("create", name, "sql")
}

// Fix applies sequential ordering to migrations.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Fix() error {
	return execGooseWoDump("fix")
}
