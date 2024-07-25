package main

import (
	"context"
	"database/sql"
	"github.com/knaka/biblioseeq/db"
	"github.com/knaka/biblioseeq/db/sqlcgen"
	_ "github.com/mattn/go-sqlite3"
	"log"

	. "github.com/knaka/go-utils"
)

func foo() {
	var dbConn *sql.DB = nil
	ctx := context.Background()
	dbFilePath := "/tmp/tmp.db"
	db.Migrate(dbFilePath)
	dbConn = V(sql.Open("sqlite3", dbFilePath))
	ctx = context.Background()
	store := sqlcgen.New(dbConn)
	ver := V(store.GetVersion(ctx))
	log.Println("Version:", ver)
}

func main() {
	foo()
	//config := V(conf.ReadConfig())
	//opts := []fts.Option{}
	//for _, directory := range config.Directories {
	//	opts = append(opts, fts.WithDirectory(directory.Path, directory.FileExtensions))
	//}
	//engine := fts.NewIndex(opts...)
	//go func() {
	//	engine.Watch()
	//}()
	//for {
	//	// Read query strings from stdin
	//	scanner := bufio.NewScanner(os.Stdin)
	//	V0(os.Stdout.WriteString("> "))
	//	for scanner.Scan() {
	//		query := scanner.Text()
	//		log.Println("Query:", query)
	//		V0(os.Stdout.WriteString("> "))
	//	}
	//}
}
