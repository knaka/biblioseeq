package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/knaka/biblioseeq/log"
	"github.com/knaka/biblioseeq/search"
	"github.com/knaka/biblioseeq/tokenizer"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
	"strings"
	"sync"

	. "github.com/knaka/go-utils"
)

func main() {
	Debugger()

	verbose := flag.Bool("v", false, "verbose")
	flag.Parse()
	log.SetOutput(Ternary[io.Writer](*verbose, os.Stderr, io.Discard))

	log.SetOutput(os.Stderr)
	ftsOpts := []search.Option{
		search.ShouldMigratesDB(true),
	}
	searchEngine := search.NewEngine(ftsOpts...)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer (func() { cancel() })()
	wg.Add(1)
	go (func() {
		searchEngine.Serve(ctx)
		wg.Done()
	})()

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
		results, err := searchEngine.Query(query)
		if err != nil {
			log.Println("Error:", err)
			V0(os.Stdout.WriteString("> "))
			continue
		}
		for _, result := range results {
			snippet := result.Snippet
			snippet = strings.ReplaceAll(snippet, "\r", "")
			snippet = strings.ReplaceAll(snippet, "\n", " ")
			snippet = tokenizer.RemoveZWSP(snippet)
			V0(os.Stdout.WriteString(fmt.Sprintln(result.Path, snippet)))
		}
	}
}
