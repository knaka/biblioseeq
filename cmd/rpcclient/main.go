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
	x := V(clt.GetVersionInfo(ctx,
		&connect.Request[v1.GetVersionInfoRequest]{Msg: &v1.GetVersionInfoRequest{}},
	))
	log.Println(x)
}
