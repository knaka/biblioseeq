//go:build sqlite_fts5 || fts5

package db

import (
	_ "embed"
	"github.com/sqldef/sqldef"
	"github.com/sqldef/sqldef/database"
	"github.com/sqldef/sqldef/database/sqlite3"
	"github.com/sqldef/sqldef/parser"
	"github.com/sqldef/sqldef/schema"

	. "github.com/knaka/go-utils"
)

//go:embed schema.sql
var dbSchema string

func Migrate(dbFilePath string) {
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
