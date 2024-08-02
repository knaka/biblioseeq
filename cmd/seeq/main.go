package main

import (
	"flag"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/knaka/biblioseeq/conf"
	"github.com/knaka/biblioseeq/log"
	"github.com/knaka/biblioseeq/search"
	"github.com/knaka/biblioseeq/search/internal/tokenizer"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
	"strings"

	. "github.com/knaka/go-utils"
)

func main() {
	Debugger()

	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()
	log.SetOutput(Ternary[io.Writer](*verbose, os.Stderr, io.Discard))

	config := V(conf.Read())

	ftsOpts := []search.Option{
		search.WithDefaultDBFilePath(),
		search.MigratesDB(),
	}
	for _, confDir := range config.Directories {
		ftsOpts = append(ftsOpts, search.WithTargetDirectory(
			confDir.AbsPath,
			confDir.FileExtensions,
		))
	}
	ftsIndexer := search.NewIndexer(ftsOpts...)

	log.Println("Starting indexer.")
	go ftsIndexer.WatchDirsContinuously()

	log.Println("Waiting for initial scan to finish.")
	ftsIndexer.WaitForInitialScanFinished()
	log.Println("Initial scan finished.")

	rl := V(readline.NewEx(&readline.Config{
		Prompt: "> ",
	}))
	defer (func() { V0(rl.Close()) })()
	rl.CaptureExitSignal()
	for {
		line := PR(rl.Readline()).NilIf(io.EOF)
		if line == nil {
			break
		}
		*line = strings.TrimSpace(*line)
		if *line == "" {
			continue
		}
		query := tokenizer.SeparateJapanese(*line)
		log.Println("query:", query)
		results, err := ftsIndexer.Query(query)
		if err != nil {
			log.Println("Error:", err)
			V0(os.Stdout.WriteString("> "))
			continue
		}
		for _, result := range results {
			snippet := result.Snippet
			snippet = strings.ReplaceAll(snippet, "\r", "")
			snippet = strings.ReplaceAll(snippet, "\n", " ")
			snippet = tokenizer.RemoveZwsp(snippet)
			V0(os.Stdout.WriteString(fmt.Sprintln(result.Path, snippet)))
		}
	}
}
