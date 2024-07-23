package main

import (
	"context"
	"flag"
	"github.com/knaka/biblioseeq/web"
	"log"
)

func main() {
	ctx := context.Background()
	hostArg := flag.String("host", "", "host")
	portArg := flag.Int("port", 0, "port")
	server, host, port, err := web.NewServer(ctx, *hostArg, *portArg)
	if err != nil {
		log.Panicf("panic 9b807b2 (%v)", err)
	}
	log.Printf("Listening on http://%s:%d .", host, port)
	// Always returns a non-nil error.
	_ = server.ListenAndServe()
}
