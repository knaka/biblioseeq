//go:build sqlite_fts5 || fts5

package db

import (
	"context"
	"database/sql"
	_ "embed"
	. "github.com/knaka/go-utils"
	"github.com/sqldef/sqldef"
	"github.com/sqldef/sqldef/database"
	"github.com/sqldef/sqldef/database/sqlite3"
	"github.com/sqldef/sqldef/parser"
	"github.com/sqldef/sqldef/schema"
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

// UpdateDocument updates a document in the database.
// Unfortunately, sqlcgen does not support rowid of FTS5 table.
func UpdateDocument(ctx context.Context, db *sql.DB, body string, rowid int64) (err error) {
	_, err = db.ExecContext(ctx,
		"UPDATE documents SET body = ? WHERE rowid = ?",
		body, rowid,
	)
	return
}
