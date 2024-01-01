package common

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// noinspection GoUnusedExportedType, GoUnnecessarilyExportedIdentifiers
type Db mg.Namespace

func execMainDatabaseDdl(ddlWithDbNamePlaceholder string) error {
	adminDbUrl := os.Getenv("ADMIN_DB_URL")
	dbUrl := os.Getenv("DB_URL")
	urlDb, err := url.Parse(dbUrl)
	if err != nil {
		return err
	}
	dbName := strings.Replace(urlDb.Path, "/", "", 1)
	if adminDbUrl == "" {
		urlDb.Path = "/template1"
		adminDbUrl = urlDb.String()
	}
	db, err := sql.Open("pgx", adminDbUrl)
	if err != nil {
		return err
	}
	defer (func() { _ = db.Close() })()
	_, err = db.Exec(fmt.Sprintf(ddlWithDbNamePlaceholder, dbName))
	if err != nil {
		return err
	}
	return nil
}

// Create creates a new PostgreSQL database.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Create() error {
	err := execMainDatabaseDdl("CREATE DATABASE %s")
	if err != nil {
		return err
	}
	return nil
}

func execDbQuery(query string) error {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return nil
	}
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		return err
	}
	defer (func() { _ = db.Close() })()
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		log.Panicf("panic 27d2d65 (%v)", err)
	}
	log.Println("Rows affected:", n)
	return nil
}

// Drop drops a PostgreSQL database if exists.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Drop() error {
	return execMainDatabaseDdl("DROP DATABASE IF EXISTS %s")
}

// Seed inserts seed data into the database.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Db) Seed() error {
	return sh.Run(mg.GoCmd(), "run", filepath.Join(".", "cmd", "db-seed"))
}

// Migrate migrates the database to the most recent version available.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (db Db) Migrate() error {
	err := execGoose("up")
	if err != nil {
		return err
	}
	return db.Converge()
}
