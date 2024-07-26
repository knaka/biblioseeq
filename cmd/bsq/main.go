package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"github.com/knaka/biblioseeq/conf"
	"github.com/knaka/biblioseeq/db"
	"github.com/knaka/biblioseeq/fts"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"

	. "github.com/knaka/go-utils"
)

func main() {
	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()

	if *verbose {
		AddLogWriter(os.Stderr)
	}

	//ctx := context.Background()
	dbFilePath := "/tmp/tmp.db"
	db.Migrate(dbFilePath)

	config := V(conf.ReadConfig())
	opts := []fts.Option{}
	for _, directory := range config.Directories {
		opts = append(opts, fts.WithDirectory(directory.Path, directory.FileExtensions))
	}
	dbConn := V(sql.Open("sqlite3", dbFilePath))
	indexer := fts.NewIndexer(dbConn, opts...)
	LogPrintln("Starting indexer.")
	go indexer.WatchContinuously()
	LogPrintln("Waiting for initial scan to finish.")
	indexer.WaitForInitialScanFinished()
	LogPrintln("Initial scan finished.")
root:
	for {
		scanner := bufio.NewScanner(os.Stdin)
		V0(os.Stdout.WriteString("> "))
		for scanner.Scan() {
			query := strings.TrimSpace(scanner.Text())
			if query == "" {
				break root
			}
			query = fts.DivideJapaneseToWords(query)
			log.Println("Query:", query)
			results, err := indexer.Query(query)
			if err != nil {
				log.Println("Error:", err)
				V0(os.Stdout.WriteString("> "))
				continue
			}
			for _, result := range results {
				snippet := result.Snippet
				snippet = strings.ReplaceAll(snippet, "\r", "")
				snippet = strings.ReplaceAll(snippet, "\n", " ")
				snippet = fts.RemoveZwsp(snippet)
				V0(os.Stdout.WriteString(fmt.Sprintln(result.Path, snippet)))
			}
			V0(os.Stdout.WriteString("> "))
		}
	}
}
