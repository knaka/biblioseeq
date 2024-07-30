package main

import (
	"context"
	"github.com/bufbuild/connect-go"
	"github.com/knaka/biblioseeq/pbgen/v1"
	"github.com/knaka/biblioseeq/pbgen/v1/v1connect"
	"log"
	"net/http"

	. "github.com/knaka/go-utils"
)

func main() {
	httpClt := &http.Client{Transport: &http.Transport{}}
	clt := v1connect.NewMainServiceClient(httpClt, "http://127.0.0.1:3000")
	ctx := context.Background()
	ctx.Done()
	x := V(clt.VersionInfo(ctx,
		&connect.Request[v1.VersionInfoRequest]{Msg: &v1.VersionInfoRequest{}},
	))
	log.Println(x)
}
