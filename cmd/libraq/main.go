package main

import (
	"context"
	"database/sql"
	. "github.com/knaka/go-utils"
	"github.com/knaka/libraq/db"
	"github.com/knaka/libraq/db/sqlcgen"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
)

func main() {
	dbFile := V(filepath.Abs(os.Args[1]))
	db.Migrate(dbFile)
	ctx := context.Background()
	dbConn := V(sql.Open("sqlite3", dbFile))
	store := sqlcgen.New(dbConn)
	version := V(store.GetVersion(ctx))
	log.Printf("sqlite version: %s", version)
}
