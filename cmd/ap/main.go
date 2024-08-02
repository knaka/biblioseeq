package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bufbuild/connect-go"
	"github.com/chzyer/readline"
	ftslog "github.com/knaka/biblioseeq/log"
	v1 "github.com/knaka/biblioseeq/pbgen/v1"
	"github.com/knaka/biblioseeq/pbgen/v1/v1connect"
	"github.com/knaka/biblioseeq/search"
	"github.com/knaka/biblioseeq/tokenizer"
	"github.com/knaka/biblioseeq/web"
	. "github.com/knaka/go-utils"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	hostArg := flag.String("host", "", "host")
	portArg := flag.Int("port", 3000, "port")
	flag.Parse()
	host := *hostArg
	port := *portArg
	ftslog.SetOutput(os.Stderr)
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
	server := V(web.NewServer(host, port, searchEngine))
	host, port = V2(web.ParseServerAddr(server.Addr))
	serverUrl := fmt.Sprintf("http://%s:%d", Elvis(host, "localhost"), port)
	log.Printf("Listening on %s .", serverUrl)
	wg.Add(1)
	go (func() {
		Expect(server.ListenAndServe(), http.ErrServerClosed)
		wg.Done()
	})()
	httpClt := &http.Client{Transport: &http.Transport{}}
	clt := v1connect.NewMainServiceClient(httpClt, serverUrl)
	ctx.Done()
	currentTime := V(clt.CurrentTime(ctx,
		&connect.Request[v1.CurrentTimeRequest]{Msg: &v1.CurrentTimeRequest{}},
	))
	log.Println("e289984", currentTime.Msg.Timestamp.AsTime())
	searchEngine.WaitForInitialScanFinished(ctx)
	for {
		status := V(clt.Status(ctx, &connect.Request[v1.StatusRequest]{Msg: &v1.StatusRequest{}}))
		if status.Msg.InitialScanFinished {
			break
		}
		time.Sleep(1 * time.Second)
	}
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
		queryResult, err := clt.Query(ctx, &connect.Request[v1.QueryRequest]{Msg: &v1.QueryRequest{Query: *line}})
		if err != nil {
			log.Println("Error:", err)
			V0(os.Stdout.WriteString("> "))
			continue
		}
		for _, result := range queryResult.Msg.Results {
			snippet := result.Snippet
			snippet = strings.ReplaceAll(snippet, "\r", "")
			snippet = strings.ReplaceAll(snippet, "\n", " ")
			snippet = tokenizer.RemoveZWSP(snippet)
			V0(os.Stdout.WriteString(fmt.Sprintln(result.Path, snippet)))
		}
	}
	V0(server.Shutdown(ctx))
	cancel()
	wg.Wait()
}
