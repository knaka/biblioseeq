package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/chzyer/readline"
	"github.com/knaka/biblioseeq/conf"
	"github.com/knaka/biblioseeq/fts"
	ftslog "github.com/knaka/biblioseeq/log"
	"github.com/knaka/biblioseeq/web"
	. "github.com/knaka/go-utils"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	hostArg := flag.String("host", "", "host")
	portArg := flag.Int("port", 3000, "port")
	flag.Parse()
	host := *hostArg
	port := *portArg
	ftslog.SetOutput(os.Stderr)
	config := V(conf.ReadConfig())
	ftsOpts := []fts.Option{
		fts.WithDefaultDBFilePath(),
		fts.MigratesDB(),
	}
	for _, confDir := range config.Directories {
		ftsOpts = append(ftsOpts, fts.WithTargetDirectory(
			confDir.Path,
			confDir.FileExtensions,
		))
	}
	ftsIndexer := fts.NewIndexer(ftsOpts...)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer (func() { cancel() })()
	wg.Add(1)
	go (func() {
		ftsIndexer.WatchContinuously(ctx)
		wg.Done()
	})()
	server := V(web.NewServer(host, port, ftsIndexer))
	host, port = V2(web.ParseServerAddr(server.Addr))
	log.Printf("Listening on http://%s:%d .",
		Elvis(host, "localhost"),
		port,
	)
	wg.Add(1)
	go (func() {
		Ignore(server.ListenAndServe())
		wg.Done()
	})()
	ftsIndexer.WaitForInitialScanFinished(ctx)
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
			break
		}
		query := fts.SeparateJapanese(*line)
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
			snippet = fts.RemoveZwsp(snippet)
			V0(os.Stdout.WriteString(fmt.Sprintln(result.Path, snippet)))
		}
	}
	V0(server.Shutdown(ctx))
	cancel()
	wg.Wait()
}
