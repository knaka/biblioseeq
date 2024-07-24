package main

import (
	"context"
	"flag"
	"github.com/knaka/biblioseeq/web"
	. "github.com/knaka/go-utils"
	"log"
)

func Main(host string, port int) (err error) {
	defer Catch(&err)
	server := V(web.NewServer(context.Background(), host, port))
	host, port = V2(web.ParseServerAddr(server.Addr))
	log.Printf("Listening on http://%s:%d .",
		Elvis(host, "localhost"),
		port,
	)
	V0(server.ListenAndServe())
	return
}

func main() {
	hostArg := flag.String("host", "", "host")
	portArg := flag.Int("port", 0, "port")
	flag.Parse()
	err := Main(*hostArg, *portArg)
	if err != nil {
		log.Panicf("panic 9f0b7b2 (%v)", err)
	}
}
