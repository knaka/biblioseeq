package web

import (
	"context"
	"github.com/knaka/biblioseeq/pbgen/v1/v1connect"
	"github.com/knaka/biblioseeq/web/rpc"
	"net"
	"net/http"
)

func ListenAndServe(ctx context.Context, addr string) error {
	server := &http.Server{Addr: addr, Handler: GetWrappedRouter()}
	server.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}
	return server.ListenAndServe()
}

func GetWrappedRouter() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(v1connect.NewMainServiceHandler(&rpc.MainServiceHandlerImpl{}))
	return mux
}
