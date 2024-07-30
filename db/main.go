package db

import (
	"database/sql"
	_ "embed"
	_ "github.com/knaka/go-sqlite3-fts5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sqldef/sqldef"
	"github.com/sqldef/sqldef/database"
	"github.com/sqldef/sqldef/database/sqlite3"
	"github.com/sqldef/sqldef/parser"
	"github.com/sqldef/sqldef/schema"
	"os"
	"path/filepath"

	. "github.com/knaka/go-utils"
)

//go:embed schema.sql
var dbSchema string

func EnsureDBFile(dbFilePath string) {
	if _, err := os.Stat(dbFilePath); err != nil && os.IsNotExist(err) {
		V0(os.MkdirAll(filepath.Dir(dbFilePath), 0755))
		V0(V(sql.Open("sqlite3", dbFilePath)).Close())
	}
}

func Migrate(dbFilePath string) {
	EnsureDBFile(dbFilePath)
	db := V(sqlite3.NewDatabase(database.Config{DbName: dbFilePath}))
	defer (func() { V0(db.Close()) })()
	sqldef.Run(schema.GeneratorModeSQLite3, db,
		database.NewParser(parser.ParserModeSQLite3),
		&sqldef.Options{
			DesiredDDLs: dbSchema,
			Config:      database.GeneratorConfig{},
		},
	)
}
